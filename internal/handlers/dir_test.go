package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupDir(t *testing.T) {
	AttachDir(testRouter.PathPrefix("/dir").Subrouter())
	cookies = testlogin(t, 0, false)
}

type creationResponse struct {
	AffectedFiles int    `json:"affectedFiles"`
	ID            string `json:"id"`
}

const dirname = "testdir"

func TestDirAPI(t *testing.T) {
	setupDir(t)
	file := testFiles[0]
	fid := file.file.GetID()
	t.Run("CreateEmpty", func(t *testing.T) {
		vals := map[string]string{
			"newname": dirname,
			"content": "",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", "/api/dir", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results creationResponse
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if results.AffectedFiles != 0 {
			t.Fatalf("Expected no affected files")
		}
	})
	t.Run("Create", func(t *testing.T) {
		vals := map[string]string{
			"newname": dirname,
			"content": fid.String(),
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", "/api/dir", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results creationResponse
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if results.AffectedFiles != 1 {
			t.Fatalf("Expected no affected files: Response:%+v", results)
		}
		if results.ID != dirname {
			t.Fatalf("Expected given directory name to be returned. Received'%s', Expected'%s'", results.ID, dirname)
		}
	})
	t.Run("GetAll", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/dir", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Folders []string `json:"folders"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(results.Folders) != 1 || results.Folders[0] != dirname {
			t.Fatalf("Expected %s to be the sole returned folder. Received %+v", dirname, results.Folders)
		}
	})
	t.Run("DirInfo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/dir/"+dirname, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Files []string `json:"files"`
			Name  string   `json:"name"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if results.Name != dirname {
			t.Fatalf("Expected queried dir name to be returned. Received'%s', expected'%s'", results.Name, dirname)
		}
		if len(results.Files) != 1 || results.Files[0] != fid.String() {
			t.Fatalf("Expected %s to be the sole returned file. Received %s", fid.String(), results.Files[0])
		}
	})
	t.Run("DirInfoNonExistent", func(t *testing.T) {
		missingDirName := "thisdirshouldnotexist"
		req, _ := http.NewRequest("GET", "/api/dir/"+missingDirName, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Files []string `json:"files"`
			Name  string   `json:"name"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if results.Name != missingDirName {
			t.Fatalf("Expected queried dir name to be returned. Received'%s', expected'%s'", results.Name, dirname)
		}
		if len(results.Files) != 0 {
			t.Fatalf("Expected dir to be empty. Received %+v", results.Files)
		}
	})
	t.Run("SearchDir", func(t *testing.T) {
		fileContent := strings.Split(file.content, " ")
		if len(fileContent) < 1 {
			t.Fatalf("Test error: file should have some content.")
		}
		vals := map[string]string{
			"find": fileContent[0],
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("GET", "/api/dir/"+dirname+"/search", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Matches []string `json:"matches"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		expected := fid.String()
		if len(results.Matches) != 1 || results.Matches[0] != expected {
			t.Fatalf("Expected %s to be the sole returned file. Received %s.", expected, results.Matches[0])
		}
	})
	t.Run("SearchDirNoResults", func(t *testing.T) {
		vals := map[string]string{
			"find": "thisshouldntappearinthefile",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("GET", "/api/dir/"+dirname+"/search", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Matches []string `json:"matches"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(results.Matches) != 0 {
			t.Fatalf("Expected no returned files. Received %+v.", results.Matches)
		}
	})
	t.Run("RemoveFile", func(t *testing.T) {
		vals := map[string]string{
			"id": fid.String(),
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("DELETE", "/api/dir/"+dirname+"/content", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that file has been removed
		req, _ = http.NewRequest("GET", "/api/dir/"+dirname, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Files []string `json:"files"`
			Name  string   `json:"name"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(results.Files) != 0 {
			t.Fatalf("Expected file to be removed from dir, leaving an empty dir. Receieved %+v", results.Files)
		}
	})
	t.Run("AddFile", func(t *testing.T) {
		fid := fid.String()
		vals := map[string]string{
			"id": fid,
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("POST", "/api/dir/"+dirname+"/content", bytes.NewReader(jsonbytes))
		req.Header.Add("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that file has been added
		req, _ = http.NewRequest("GET", "/api/dir/"+dirname, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Files []string `json:"files"`
			Name  string   `json:"name"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(results.Files) != 1 || results.Files[0] != fid {
			t.Fatalf("Expected %s to be sole file in dir. Receieved %+v", fid, results.Files)
		}
	})
	t.Run("DeleteDir", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/dir/"+dirname, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		//check that dir is empty
		req, _ = http.NewRequest("GET", "/api/dir/"+dirname, nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
		}
		var results struct {
			Files []string `json:"files"`
			Name  string   `json:"name"`
		}
		err := json.NewDecoder(res.Result().Body).Decode(&results)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		if len(results.Files) != 0 {
			t.Fatalf("Expected dir to be empty. Receieved %+v", results.Files)
		}
	})
}
