package srverror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

//TODO vanilla error test
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
}
