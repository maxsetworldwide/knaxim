// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func setupPerm(t *testing.T) {
	AttachPerm(testRouter.PathPrefix("/perm").Subrouter())
	cookies = testlogin(t, 0, false)
	admincookies = testlogin(t, 0, true) // needed for public perms
}

//mux type: "file" and "group"
//mux id: file id
//form id: user id to share with

//TODO: add tests for "group" type
func TestPermAPI(t *testing.T) {
	setupPerm(t)
	type returnObj struct {
		IsOwned    bool   `json:"isOwned"`
		Owner      string `json:"owner"`
		Permission struct {
			View []string `json:"view"`
		} `json:"permission"`
	}
	t.Run("GetFilePermissionAndEmpty", func(t *testing.T) {
		fid := testFiles[0].file.GetID().String()
		req, _ := http.NewRequest("GET", "/api/perm/file/"+fid, nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results returnObj
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if len(results.Permission.View) != 0 {
			t.Fatalf("Expected no viewers for file: response: %+#v\nBody:%s", results, responseBodyString(res))
		}
	})
	t.Run("SetFilePermissionDenied", func(t *testing.T) {
		fid := testFiles[1].file.GetID().String()
		uid := testUsers["users"][0]["id"]
		vals := map[string]string{
			"id": uid,
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", "/api/perm/file/"+fid, bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("expected 403: %+#v\nBody:%s", res, responseBodyString(res))
		}

	})
	t.Run("SetFilePermissionTrue", func(t *testing.T) {
		fid := testFiles[0].file.GetID().String()
		uid := testUsers["users"][1]["id"]
		vals := map[string]string{
			"id": uid,
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", "/api/perm/file/"+fid, bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that permission changed
		req, _ = http.NewRequest("GET", "/api/perm/file/"+fid, nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results returnObj
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if len(results.Permission.View) != 1 || results.Permission.View[0] != uid {
			t.Fatalf("Expected sole viewer to be %s: response: %+#v\nBody:%s", uid, results, responseBodyString(res))
		}
	})
	t.Run("SetFilePermissionFalse", func(t *testing.T) {
		fid := testFiles[0].file.GetID().String()
		uid := testUsers["users"][1]["id"]
		vals := map[string]string{
			"id": uid,
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("DELETE", "/api/perm/file/"+fid, bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that permission changed
		req, _ = http.NewRequest("GET", "/api/perm/file/"+fid, nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results returnObj
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if len(results.Permission.View) != 0 {
			t.Fatalf("Expected no viewers for file: response: %+#v\nBody:%s", results, responseBodyString(res))
		}
	})
	t.Run("SetPublicNoAdmin", func(t *testing.T) {
		fid := testFiles[0].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/perm/file/"+fid+"/public", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("expected 403: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("SetPublicNotOwned", func(t *testing.T) {
		fid := testFiles[0].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/perm/file/"+fid+"/public", nil)
		for _, cookie := range admincookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("expected 403: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("SetPublicTrue", func(t *testing.T) {
		fid := adminFiles[0].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/perm/file/"+fid+"/public", nil)
		for _, cookie := range admincookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that permission changed
		req, _ = http.NewRequest("GET", "/api/perm/file/"+fid, nil)
		for _, cookie := range admincookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results returnObj
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if len(results.Permission.View) != 1 {
			t.Fatalf("Expected a single viewer for file: response: %+#v\nBody:%s", results, responseBodyString(res))
		}
		receivedID := results.Permission.View[0]
		publicID := types.Public.GetID().String()
		if receivedID != publicID {
			t.Fatalf("Expected received ID to be the public owner: got %s, expected %s", receivedID, publicID)
		}
	})
	t.Run("SetPublicFalse", func(t *testing.T) {
		fid := adminFiles[0].file.GetID().String()
		req, _ := http.NewRequest("DELETE", "/api/perm/file/"+fid+"/public", nil)
		for _, cookie := range admincookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that permission changed
		req, _ = http.NewRequest("GET", "/api/perm/file/"+fid, nil)
		for _, cookie := range admincookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results returnObj
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if len(results.Permission.View) != 0 {
			t.Fatalf("Expected no viewers for file: response: %+#v\nBody:%s", results, responseBodyString(res))
		}
	})
}
