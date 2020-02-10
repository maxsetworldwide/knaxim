package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupGroup(t *testing.T) {
	AttachGroup(testRouter.PathPrefix("/group").Subrouter())
	cookies = testlogin(t, 0)
}

func TestGroupAPI(t *testing.T) {
	setupGroup(t)
	t.Run("Create", func(t *testing.T) {
		vals := map[string]string{
			"newname": "TestGroup",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", "/api/group", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("GetGroups", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/group/options", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var vals struct {
			Member []map[string]string `json:"member"`
			Own    []map[string]string `json:"own"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&vals)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if vals.Own[0]["name"] != "TestGroup" {
			t.Fatalf("Expected an owned group named 'TestGroup': %+#v", vals.Own[0])
		}
	})
	var firstGroupID string
	t.Run("Lookup", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/group/name/TestGroup", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var vals struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&vals)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if vals.Name != "TestGroup" {
			t.Fatalf("Expected requested group name, got '%s'", vals.Name)
		}
		firstGroupID = vals.ID
	})
	t.Run("LookupMissing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/group/name/missingGroup", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 404 {
			t.Fatalf("expected 404 status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("GroupInfo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/group/"+firstGroupID, nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var vals struct {
			ID string `json:"id"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&vals)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		if vals.ID != firstGroupID {
			t.Fatalf("Expected id '%s' to be returned, but got '%s'", firstGroupID, vals.ID)
		}
	})
	t.Run("AddMember", func(t *testing.T) {
		params := map[string]string{
			"id": testUsers["users"][1]["id"],
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("POST", "/api/group/"+firstGroupID+"/member", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("RemoveMember", func(t *testing.T) {
		params := map[string]string{
			"id": testUsers["users"][1]["id"],
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("DELETE", "/api/group/"+firstGroupID+"/member", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("AddMemberBadID", func(t *testing.T) {
		params := map[string]string{
			"id": "badID",
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("POST", "/api/group/"+firstGroupID+"/member", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 400 {
			t.Fatalf("expected status code 400: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("RemoveMemberBadID", func(t *testing.T) {
		params := map[string]string{
			"id": "badID",
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("DELETE", "/api/group/"+firstGroupID+"/member", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 400 {
			t.Fatalf("expected status code 400: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	var parentID string
	t.Run("AddGroupToGroup", func(t *testing.T) {
		{
			// create parent group
			params := map[string]string{
				"newname": "ParentGroup",
			}
			jsonbytes, _ := json.Marshal(params)
			req, _ := http.NewRequest("PUT", "/api/group", bytes.NewReader(jsonbytes))
			req.Header.Add("Content-Type", "application/json")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			res := httptest.NewRecorder()
			testRouter.ServeHTTP(res, req)
			if res.Code != 200 {
				t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
			}
		}
		{
			// get id of new group
			req, _ := http.NewRequest("GET", "/api/group/name/ParentGroup", nil)
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			res := httptest.NewRecorder()
			testRouter.ServeHTTP(res, req)
			if res.Code != 200 {
				t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
			}
			var vals struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			}
			err := json.NewDecoder(res.Result().Body).Decode(&vals)
			if err != nil {
				t.Fatalf("error reading response:%s", err)
			}
			if vals.Name != "ParentGroup" {
				t.Fatalf("Expected requested group name, got '%s'", vals.Name)
			}
			parentID = vals.ID
		}
		// add first group under parent group
		params := map[string]string{
			"id": firstGroupID,
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("POST", "/api/group/"+parentID+"/member", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("expected status code 400: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
	t.Run("GetGroupsGroups", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/group/options/"+parentID, nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		type memberEntry struct {
			ID      string   `json:"id"`
			Name    string   `json:"name"`
			Members []string `json:"members"`
		}
		var jsonRes struct {
			Entry []memberEntry `json:"member"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&jsonRes)
		if err != nil {
			t.Fatalf("error reading response:%s", err)
		}
		found := false
		for i := 0; i < len(jsonRes.Entry); i++ {
			curr := jsonRes.Entry[i]
			for j := 0; j < len(curr.Members) && !found; j++ {
				found = curr.Members[j] == firstGroupID
			}
			if !found {
				t.Fatalf("Expected group '%s' to be a member of '%s'. List:%+#v", firstGroupID, parentID, jsonRes.Entry)
			}
		}
	})
	t.Run("SearchGroupFiles", func(t *testing.T) {
		params := map[string]string{
			"find": "searchTerm",
		}
		jsonbytes, _ := json.Marshal(params)
		req, _ := http.NewRequest("GET", "/api/group/"+firstGroupID+"/search", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
	})
}
