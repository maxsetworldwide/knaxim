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

package srverror

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
 *  This test will write to the file system in the current directory. It will
 *  create a log folder specified below, then remove it. A partial UUID has been
 *  added to the directory name to help ensure that no unnecessary deletions are
 *  done.
 */
var testLogDir = "./testLogDir13d8bb899d9a"

func TestReport(t *testing.T) {
	t.Run("LogString", func(t *testing.T) {
		testURL := "/test/path"
		testCookies := []http.Cookie{
			http.Cookie{
				Name:  "testCookieName",
				Value: "testCookieValue",
			},
			http.Cookie{
				Name:  "testCookieNameTwo",
				Value: "testCookieValueTwo",
			},
		}
		testError := New(errors.New("wrapped error"), 555, "extra message one", "extra message two")
		testMethod := "POST"
		testBody := "test body for purposes of giving content length to request"
		testContentLength := len(testBody)
		testReq, _ := http.NewRequest(testMethod, testURL, strings.NewReader(testBody))
		var expectedCookieContent []string
		for _, c := range testCookies {
			testReq.AddCookie(&c)
			expectedCookieContent = append(expectedCookieContent, c.Value)
		}
		remoteAddr := "testRemoteAddr"
		testReq.RemoteAddr = remoteAddr
		testRes := httptest.NewRecorder()
		testResHeadKey := "Testresponseheaderkey"
		testResHeadVal := "testResponseHeaderValue"
		testRes.Header().Add(testResHeadKey, testResHeadVal)
		timestamp := time.Now().Format(time.RFC1123) // possible to not match actual
		// time during call. If that becomes a problem, might be better to regex for
		// existence of some timestamp. However, rerunning tests will generally work

		expectedContent := []string{
			timestamp,
			testURL,
			testMethod,
			strconv.Itoa(testError.Status()),
			strconv.Itoa(testContentLength),
			testError.Error(),
			testResHeadKey,
			testResHeadVal,
			remoteAddr,
		}
		expectedContent = append(expectedContent, expectedCookieContent...)

		result := LogString(testError, testReq, testRes)
		t.Logf("Result log string: \n%s\n", result)
		for _, expect := range expectedContent {
			if strings.Index(result, expect) == -1 {
				t.Errorf("Expected log string to contain '%s'", expect)
			}
		}
	})
	t.Run("File Write", func(t *testing.T) {
		LogPath = testLogDir
		currTime := time.Now()
		testLogPath := filepath.Join(testLogDir, currTime.Format("2006"), currTime.Format("01"), currTime.Format("02")+".log")
		testLogContent := "testLogMessage\n"
		defer func() {
			if err := os.RemoveAll(testLogDir); err != nil {
				t.Errorf(err.Error())
			}
		}()
		err := WriteToFile(testLogContent)
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf, err := ioutil.ReadFile(testLogPath)
		if err != nil {
			t.Fatalf("Unable to open output file: %s", err.Error())
		}
		if bytes.Compare(buf, []byte(testLogContent)) != 0 {
			t.Fatalf("Unexpected output. Expected:\n%s Received:\n%s", testLogContent, string(buf))
		}
		t.Run("File Write Append", func(t *testing.T) {
			appendedString := "testLogMessagePartTwo\n"
			expectedString := testLogContent + appendedString
			err := WriteToFile(appendedString)
			if err != nil {
				t.Fatalf(err.Error())
			}
			buf, err = ioutil.ReadFile(testLogPath)
			if err != nil {
				t.Fatalf("Unable to open output file: %s", err.Error())
			}
			if bytes.Compare(buf, []byte(expectedString)) != 0 {
				t.Fatalf("Unexpected output. Expected:\n%s Received:\n%s", expectedString, string(buf))
			}
		})
	})
}
