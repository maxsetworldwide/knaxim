package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/mongo"

	"github.com/google/go-tika/tika"
)

var testingconfiguration = configuration{
	StaticPath:      ".",
	GracefulTimeout: time.Minute * 3,
	BasicTimeout:    time.Second * 5,
	FileTimeoutRate: 1000000000,
	MaxFileTimeout:  time.Minute * 30,
	MinFileTimeout:  time.Second * 10,
	Tika: tikaconf{
		Type: "local",
		Path: "/usr/lib/tika/tika-server-1.21.jar",
	},
	FileLimit: 50 * 1024 * 1024,
	AdminKey:  "testadminkey",
	GuestUser: &guestconf{
		Name:  "guest",
		Pass:  "guestpass",
		Email: "info@maxset.org",
	},
}

var testingdb database.Database = &mongo.Database{
	URI:    "mongodb://localhost:27017",
	DBName: "AutoTest",
}

var testingdata = map[string][]map[string]string{
	"users": []map[string]string{
		map[string]string{
			"name":     "first",
			"password": "1Password!",
			"email":    "first@example.com",
		},
		map[string]string{
			"name":     "second",
			"password": "2Password!",
			"email":    "second@example.com",
		},
	},
	"admin": []map[string]string{
		map[string]string{
			"name":     "admin",
			"password": "adminPass!",
			"email":    "admin@example.com",
		},
	},
}

var testingacronyms = map[string][]string{
	"db": []string{"database"},
	"ab": []string{"acronymbase", "another babboon"},
}

type testfile struct {
	file    database.FileI
	store   *database.FileStore
	ctype   string
	content string
}

var testingfiles = []testfile{
	testfile{
		file: &database.File{
			Name: "first.txt",
		},
		ctype:   "text/plain",
		content: "this is the first test file.",
	},
	testfile{
		file: &database.WebFile{
			File: database.File{
				Name: "second.html",
			},
			URL: "localhost/second.txt",
		},
		ctype:   "text/html",
		content: "<p>This is the second file.</p>",
	},
}

func populateDB() error {
	var err error
	setupctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	if err := db.Init(setupctx, true); err != nil {
		return err
	}
	if conf.Tika.Type == "local" {
		var err error
		tserver, err = tika.NewServer(conf.Tika.Path, conf.Tika.Port)
		if err != nil {
			return err
		}
		tserver.ChildMode(&tika.ChildOptions{
			MaxFiles:          conf.Tika.MaxFiles,
			TaskPulseMillis:   conf.Tika.TaskPulse,
			TaskTimeoutMillis: conf.Tika.TaskTimeout,
			PingPulseMillis:   conf.Tika.PingPulse,
			PingTimeoutMillis: conf.Tika.PingTimeout,
		})
		defer cancel()
		if err := tserver.Start(setupctx); err != nil {
			return err
		}
		tikapath = tserver.URL()
	} else if conf.Tika.Type == "external" {
		if conf.Tika.Port == "" {
			conf.Tika.Port = "9998"
		}
		tikapath = conf.Tika.Path + ":" + conf.Tika.Port
	} else {
		return fmt.Errorf("unrecognized Tika Type")
	}
	userbase := db.Owner(setupctx)
	defer userbase.Close(setupctx)
	if conf.GuestUser != nil {
		guestUser := database.NewUser(conf.GuestUser.Name, conf.GuestUser.Pass, conf.GuestUser.Email)
		guestUser.SetRole("Guest", true)
		if _, err := userbase.FindUserName(conf.GuestUser.Name); err == database.ErrNotFound {
			if guestUser.ID, err = userbase.Reserve(guestUser.ID, guestUser.Name); err != nil {
				return err
			}
			if err := userbase.Insert(guestUser); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	//loading users
	for i, userdata := range testingdata["users"] {
		user := database.NewUser(userdata["name"], userdata["password"], userdata["email"])
		if user.ID, err = userbase.Reserve(user.ID, user.Name); err != nil {
			return err
		}
		if err = userbase.Insert(user); err != nil {
			return err
		}
		userdata["id"] = user.GetID().String()
		switch v := testingfiles[i].file.(type) {
		case *database.File:
			v.Own = user
		case *database.WebFile:
			v.Own = user
		}
	}
	for _, admindata := range testingdata["admin"] {
		admin := database.NewUser(admindata["name"], admindata["password"], admindata["email"])
		admin.SetRole("admin", true)
		if admin.ID, err = userbase.Reserve(admin.ID, admin.Name); err != nil {
			return err
		}
		if err = userbase.Insert(admin); err != nil {
			return err
		}
		admindata["id"] = admin.GetID().String()
	}
	for i, tf := range testingfiles {
		testingfiles[i].store, err = database.InjestFile(setupctx, tf.file, tf.ctype, strings.NewReader(tf.content), db)
		if err != nil {
			return err
		}
		err = processContent(setupctx, nil, tf.file, testingfiles[i].store)
		if err != nil {
			return err
		}
	}
	ab := db.Acronym(setupctx)
	defer ab.Close(setupctx)
	for acronym, options := range testingacronyms {
		for _, option := range options {
			err := ab.Put(acronym, option)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testlogin(t *testing.T, i int) []*http.Cookie {
	t.Helper()
	ldata := map[string]string{
		"name": testingdata["users"][i]["name"],
		"pass": testingdata["users"][i]["password"],
	}
	loginbody, _ := json.Marshal(ldata)
	loginreq, _ := http.NewRequest("POST", server.URL+"/api/user/login", bytes.NewReader(loginbody))
	loginreq.Header.Add("Content-Type", "application/json")
	res, err := server.Client().Do(loginreq)
	if err != nil {
		t.Fatal("Unable to login", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("login non success status code: %+#v", res)
	}
	return res.Cookies()
}
