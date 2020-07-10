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

package srverror

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestSrvError(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		msgs := []string{
			"message one",
			"message 2",
		}
		statusCode := 555
		srverr := Basic(statusCode, msgs...)
		t.Run("Message", func(t *testing.T) {
			resultMessage := srverr.Error()
			for _, msg := range msgs {
				if strings.Index(resultMessage, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, resultMessage)
				}
			}
		})
		t.Run("Extend", func(t *testing.T) {
			extensions := []string{
				"extension 1",
				"second extension",
			}
			allMsgs := append(msgs, extensions...)
			extendedError := srverr.Extend(extensions...)
			resultMessage := extendedError.Error()
			for _, msg := range allMsgs {

				if strings.Index(resultMessage, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, resultMessage)
				}
			}
		})
		t.Run("Status", func(t *testing.T) {
			if srverr.Status() != statusCode {
				t.Fatalf("Provided status code was not returned back. Expected %d, received %d", statusCode, srverr.Status())
			}
		})
		t.Run("ServeHTTPDebugOn", func(t *testing.T) {
			testServeHTTP(t, true, srverr, statusCode, msgs)
		})
		t.Run("ServeHTTPDebugOff", func(t *testing.T) {
			testServeHTTP(t, false, srverr, statusCode, msgs)
		})
	})
	t.Run("New", func(t *testing.T) {
		wrappedErrMsg := "this is the message in the wrapped error"
		wrappedErr := errors.New(wrappedErrMsg)
		statusCode := 525
		msgs := []string{
			"new message one",
			"second message for new",
		}
		srverr := New(wrappedErr, statusCode, msgs...)
		t.Run("Message", func(t *testing.T) {
			resultMessage := srverr.Error()
			for _, msg := range msgs {
				if strings.Index(resultMessage, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, resultMessage)
				}
			}
		})
		t.Run("Extend", func(t *testing.T) {
			extensions := []string{
				"extension for new one",
				"second extension messsage",
			}
			allMsgs := append(msgs, extensions...)
			extendedError := srverr.Extend(extensions...)
			resultMessage := extendedError.Error()
			for _, msg := range allMsgs {
				if strings.Index(resultMessage, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, resultMessage)
				}
			}
		})
		t.Run("Status", func(t *testing.T) {
			if srverr.Status() != statusCode {
				t.Fatalf("Provided status code was not returned back. Expected %d, received %d", statusCode, srverr.Status())
			}
		})
		t.Run("Unwrap", func(t *testing.T) {
			unwrapped := srverr.Unwrap()
			if unwrapped != wrappedErr {
				t.Fatalf("Expected unwrapped error to be same error object that was initially wrapped.\nExpected %+#v\nReceived %+#v", wrappedErr, unwrapped)
			}
		})
		t.Run("ServeHTTPDebugOn", func(t *testing.T) {
			testServeHTTP(t, true, srverr, statusCode, msgs)
		})
		t.Run("ServeHTTPDebugOff", func(t *testing.T) {
			testServeHTTP(t, false, srverr, statusCode, msgs)
		})
	})
	/*
	 * Nil error is being tested, but outside packages will not be able to
	 * replicate the behavior, since Error is the exported type and is an
	 * interface and these functions cannot be called with a nil interface.
	 * They can, however, be called with a nil srverr pointer.
	 */
	t.Run("NilSrverror", func(t *testing.T) {
		var errObj *srverr
		t.Run("Message", func(t *testing.T) {
			msg := errObj.Error()
			if msg != "" {
				t.Fatalf("Expected empty message from Error(). Received '%s'", msg)
			}
		})
		t.Run("Extend", func(t *testing.T) {
			expectedCode := 500
			extensions := []string{
				"nil error extend message",
				"extension for nil part two",
			}
			extended := errObj.Extend(extensions...)
			errMsg := extended.Error()
			if extended.Status() != expectedCode {
				t.Fatalf("Expected internal error code %d. Got %d", expectedCode, extended.Status())
			}
			for _, msg := range extensions {
				if strings.Index(errMsg, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, errMsg)
				}
			}
		})
		t.Run("Unwrap", func(t *testing.T) {
			unwrapped := errObj.Unwrap()
			if unwrapped != nil {
				t.Fatalf("Expected unwrapped nil error to be nil. Received %+v", unwrapped)
			}
		})
		t.Run("ServeHTTP", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			errObj.ServeHTTP(recorder, nil)
			expectedCode := 200
			t.Run("Status Code", func(t *testing.T) {
				if recorder.Result().StatusCode != expectedCode {
					t.Fatalf("Expected status code %d from nil error. Received %d.", expectedCode, recorder.Result().StatusCode)
				}
			})
		})
	})
	t.Run("Nested Error", func(t *testing.T) {
		innerMsgs := []string{
			"inner one",
			"second inner",
		}
		innerStatus := 987
		innerErr := Basic(innerStatus, innerMsgs...)
		outerMsgs := []string{
			"outer part 1",
			"this one is the second on the outside",
		}
		outerStatus := 654
		outerErr := New(innerErr, outerStatus, outerMsgs...)
		t.Run("Message", func(t *testing.T) {
			allMsgs := append(innerMsgs, outerMsgs...)
			resultMessage := outerErr.Error()
			for _, msg := range allMsgs {
				if strings.Index(resultMessage, msg) == -1 {
					t.Fatalf("Expected '%s' to appear in Error() string. Received Error() string '%s'", msg, resultMessage)
				}
			}
			if strings.Index(resultMessage, strconv.Itoa(innerStatus)) == -1 {
				t.Fatalf("Expected inner status code to be mentioned in wrapping error messages. Expected %d, Received %s", innerStatus, resultMessage)
			}
		})
		t.Run("Status", func(t *testing.T) {
			if outerErr.Status() != outerStatus {
				t.Fatalf("Provided status code was not returned back. Expected %d, received %d", outerStatus, outerErr.Status())
			}
		})
		t.Run("Unwrap", func(t *testing.T) {
			unwrapped := outerErr.Unwrap()
			if unwrapped != innerErr.Unwrap() {
				t.Fatalf("Expected unwrapped error to be the wrapped error of the original inner error.\nReceived %+v\nExpected %+v", unwrapped, innerErr)
			}
		})
	})
	t.Run("Nil New", func(t *testing.T) {
		status := 1000
		msgs := []string{
			"these messages",
			"should not appear",
			"in the error",
		}
		newErr := New(nil, status, msgs...)
		if newErr != nil {
			t.Fatalf("Expected nil Error to be returned when providing nil error.\nReceived %+v", newErr)
		}
	})
}

func testServeHTTP(t *testing.T, debug bool, srverr Error, statusCode int, msgs []string) {
	DEBUG = debug
	recorder := httptest.NewRecorder()
	srverr.ServeHTTP(recorder, nil)
	var body struct {
		Message string `json:"message"`
	}
	err := json.NewDecoder(recorder.Result().Body).Decode(&body)
	if err != nil {
		t.Fatalf("Error reading http response: %s", err)
	}
	t.Run("Status Code", func(t *testing.T) {
		if recorder.Result().StatusCode != statusCode {
			t.Fatalf("Received incorrect status code in http response:\nExpected %d\nReceived%+#v", statusCode, recorder)
		}
	})
	t.Run("Message", func(t *testing.T) {
		if debug {
			for _, msg := range msgs {
				if strings.Index(body.Message, msg) == -1 {
					t.Fatalf("Message not found in http response body.\nExpected %s\nReceived %s", msg, body.Message)
				}
			}
		} else {
			if strings.Index(body.Message, msgs[0]) == -1 {
				t.Fatalf("Message not found in http response body.\nExpected %s\nReceived %s", msgs[0], body.Message)
			}
		}
	})
}
