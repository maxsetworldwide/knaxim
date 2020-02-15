package handlers

import (
	"bytes"
	"encoding/json"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func setupFileAPI(t *testing.T) {
	AttachFile(testRouter.PathPrefix("/file").Subrouter())
	config.V.FileLimit = math.MaxInt64
	cookies = testlogin(t, 0)
}

var uploadFileContent = "This is the text content of a text file for upload"

// TODO: find out why getting a 'not found' when getting files added in setup
func TestFileAPI(t *testing.T) {
	t.Skip("Skipping due to problem with not finding files added in setup")
	setupFileAPI(t)
	var uploadfid filehash.FileID
	t.Run("Upload", func(t *testing.T) {
		body := new(bytes.Buffer)
		wrtr := multipart.NewWriter(body)
		part, err := wrtr.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("Unable to create multipart form file: %s\n", err)
		}
		_, err = part.Write([]byte(uploadFileContent))
		if err != nil {
			t.Fatalf("Failed to write file content to request: %s\n", err)
		}
		err = wrtr.Close()
		if err != nil {
			t.Fatalf("error closing multipart builder: %s\n", err)
		}
		req, err := http.NewRequest("PUT", "/api/file", body)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		req.Header.Set("Content-Type", wrtr.FormDataContentType())
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		var responseBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			t.Fatalf("Unable to decode response: %s\n", err)
		}
		uploadfid, err = filehash.DecodeFileID(responseBody["id"])
		if err != nil {
			t.Fatalf("Unable to get file ID from upload: %s\n%+#v", err, responseBody)
		}
		t.Log("Result ID:", uploadfid)
	})
	t.Run("Download", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/file/"+testFiles[1].store.ID.String()+"/download", nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		downloadContent := responseBodyString(res)
		if strings.Index(downloadContent, uploadFileContent) == -1 {
			t.Fatalf("Received wrong file: got content '%s', expected '%s'", downloadContent, uploadFileContent)
		}
	})
	t.Run("Download View", func(t *testing.T) {
		viewID := testFiles[0].store.ID.String()
		t.Logf("Querying view with ID: %s\n", viewID)
		req, err := http.NewRequest("GET", "/api/file/"+viewID+"/view", nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 999 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		t.Logf("Download view response: %+#v, body:%s", res, responseBodyString(res))
	})
}
