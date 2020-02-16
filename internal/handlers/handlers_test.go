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
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/process"
	"github.com/gorilla/mux"
)

var testRouter *mux.Router

var cookies []*http.Cookie

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
	},
	"admin": []map[string]string{
		map[string]string{
			"name":     "admin",
			"password": "adminPass!",
			"email":    "admin@example.com",
		},
	},
}

type testFile struct {
	file    database.FileI
	store   *database.FileStore
	ctype   string
	content string
}

var testFiles = []testFile{
	testFile{
		file: &database.File{
			Name: "first.txt",
		},
		ctype:   "text/plain",
		content: "this is the first test file.",
	},
	testFile{
		file: &database.File{
			Name: "second.txt",
		},
		ctype:   "text/plain",
		content: "This is the second file.",
	},
}

func TestMain(m *testing.M) {
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
	configTimeout.Duration = time.Duration(9999999999)
	config.V.UserTimeouts.Inactivity = configTimeout
	config.V.UserTimeouts.Total = configTimeout

	os.Exit(m.Run())
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
	userbase := config.DB.Owner(setupctx)
	defer userbase.Close(setupctx)
	for i, userdata := range testUsers["users"] {
		user := database.NewUser(userdata["name"], userdata["password"], userdata["email"])
		if _, err = userbase.Reserve(user.ID, user.Name); err != nil {
			return
		}
		if err = userbase.Insert(user); err != nil {
			return
		}
		userdata["id"] = user.GetID().String()
		switch v := testFiles[i].file.(type) {
		case *database.File:
			v.Own = user
		case *database.WebFile:
			v.Own = user
		}
	}
	for i, file := range testFiles {
		testFiles[i].store, err = process.InjestFile(setupctx, file.file, file.ctype, strings.NewReader(file.content), userbase)
		if err != nil {
			return
		}
		err = processContent(setupctx, nil, testFiles[i].file, testFiles[i].store)
		if err != nil {
			return
		}
	}

	return nil
}

func responseBodyString(res *httptest.ResponseRecorder) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Result().Body)
	return buf.String()
}

func testlogin(t *testing.T, i int) []*http.Cookie {
	ldata := map[string]string{
		"name": testUsers["users"][i]["name"],
		"pass": testUsers["users"][i]["password"],
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
