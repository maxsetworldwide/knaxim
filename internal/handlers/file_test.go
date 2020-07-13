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
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"sort"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/types"
)

var fileUserIdx = 2

func setupFileAPI(t *testing.T) {
	AttachFile(testRouter.PathPrefix("/file").Subrouter())
	config.V.FileLimit = math.MaxInt64
	cookies = testlogin(t, fileUserIdx, false)
}

var uploadFileName = "fileUploadTest.txt"
var uploadFileSentences = []string{
	"This is the text content of a text file for upload.",
	" We want it to have multiple sentences for the search feature.",
	" This is a third sentence!",
	" And this is the last one.",
	"\n", /* a newline is injected into uploaded files. Putting this here keeps
	the newline out of our last sentence. */
}
var uploadFileContent = strings.Join(uploadFileSentences, "")

type bound struct {
	start int
	end   int
}
type lineTest struct {
	fileidx    int // 0 => use the file uploaded during tests. !0 => use testFiles
	bounds     bound
	find       string
	expected   []int // indeces of uploadFileSentences
	statusCode int
}

var searchTests = []lineTest{
	lineTest{
		bounds: bound{
			start: 0,
			end:   3,
		},
		find:       "is",
		expected:   []int{0, 2},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "last",
		expected:   []int{3},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "nothing",
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   2,
		},
		find:       "last",
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 10,
			end:   20,
		},
		find:       "is",
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 3,
			end:   0,
		},
		find:       "is",
		expected:   []int{},
		statusCode: 400,
	},
	lineTest{
		bounds: bound{
			start: 3,
			end:   0,
		},
		find:       "is",
		expected:   []int{},
		statusCode: 400,
	},
	lineTest{
		fileidx: 1,
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "the",
		expected:   []int{},
		statusCode: 403,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "third last",
		expected:   []int{2, 3},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "\"third last\"",
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "\"multiple sentences\"",
		expected:   []int{1},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "\"is the\"",
		expected:   []int{0, 3},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   5,
		},
		find:       "",
		expected:   []int{0, 1, 2, 3},
		statusCode: 200,
	},
}

var sliceTests = []lineTest{
	lineTest{
		bounds: bound{
			start: 0,
			end:   3,
		},
		expected:   []int{0, 1, 2},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 0,
			end:   1,
		},
		expected:   []int{0},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 3,
			end:   4,
		},
		expected:   []int{3},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 2,
			end:   2,
		},
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 10,
			end:   20,
		},
		expected:   []int{},
		statusCode: 200,
	},
	lineTest{
		bounds: bound{
			start: 5,
			end:   0,
		},
		expected:   []int{},
		statusCode: 400,
	},
	lineTest{
		fileidx: 1,
		bounds: bound{
			start: 0,
			end:   5,
		},
		expected:   []int{},
		statusCode: 403,
	},
}

type lineResults struct {
	Lines []struct {
		ID struct {
			Hash  uint32 `json:"Hash"`
			Stamp uint16 `json:"Stamp"`
		} `json:"ID"`
		Position int      `json:"Position"`
		Content  []string `json:"Content"`
	} `json:"lines"`
	Size int `json:"size"`
}

func executeLineTest(t *testing.T, test lineTest, uploadfid types.FileID, search bool) {
	// if the test has a specified fileidx other than 0, testFiles[fileidx] will
	// take precedence over the given uploadfid param
	var fid types.FileID
	if test.fileidx != 0 {
		fid = testFiles[test.fileidx].file.GetID()
	} else {
		fid = uploadfid
	}
	t.Logf("fid: %s", fid.String())
	var req *http.Request
	var err error
	if search { // search bool used instead of checking struct field so we can search for empty strings
		params := map[string]string{
			"find": test.find,
		}
		jsonbytes, _ := json.Marshal(params)
		url := fmt.Sprintf("/api/file/%s/search/%d/%d", fid.String(), test.bounds.start, test.bounds.end)
		req, err = http.NewRequest("GET", url, bytes.NewReader(jsonbytes))
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		req.Header.Add("Content-Type", "application/json")
	} else {
		url := fmt.Sprintf("/api/file/%s/slice/%d/%d", fid.String(), test.bounds.start, test.bounds.end)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	if res.Code != test.statusCode {
		t.Fatalf("Unexpected status code. Expected %d\nResponse:%+#v\nBody:%s\n", test.statusCode, res, responseBodyString(res))
	}
	var jsonResponse lineResults
	err = json.NewDecoder(res.Result().Body).Decode(&jsonResponse)
	if err != nil {
		t.Fatalf("JSON Decode error:\n%s", err)
	}
	t.Logf("json response: %+v", jsonResponse)
	var allMatches []string
	var allPositions []int
	var receivedID types.StoreID
	for _, line := range jsonResponse.Lines {
		allMatches = append(allMatches, line.Content[0])
		allPositions = append(allPositions, line.Position)
		nextID := types.StoreID{
			Hash:  line.ID.Hash,
			Stamp: line.ID.Stamp,
		}
		if receivedID.ToNum() != 0 && !nextID.Equal(receivedID) {
			t.Fatalf("Expected one file ID, but got multiple: %s, %s", receivedID.String(), nextID.String())
		}
		receivedID = nextID
	}
	if receivedID.ToNum() != 0 && !receivedID.Equal(fid.StoreID) {
		t.Fatalf("Expected store ID to match given ID:\nReceieved:%s\nExpected:%s", receivedID.String(), fid.String())
	}
	if len(allMatches) != jsonResponse.Size {
		t.Fatalf("Returned size does not match number of received results.\nReceived %d, expected %d", jsonResponse.Size, len(allMatches))
	}
	sort.Ints(allPositions)
	sort.Ints(test.expected)
	if len(allPositions) != len(test.expected) {
		t.Fatalf("Received incorrect positions:\nReceived%+#v\nExpected%+#v", allPositions, test.expected)
	}
	for i := range allPositions {
		if allPositions[i] != test.expected[i] {
			t.Fatalf("Received incorrect positions:\nReceived%+#v\nExpected%+#v", allPositions, test.expected)
		}
	}
	var expectedMatches []string
	for _, idx := range test.expected {
		expectedMatches = append(expectedMatches, uploadFileSentences[idx])
	}
	t.Logf("Testing:\nExpected:%+#v\nReceived:%+#v", expectedMatches, allMatches)
	if len(allMatches) != len(expectedMatches) {
		t.Fatalf("Received wrong lines.\nExpected %+#v\nReceived %+#v", expectedMatches, allMatches)
	}
	for _, line := range allMatches {
		if !sliceContains(expectedMatches, line) {
			t.Fatalf("Received wrong lines.\nExpected %+#v\nReceived %+#v", expectedMatches, allMatches)
		}
	}
}

func TestFileAPI(t *testing.T) {
	setupFileAPI(t)
	var uploadfid types.FileID
	var uploadTime time.Time
	t.Run("Upload", func(t *testing.T) {
		body := new(bytes.Buffer)
		wrtr := multipart.NewWriter(body)
		mimeHead := make(textproto.MIMEHeader)
		mimeHead.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, uploadFileName))
		mimeHead.Set("Content-Type", "text/plain")
		part, err := wrtr.CreatePart(mimeHead)
		if err != nil {
			t.Fatalf("Unable to create multipart form file: %s\n", err)
		}
		_, err = part.Write([]byte(uploadFileContent))
		if err != nil {
			t.Fatalf("Failed to write file content to request: %s\n", err)
		}
		err = wrtr.Close()
		if err != nil {
			t.Fatalf("error closing multipart builder: %s\n", err)
		}
		req, err := http.NewRequest("PUT", "/api/file", body)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		req.Header.Set("Content-Type", wrtr.FormDataContentType())
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		var responseBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			t.Fatalf("Unable to decode response: %s\n", err)
		}
		uploadfid, err = types.DecodeFileID(responseBody["id"])
		uploadTime = time.Now()
		if err != nil {
			t.Fatalf("Unable to get file ID from upload: %s\n%+#v", err, responseBody)
		}
	})
	t.Run("FileInfo", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/file/"+uploadfid.String(), nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		var jsonResponse struct {
			File struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Date struct {
					Upload string `json:"upload"`
				} `json:"date"`
			} `json:"file"`
		}
		err = json.NewDecoder(res.Result().Body).Decode(&jsonResponse)
		if err != nil {
			t.Fatalf("JSON Decode error:\n%s", err)
		}
		// check that date is passable
		var returnedUpload time.Time
		returnedUpload, err = time.Parse(time.RFC3339Nano, jsonResponse.File.Date.Upload)
		if err != nil {
			t.Fatalf("Error parsing returned date: %s", err)
		}
		timeDifference := math.Abs(uploadTime.Sub(returnedUpload).Seconds())
		if timeDifference > 10.0 {
			t.Logf("Time difference in seconds: %f", timeDifference)
			t.Fatalf("Received upload time too different from time file was uploaded: logged time at %s, received time of %s", uploadTime, returnedUpload)
		}
		//check that id and name match
		if jsonResponse.File.ID != uploadfid.String() {
			t.Fatalf("Received incorrect file ID: received %s, expected %s", jsonResponse.File.ID, uploadfid.String())
		}
		if jsonResponse.File.Name != uploadFileName {
			t.Fatalf("Received incorrect file name: received %s, expected %s", jsonResponse.File.Name, uploadFileName)
		}
	})
	var uploadWebFid types.FileID
	t.Run("WebpageUpload", func(t *testing.T) {
		params := map[string]string{
			"url": "https://www.google.com/",
		}
		jsonbytes, _ := json.Marshal(params)
		req, err := http.NewRequest("PUT", "/api/file/webpage", bytes.NewReader(jsonbytes))
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		var responseBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			t.Fatalf("Unable to decode response: %s\n", err)
		}
		uploadWebFid, err = types.DecodeFileID(responseBody["id"])
		if err != nil {
			t.Fatalf("Unable to get webpage file ID from webpage upload: %s\n%+#v", err, responseBody)
		}
	})
	t.Run("Download", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/file/"+uploadfid.String()+"/download", nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		downloadContent := responseBodyString(res)
		if strings.Index(downloadContent, uploadFileContent) == -1 {
			t.Fatalf("Received wrong file: got content '%s', expected '%s'", downloadContent, uploadFileContent)
		}
	})
	t.Run("DownloadFileView", func(t *testing.T) {
		done := false
		timeout := config.V.MinFileTimeout.Seconds()
		totalWait := 0.0
		waitInterval := 1.5 //seconds
		var finishedResponse *httptest.ResponseRecorder
		for !done {
			req, err := http.NewRequest("GET", "/api/file/"+uploadfid.String()+"/view", nil)
			if err != nil {
				t.Fatal("Error creating http request: ", err)
			}
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			res := httptest.NewRecorder()
			testRouter.ServeHTTP(res, req)
			totalWait += waitInterval
			finishedResponse = res
			done = res.Code != 202 || totalWait > timeout
			t.Logf("Code: %d", res.Code)
			if !done {
				t.Logf("File still processing. Trying again in %v sec", waitInterval)
				time.Sleep(time.Duration(waitInterval * float64(time.Second)))
			}
		}
		if finishedResponse.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", finishedResponse, responseBodyString(finishedResponse))
		}
		if finishedResponse.Header().Get("Content-Type") != "application/pdf" {
			t.Fatalf("Expected application/pdf, got %s", finishedResponse.Header().Get("Content-Type"))
		}
		expectedFileName := uploadFileName[:strings.LastIndex(uploadFileName, ".")] + ".pdf"
		disposition := finishedResponse.Header().Get("Content-Disposition")
		if strings.Index(disposition, expectedFileName) == -1 {
			t.Fatalf("Expected attachment to be uploaded file name with pdf extension. Got %s, expected %s", disposition, expectedFileName)
		}
	})
	t.Run("DownloadWebPageView", func(t *testing.T) {
		done := false
		timeout := config.V.MinFileTimeout.Seconds()
		totalWait := 0.0
		waitInterval := 1.5 //seconds
		var finishedResponse *httptest.ResponseRecorder
		for !done {
			req, err := http.NewRequest("GET", "/api/file/"+uploadWebFid.String()+"/view", nil)
			if err != nil {
				t.Fatal("Error creating http request: ", err)
			}
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			res := httptest.NewRecorder()
			testRouter.ServeHTTP(res, req)
			totalWait += waitInterval
			finishedResponse = res
			done = res.Code != 202 || totalWait > timeout
			t.Logf("Code: %d", res.Code)
			if !done {
				t.Logf("File still processing. Trying again in %v sec", waitInterval)
				time.Sleep(time.Duration(waitInterval * float64(time.Second)))
			}
		}
		if finishedResponse.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", finishedResponse, responseBodyString(finishedResponse))
		}
		if finishedResponse.Header().Get("Content-Type") != "application/pdf" {
			t.Fatalf("Expected application/pdf, got %s", finishedResponse.Header().Get("Content-Type"))
		}
	})
	t.Run("FileSearch", func(t *testing.T) {
		for i, test := range searchTests {
			t.Run(fmt.Sprintf("FileSearch-%d", i), func(t *testing.T) {
				executeLineTest(t, test, uploadfid, true)
			})
		}
	})
	t.Run("Slices", func(t *testing.T) {
		for i, test := range sliceTests {
			t.Run(fmt.Sprintf("Slices-%d", i), func(t *testing.T) {
				executeLineTest(t, test, uploadfid, false)
			})
		}
	})
	t.Run("DeleteRecord", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/file/"+uploadfid.String(), nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		//check that file is gone
		req, err = http.NewRequest("GET", "/api/file/"+uploadfid.String(), nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 404 {
			t.Fatalf("Expected 404: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
	})
	t.Run("DeleteRecordPermissionDenied", func(t *testing.T) {
		fid := testFiles[1].file.GetID()
		req, err := http.NewRequest("DELETE", "/api/file/"+fid.String(), nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
		//check that file is still there
		cookies = testlogin(t, 1, false)
		req, err = http.NewRequest("GET", "/api/file/"+fid.String(), nil)
		if err != nil {
			t.Fatal("Error creating http request: ", err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res = httptest.NewRecorder()
		testRouter.ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("Non success status code: %+#v\nBody:%s\n", res, responseBodyString(res))
		}
	})
}
