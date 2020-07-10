/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
