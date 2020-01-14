package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestFile(t *testing.T) {
	cookies := testlogin(t, 1)
	LogBuffer(t)
	var uploadfid filehash.FileID
	t.Run("Upload", func(t *testing.T) {
		body := new(bytes.Buffer)
		wrtr := multipart.NewWriter(body)
		part, err := wrtr.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatal("Unable to create multipart form file: ", err)
		}
		_, err = part.Write([]byte("This is a test upload file."))
		if err != nil {
			t.Fatal("Failed to write to part: ", err)
		}
		err = wrtr.Close()
		if err != nil {
			t.Fatal("error closing multipart builder: ", err)
		}
		req, _ := http.NewRequest("PUT", server.URL+"/api/file", body)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Set("Content-Type", wrtr.FormDataContentType())
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatal("unable to decode response: ", err)
		}
		uploadfid, err = filehash.DecodeFileID(result["id"])
		if err != nil {
			t.Fatal("Unable to get file id", err, result)
		}
	})
	t.Run("UploadWeb", func(t *testing.T) {
		vals := map[string]string{
			"url": "https://en.wikipedia.org/wiki/Main_Page",
		}
		jsonbytes, _ := json.Marshal(vals)
		req, _ := http.NewRequest("PUT", server.URL+"/api/file/webpage", bytes.NewReader(jsonbytes))
		for _, c := range cookies {
			req.AddCookie(c)
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatal("unable to decode response: ", err)
		}
		_, err = filehash.DecodeFileID(result["id"])
		if err != nil {
			t.Fatal("Unable to get file id", err, result)
		}
	})
	t.Run("FileInfo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/file/"+uploadfid.String(), nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	t.Run("Slice", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/file/"+uploadfid.String()+"/slice/0/1", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	t.Run("Download", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/file/"+uploadfid.String()+"/download", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", server.URL+"/api/file/"+uploadfid.String(), nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
	})
}
