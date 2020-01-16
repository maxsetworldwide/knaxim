package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUserAPI(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		vals := map[string]string{
			"name":  "testuser",
			"pass":  "testPass1!",
			"email": "test@example.com",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", server.URL+"/api/user", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		client := server.Client()
		res, err := client.Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("error returned from client", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	t.Run("CreateAdmin", func(t *testing.T) {
		vals := map[string]string{
			"name":     "testadmin",
			"pass":     "testPassAdmin!1",
			"email":    "admin@example.com",
			"adminkey": conf.AdminKey,
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", server.URL+"/api/user/admin", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		client := server.Client()
		res, err := client.Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("error returned from client", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	//login
	ldata := map[string]string{
		"name": "testadmin",
		"pass": "testPassAdmin!1",
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
	cookies := res.Cookies()
	t.Run("UserInfo=Self", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Unable to Get User Info", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("UserInfo=Other", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user?id="+testingdata["users"][0]["id"], nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Unable to Get User Info", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("UserComplete", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user/complete", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Unable to Get User Info", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user/data", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("unable to get data", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("ChangePass", func(t *testing.T) {
		vals := map[string]string{
			"oldpass": "testPassAdmin!1",
			"newpass": "testPassAdmin!2",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", server.URL+"/api/user/pass", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Unable to Change password", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("Logout", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", server.URL+"/api/user", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Unable to Log Out", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
}
