/*
this test requires Gotenberg and Tika to be running, and the TIKA_PATH and
GOTENBERG_PATH env variables to both be set to the correct URLs of the services
*/

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/decode"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"github.com/gorilla/mux"
)

var testRouter *mux.Router

var cookies []*http.Cookie
var admincookies []*http.Cookie

var testUsers = map[string][]map[string]string{
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
		map[string]string{
			"name":     "fileUser",
			"password": "filePassword",
			"email":    "third@example.com",
		},
	},
	"admin": []map[string]string{
		map[string]string{
			"name":     "admin",
			"password": "adminPass!",
			"email":    "admin@example.com",
		},
		map[string]string{ // will own public files
			"name":     "adminTwo",
			"password": "adminPass!",
			"email":    "admintwo@example.com",
		},
	},
}

type testFile struct {
	file    types.FileI
	store   *types.FileStore
	ctype   string
	content string
	tags    []tag.Tag
}

var testFiles = []testFile{
	testFile{
		file: &types.File{
			Name: "first.txt",
		},
		ctype:   "text/plain",
		content: "this is the first test file.",
		tags: []tag.Tag{
			tag.Tag{
				Word: "a",
				Type: tag.TOPIC,
				Data: tag.Data{
					tag.TOPIC: map[string]interface{}{
						"significance": 0,
						"count":        42,
						"first":        0,
					},
				},
			},
			tag.Tag{
				Word: "b",
				Type: tag.TOPIC,
				Data: tag.Data{
					tag.TOPIC: map[string]interface{}{
						"significance": 1,
						"count":        41,
						"first":        0,
					},
				},
			},
			tag.Tag{
				Word: "c",
				Type: tag.TOPIC,
				Data: tag.Data{
					tag.TOPIC: map[string]interface{}{
						"significance": 2,
						"count":        40,
						"first":        0,
					},
				},
			},
		},
	},
	testFile{
		file: &types.File{
			Name: "second.txt",
		},
		ctype:   "text/plain",
		content: "This is the second file.",
	},
	testFile{
		file: &types.File{
			Name: "third.txt",
		},
		ctype:   "text/plain",
		content: "This is the third file.",
	},
}

var adminFiles = []testFile{
	testFile{
		file: &types.File{
			Name: "admin.txt",
		},
		ctype:   "text/plain",
		content: "this is an admin's file.",
	},
	testFile{
		file: &types.File{
			Name: "secrets.txt",
		},
		ctype:   "text/plain",
		content: "this is the second admin's file.",
	},
}

var publicFiles = []testFile{
	testFile{
		file: &types.File{
			Name: "public1.txt",
		},
		ctype:   "text/plain",
		content: "this is a public file.",
	},
	testFile{
		file: &types.File{
			Name: "public2.txt",
		},
		ctype:   "text/plain",
		content: "This is public file number two. It has two sentences!",
	},
	testFile{
		file: &types.File{
			Name: "public3.txt",
		},
		ctype:   "text/plain",
		content: "The quick brown fox jumped over the lazy dog.",
	},
	testFile{
		file: &types.File{
			Name: "public4.txt",
		},
		ctype:   "text/plain",
		content: "Public files can come in all shapes and sizes!",
	},
}

func sliceContains(slice []string, s string) bool {
	for _, candidate := range slice {
		if candidate == s {
			return true
		}
	}
	return false
}

func TestMain(m *testing.M) {
	srverror.DEBUG = true
	testRouter = mux.NewRouter().PathPrefix("/api").Subrouter()
	testRouter.Use(Recovery)
	if err := populateDB(); err != nil {
		fmt.Println("error populating DB:", err)
		os.Exit(1)
	}

	// Attaching user handler, adding users, and setting timeouts so tests may
	// have a valid cookie
	AttachUser(testRouter.PathPrefix("/user").Subrouter())
	// probably attach all handlers together so they're in one place
	var configTimeout config.Duration
	configTimeout.Duration = time.Duration(10 * time.Second)
	config.V.UserTimeouts.Inactivity = configTimeout
	config.V.UserTimeouts.Total = configTimeout
	config.V.MinFileTimeout = configTimeout
	config.V.MaxFileTimeout = configTimeout

	status := m.Run()
	if status == 0 {
		if oc := memory.CurrentOpenConnections(); oc != 0 {
			status = 2
			fmt.Printf("Database Connections not handled: %d connections\n", oc)
		}
	}
	os.Exit(status)
}

// TODO: move config stuff to separate function
func populateDB() (err error) {
	setupctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	config.DB = new(memory.Database)
	if err = config.DB.Init(setupctx, true); err != nil {
		return
	}
	tikapath := os.Getenv("TIKA_PATH")
	if len(tikapath) == 0 {
		tikapath = "http://localhost:9998"
	}
	gotenpath := os.Getenv("GOTENBERG_PATH")
	if len(gotenpath) == 0 {
		gotenpath = "http://localhost:3000"
	}
	config.T.Path = tikapath
	config.V.GotenPath = gotenpath
	db, err := config.DB.Connect(setupctx)
	if err != nil {
		return
	}
	defer db.Close(setupctx)
	userbase := db.Owner()
	tagbase := userbase.Tag()
	for i, userdata := range testUsers["users"] {
		user := types.NewUser(userdata["name"], userdata["password"], userdata["email"])
		if _, err = userbase.Reserve(user.ID, user.Name); err != nil {
			return
		}
		if err = userbase.Insert(user); err != nil {
			return
		}
		userdata["id"] = user.GetID().String()
		switch v := testFiles[i].file.(type) {
		case *types.File:
			v.Own = user
		case *types.WebFile:
			v.Own = user
		}
	}
	for i, admindata := range testUsers["admin"] {
		admin := types.NewUser(admindata["name"], admindata["password"], admindata["email"])
		admin.SetRole("admin", true)
		if admin.ID, err = userbase.Reserve(admin.ID, admin.Name); err != nil {
			return err
		}
		if err = userbase.Insert(admin); err != nil {
			return err
		}
		admindata["id"] = admin.GetID().String()
		switch v := adminFiles[i].file.(type) {
		case *types.File:
			v.Own = admin
		case *types.WebFile:
			v.Own = admin
		}
	}
	for i, file := range adminFiles {
		adminFiles[i].store, err = process.InjestFile(setupctx, file.file, file.ctype, strings.NewReader(file.content), userbase)
		if err != nil {
			return
		}
		decode.Read(setupctx, nil, adminFiles[i].store, config.DB, config.T.Path, config.V.GotenPath)
		if err != nil {
			return
		}
	}
	var publicOwnerID types.OwnerID
	publicOwnerID, err = types.DecodeObjectIDString(testUsers["admin"][1]["id"])
	if err != nil {
		return
	}
	var publicOwner types.Owner
	publicOwner, err = userbase.Get(publicOwnerID)
	if err != nil {
		return
	}
	for i, file := range publicFiles {
		switch f := file.file.(type) {
		case *types.File:
			f.Own = publicOwner
		case *types.WebFile:
			f.Own = publicOwner
		}
		file.file.SetPerm(types.Public, "view", true)
		publicFiles[i].store, err = process.InjestFile(setupctx, file.file, file.ctype, strings.NewReader(file.content), userbase)
		if err != nil {
			fmt.Printf("injest file failed")
			return
		}
		decode.Read(setupctx, nil, publicFiles[i].store, config.DB, config.T.Path, config.V.GotenPath)
		if err != nil {
			return
		}
	}
	for i, file := range testFiles {
		testFiles[i].store, err = process.InjestFile(setupctx, file.file, file.ctype, strings.NewReader(file.content), userbase)
		if err != nil {
			return
		}
		decode.Read(setupctx, nil, testFiles[i].store, config.DB, config.T.Path, config.V.GotenPath)
		if err != nil {
			return
		}
		// fmt.Printf("i:%d, ID:%+#v", i, testFiles[i].file.GetID())
		if i > 0 {
			perm := file.file.(types.PermissionI)
			var targetUser types.UserI
			targetUser, err = userbase.FindUserName(testUsers["users"][i-1]["name"])
			if err != nil {
				return
			}
			perm.SetPerm(targetUser, "view", true)
		}
		for _, t := range file.tags {
			ftag := tag.FileTag{
				File:  file.file.GetID(),
				Owner: file.file.GetOwner().GetID(),
				Tag:   t,
			}
			tagbase.Upsert(ftag)
		}
	}
	return nil
}

func responseBodyString(res *httptest.ResponseRecorder) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Result().Body)
	return buf.String()
}

func testlogin(t *testing.T, i int, admin bool) []*http.Cookie {
	var userColl string
	if admin {
		userColl = "admin"
	} else {
		userColl = "users"
	}
	ldata := map[string]string{
		"name": testUsers[userColl][i]["name"],
		"pass": testUsers[userColl][i]["password"],
	}
	loginbody, _ := json.Marshal(ldata)
	loginreq, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(loginbody))
	loginreq.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, loginreq)
	if response.Code != 200 {
		t.Fatalf("login non success status code: %+#v", response)
	}
	return response.Result().Cookies()
}
