package handlers

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	cookies := testlogin(t, 0)
	LogBuffer(t)
	t.Run("UserContext", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user/search?find=file", nil)
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
		returned := new(strings.Builder)
		_, err = io.Copy(returned, res.Body)
		if err != nil {
			t.Fatal("unable to read body: ", err)
		}
		r := returned.String()
		t.Logf("returned=%s", r)
		// var matched SearchResponse
		// if err = json.NewDecoder(strings.NewReader(r)).Decode(&matched); err != nil {
		// 	t.Fatal("unable to decode responce; ", err)
		// }
		// if len(matched.Files) != 1 {
		// 	t.Fatalf("incorrect returned fileid: %v", matched)
		// }
	})
	t.Run("SearchFile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/file/"+testingfiles[0].file.GetID().String()+"/search/0/1?find=file", nil)
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
