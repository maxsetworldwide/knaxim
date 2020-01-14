package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDir(t *testing.T) {
	cookies := testlogin(t, 0)
	t.Run("Create", func(t *testing.T) {
		t.Logf("adding testing dir to %s", testingfiles[0].file.GetID().String())
		vals := map[string]string{
			"newname": "testing",
			"content": testingfiles[0].file.GetID().String(),
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", server.URL+"/api/dir", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("server client error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Non Success Returned: %+#v", res)
		}
	})
	t.Run("Info", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/dir/testing", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("server client error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Non Success Returned: %+#v", res)
		}
	})
	t.Run("RemoveFile", func(t *testing.T) {
		vals := map[string]string{
			"id": testingfiles[0].file.GetID().String(),
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("DELETE", server.URL+"/api/dir/testing/content", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("server client error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Non Success Returned: %+#v", res)
		}
	})
	t.Run("AddFile", func(t *testing.T) {
		vals := map[string]string{
			"id": testingfiles[0].file.GetID().String(),
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", server.URL+"/api/dir/testing/content", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("server client error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Non Success Returned: %+#v", res)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", server.URL+"/api/dir/testing", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("server client error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("Non Success Returned: %+#v", res)
		}
	})
}
