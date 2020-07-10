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

package util

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

var logBuf *bytes.Buffer

func resetLogger() {
	logBuf = new(bytes.Buffer)
	SetLogger(log.New(logBuf, "", log.LstdFlags))
}

// Two tests for functionality of the -v flag
func TestVerboseFlagOff(t *testing.T) {
	resetLogger()
	origLen := logBuf.Len()
	flag.Set("v", "false")
	Verbose("this is a test")
	if logBuf.Len() != origLen {
		t.Errorf("Fail: logs printed despite -v set to false")
	}
}

func TestVerboseFlagOn(t *testing.T) {
	resetLogger()
	origLen := logBuf.Len()
	flag.Set("v", "true")
	Verbose("this is a test")
	if logBuf.Len() <= origLen {
		t.Errorf("Fail: logs not printed despite -v set to true")
	}
}

// Test for inclusion of given message in the log message
func TestVerboseMessage(t *testing.T) {
	resetLogger()
	flag.Set("v", "true")
	testString := "test string"
	Verbose(testString)
	logString := logBuf.String()
	if strings.Index(logString, testString) < 0 {
		t.Errorf("Fail: given message not found in log output:\nmessage: %s,\nlog: %s", testString, logString)
	}
}

// Test functionality of Verbose's use of Printf and var args
func TestVerboseMessagePrintf(t *testing.T) {
	resetLogger()
	flag.Set("v", "true")
	addString := "test string to be placed within the formatting of the printf"
	addInt := 13679075492
	baseString := "this %d is the %s base string"
	expectedString := fmt.Sprintf(baseString, addInt, addString)
	Verbose(baseString, addInt, addString)
	logString := logBuf.String()
	if strings.Index(logString, expectedString) < 0 {
		t.Errorf("Fail: printf string (%s) not found in log output (%s)", expectedString, logString)
	}
}

// Test for inclusion of a request's information and printf message
func TestVerboseRequest(t *testing.T) {
	resetLogger()
	flag.Set("v", "true")
	testReq := httptest.NewRequest("POST", "https://knaxim.net/test/handle", nil)
	testBaseMessage := "test %s verbose %d string"
	addString := "this is a test string to be placed within the base"
	addInt := 7159402916
	expectedString := fmt.Sprintf(testBaseMessage, addString, addInt)
	wantedInfo := []string{
		testReq.RemoteAddr,
		testReq.Method,
		testReq.URL.Path,
		expectedString,
	}
	VerboseRequest(testReq, testBaseMessage, addString, addInt)
	logString := logBuf.String()
	for i, info := range wantedInfo {
		if len(info) == 0 {
			t.Errorf("Test error: passed string has length of 0 (wantedInfo[%d])", i)
		}
		if strings.Index(logString, info) < 0 {
			t.Errorf("Fail: info '%s' not found in log output: %s", info, logString)
		}
	}
}
