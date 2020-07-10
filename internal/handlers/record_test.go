/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const userIdx = 0

func setupRecord(t *testing.T) {
	AttachRecord(testRouter.PathPrefix("/record").Subrouter())
	cookies = testlogin(t, userIdx, false)
}

func TestRecordAPI(t *testing.T) {
	setupRecord(t)
	t.Run("GetOwned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/record", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		// t.Logf("Body: %s", responseBodyString(res))
		var result struct {
			Files map[string]struct {
				Count int `json:"count"`
				File  struct {
					Name string `json:"name"`
				} `json:"file"`
			} `json:"files"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&result)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(result.Files) != 1 {
			t.Fatalf("Expected number of files to be 1, received: %+#v", result)
		}
		for _, v := range result.Files {
			name := v.File.Name
			expected := testFiles[userIdx].file.GetName()
			if name != expected {
				t.Fatalf("Expected returned file to be %s, but got %s", expected, name)
			}
		}
	})
	t.Run("GetView", func(t *testing.T) {
		/*
		 * each file in testFiles in handlers_test is shared with the i-1th user,
		 * so expect this user to have a shared file from user i + 1, as long as
		 * our setup logged us in as a user that isn't the last in the array
		 */
		req, _ := http.NewRequest("GET", "/api/record/view", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var result struct {
			Files map[string]struct {
				File struct {
					Name string `json:"name"`
					Own  string `json:"own"`
					Perm struct {
						View []string `json:"view"`
					} `json:"perm"`
				} `json:"file"`
			} `json:"files"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&result)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(result.Files) != 1 {
			t.Fatalf("Expected 1 result: %+#v", result.Files)
		}
		for _, fileWrapper := range result.Files {
			file := fileWrapper.File
			expected := testFiles[userIdx+1].file
			us := testUsers["users"][userIdx]
			sharer := testUsers["users"][userIdx+1]
			if file.Name != expected.GetName() {
				t.Fatalf("Received file name: %s, expected %s", file.Name, expected.GetName())
			}
			if file.Own != sharer["id"] {
				t.Fatalf("Received file owner id: %s, expected %s", file.Own, sharer["id"])
			}
			if !sliceContains(file.Perm.View, us["id"]) {
				t.Fatalf("Received file shared with: %+#v, expected %s", file.Perm.View, us["id"])
			}
		}
	})
	t.Run("ChangeRecordNameBadID", func(t *testing.T) {
		// id := testFiles[0].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/record/badid/name", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 400 {
			t.Fatalf("expected status code 400: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("ChangeRecordNameNoName", func(t *testing.T) {
		id := testFiles[userIdx].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/record/"+id+"/name", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 400 {
			t.Fatalf("expected status code 400: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("ChangeRecordNamePermissionDenied", func(t *testing.T) {
		id := testFiles[userIdx+1].file.GetID().String()
		req, _ := http.NewRequest("POST", "/api/record/"+id+"/name", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("expected status code 403: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("ChangeRecordNameSuccess", func(t *testing.T) {
		id := testFiles[userIdx].file.GetID().String()
		newname := "renamed.txt"
		params := map[string]string{
			"name": newname,
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("POST", "/api/record/"+id+"/name", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("expected status code 200: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that name actually changed
		{
			req, _ := http.NewRequest("GET", "/api/record", nil)
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			res := httptest.NewRecorder()
			testRouter.ServeHTTP(res, req)
			if res.Code != 200 {
				t.Fatalf("non success status code when retrieving file after changing name: %+#v\nBody:%s", res, responseBodyString(res))
			}
			var result struct {
				Files map[string]struct {
					File struct {
						Name string `json:"name"`
					} `json:"file"`
				} `json:"files"`
			}
			err := json.NewDecoder(res.Result().Body).Decode(&result)
			if err != nil {
				t.Fatalf("JSON Decode error:\n%s", err)
			}
			if len(result.Files) != 1 {
				t.Fatalf("Expected number of files to be 1, received: %+#v", result)
			}
			for _, v := range result.Files {
				name := v.File.Name
				if name != newname {
					t.Fatalf("Expected returned file to be %s, but got %s", newname, name)
				}
			}
		}
		//change name back in case other tests are expecting certain file names
		params = map[string]string{
			"name": testFiles[userIdx].file.GetName(),
		}
		jsonbytes, _ = json.Marshal(params)
		req, _ = http.NewRequest("POST", "/api/record/"+id+"/name", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("expected status code 200: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
}
