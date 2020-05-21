package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func setupOwner(t *testing.T) {
	AttachOwner(testRouter.PathPrefix("/owner").Subrouter())
}

func TestOwner(t *testing.T) {
	setupOwner(t)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/owner/%s", testUsers["users"][0]["id"]), nil)
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
	}
	var result struct {
		ID   types.OwnerID `json:"id"`
		Name string        `json:"name"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("unable to decode response body (%s): %s", err.Error(), responseBodyString(res))
	}
	if result.Name != testUsers["users"][0]["name"] {
		t.Fatalf("incorrect result: %+#v", result)
	}
}
