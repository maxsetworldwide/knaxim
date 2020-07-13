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
