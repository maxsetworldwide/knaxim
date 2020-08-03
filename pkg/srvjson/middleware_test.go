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

package srvjson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHandler struct {
	t *testing.T
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jw, ok := w.(*ResponseWriter)
	if !ok {
		mh.t.Fatal("middleware did not replace writer with json Writer")
		return
	}
	jw.Set("Test", "Hello")
	jw.Write([]byte("World"))
	jw.Write([]byte(", and again."))
	jw.WriteHeader(5)
	sb := new(strings.Builder)
	sb.WriteString(jw.data["message"].(string))
	jw.Set("message", sb)
	jw.Write([]byte(" message is now a string builder."))
}

func TestMiddleware(t *testing.T) {
	handler := JSONResponse(&mockHandler{
		t: t,
	})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, nil)
	if w.Code != 5 {
		t.Errorf("incorrect Status Code: %d", w.Code)
	}
	var result struct {
		Test    string
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("unable to decode response: %s", err)
	}
	if result.Test != "Hello" {
		t.Errorf("incorrect test field: %s", result.Test)
	}
	if result.Message != "World, and again. message is now a string builder." {
		t.Errorf("incorrect message: %s", result.Message)
	}
}
