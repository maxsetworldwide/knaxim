package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func setupNLP(t *testing.T) {
	AttachNLP(testRouter.PathPrefix("/nlp").Subrouter())
	cookies = testlogin(t, 0, false)
}

type nlpresponse struct {
	File types.FileID `json:"fid"`
	Data []struct {
		Word  string `json:"word"`
		Count int    `json:"count"`
	} `json:"info"`
}

func TestNlp(t *testing.T) {
	setupNLP(t)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/nlp/file/%s/t/0/3", testFiles[0].file.GetID().String()), nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
	}
	var result nlpresponse
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("unable to decode response body: %s", responseBodyString(res))
	}
	if !(result.File.Equal(testFiles[0].file.GetID()) &&
		len(result.Data) == 3 &&
		result.Data[0].Word == "a" && result.Data[0].Count == 42) {
		t.Fatalf("incorrect result: %+#v", result)
	}
}
