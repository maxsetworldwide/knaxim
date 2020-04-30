package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	AttachSearch(testRouter.PathPrefix("/search").Subrouter())
	cookies = testlogin(t, 0, false)
	query := `{
    "context": "%s",
    "match": "match"
  }`
	query = fmt.Sprintf(query, testUsers["users"][0]["id"])
	req, _ := http.NewRequest("POST", "/api/search/tags", strings.NewReader(query))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
	}
}
