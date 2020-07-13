package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.maxset.io/web/knaxim/internal/config"
)

func setupUserTest() {
	config.V.AdminKey = "adminKey1234"
}

func TestUserAPI(t *testing.T) {
	setupUserTest()
	t.Run("Create", func(t *testing.T) {
		vals := map[string]string{
			"name":  "testuser",
			"pass":  "testPass1!",
			"email": "test@example.com",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", "/api/user", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	// t.Run("CreateAdminWrongKey", func(t *testing.T) {
	// 	vals := map[string]string{
	// 		"name":     "testadmin",
	// 		"pass":     "testPassAdmin!1",
	// 		"email":    "admin@example.com",
	// 		"adminkey": "thisKeyShouldBeWrong",
	// 	}
	// 	jsonbytes, _ := json.Marshal(vals)
	// 	req, _ := http.NewRequest("PUT", "/api/user/admin", bytes.NewReader(jsonbytes))
	// 	req.Header.Add("Content-Type", "application/json")
	// 	res := httptest.NewRecorder()
	// 	testRouter.ServeHTTP(res, req)
	// 	if res.Code != 400 {
	// 		t.Fatalf("expected non success status code: %+#v", res)
	// 	}
	// })
	// t.Run("CreateAdmin", func(t *testing.T) {
	// 	vals := map[string]string{
	// 		"name":     "testadmin",
	// 		"pass":     "testPassAdmin!1",
	// 		"email":    "admin@example.com",
	// 		"adminkey": config.V.AdminKey,
	// 	}
	// 	jsonbytes, _ := json.Marshal(vals)
	// 	req, _ := http.NewRequest("PUT", "/api/user/admin", bytes.NewReader(jsonbytes))
	// 	req.Header.Add("Content-Type", "application/json")
	// 	res := httptest.NewRecorder()
	// 	testRouter.ServeHTTP(res, req)
	// 	if res.Code != 200 {
	// 		t.Fatalf("non success status code: %+#v", res)
	// 	}
	// })
	// t.Run("AdminLoginWrongPassword", func(t *testing.T) {
	// 	ldata := map[string]string{
	// 		"name": "testadmin",
	// 		"pass": "wrongPassword",
	// 	}
	// 	loginbody, _ := json.Marshal(ldata)
	// 	loginreq, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(loginbody))
	// 	loginreq.Header.Add("Content-Type", "application/json")
	// 	res := httptest.NewRecorder()
	// 	testRouter.ServeHTTP(res, loginreq)
	// 	if res.Code != 404 {
	// 		t.Fatalf("expected non success status code: %+#v", res)
	// 	}
	// })
	// t.Run("AdminLogin", func(t *testing.T) {
	// 	ldata := map[string]string{
	// 		"name": "testadmin",
	// 		"pass": "testPassAdmin!1",
	// 	}
	// 	loginbody, _ := json.Marshal(ldata)
	// 	loginreq, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(loginbody))
	// 	loginreq.Header.Add("Content-Type", "application/json")
	// 	res := httptest.NewRecorder()
	// 	testRouter.ServeHTTP(res, loginreq)
	// 	if res.Code != 200 {
	// 		t.Fatalf("login non success status code: %+#v", res)
	// 	}
	// 	cookies = res.Result().Cookies()
	// })
	t.Run("Login", func(t *testing.T) {
		logindata := map[string]string{
			"name": "testuser",
			"pass": "testPass1!",
		}
		body, _ := json.Marshal(logindata)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("login non success status code: %+#v", res)
		}
		cookies = res.Result().Cookies()
	})
	t.Run("UserInfo=Self", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("UserInfo=Other", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user?id="+testUsers["users"][0]["id"], nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("UserInfoMissingID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user?id=thisisabadid", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 404 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("UserComplete", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/complete", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/data", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("ChangePassWrongPass", func(t *testing.T) {
		vals := map[string]string{
			"oldpass": "wrongPassword",
			"newpass": "testPassAdmin!2",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", "/api/user/pass", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 404 {
			t.Fatalf("Expected 404 Status Code: %+#v", res)
		}
	})
	t.Run("ChangePass", func(t *testing.T) {
		vals := map[string]string{
			"oldpass": "testPass1!",
			"newpass": "testPass2!",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", "/api/user/pass", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("LookupMissingUser", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/name/missinguser", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 404 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("LookupUser", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/name/"+testUsers["users"][0]["name"], nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("SearchUserFiles", func(t *testing.T) {
		// just search for clean response
		vals := map[string]string{
			"find": "search term",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("GET", "/api/user/search", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
	t.Run("Logout", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/user", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Bad Status Code: %+#v", res)
		}
	})
}
