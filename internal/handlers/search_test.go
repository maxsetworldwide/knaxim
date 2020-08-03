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
