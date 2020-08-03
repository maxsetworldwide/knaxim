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
	t.Logf("lookup id")
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/owner/id/%s", testUsers["users"][0]["id"]), nil)
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
	}
	var result struct {
		ID   types.OwnerID `json:"id"`
		Name string        `json:"name"`
		Type string        `json:"type"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("unable to decode response body (%s): %s", err.Error(), responseBodyString(res))
	}
	if result.Name != testUsers["users"][0]["name"] {
		t.Fatalf("incorrect result: %+#v", result)
	}
	t.Logf("lookup name")
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/owner/name/%s", testUsers["users"][0]["name"]), nil)
	res = httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("non success status code: %+#v\nBody:%s", res, responseBodyString(res))
	}
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("unable to decode response body (%s): %s", err.Error(), responseBodyString(res))
	}
	if result.ID.String() != testUsers["users"][0]["id"] {
		t.Fatalf("incorrect result: %+#v", result)
	}
}
