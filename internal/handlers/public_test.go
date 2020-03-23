package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupPublic(t *testing.T) {
	AttachPublic(testRouter.PathPrefix("/public").Subrouter())
	cookies = testlogin(t, 0, false)
}

type publicSearchTest struct {
	find       string
	expected   []int //indeces of publicFiles[]
	statusCode int
}

var publicSearchTests = []publicSearchTest{
	publicSearchTest{
		find:       "",
		expected:   []int{},
		statusCode: 400,
	},
	publicSearchTest{
		find:       "is",
		expected:   []int{0, 1},
		statusCode: 200,
	},
	publicSearchTest{
		find:       "fox",
		expected:   []int{2},
		statusCode: 200,
	},
	publicSearchTest{
		find:       "public has", //public searching multiple terms is AND
		expected:   []int{1},
		statusCode: 200,
	},
	publicSearchTest{
		find:       "noresults",
		expected:   []int{},
		statusCode: 200,
	},
	publicSearchTest{
		find:       "public brown",
		expected:   []int{},
		statusCode: 200,
	},
}

type publicSearchResult struct {
	Matched []struct {
		File struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"file"`
		Count int `json:"count"`
	} `json:"matched"`
}

func TestPublic(t *testing.T) {
	setupPublic(t)
	t.Run("PublicSearch", func(t *testing.T) {
		for i, test := range publicSearchTests {
			t.Run(fmt.Sprintf("PublicSearch-%d", i), func(t *testing.T) {
				vals := map[string]string{
					"find": test.find,
				}
				jsonbytes, _ := json.Marshal(vals)
				req, _ := http.NewRequest("GET", "/api/public/search", bytes.NewReader(jsonbytes))
				req.Header.Add("Content-Type", "application/json")
				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}
				res := httptest.NewRecorder()
				testRouter.ServeHTTP(res, req)
				if res.Code != test.statusCode {
					t.Fatalf("Expected status code %d: %+#v\nBody:%s", test.statusCode, res, responseBodyString(res))
				}
				var result publicSearchResult
				err := json.NewDecoder(res.Result().Body).Decode(&result)
				if err != nil {
					t.Fatalf("JSON Decode error:\n%s", err)
				}
				var expectedFiles []testFile
				for _, val := range test.expected {
					expectedFiles = append(expectedFiles, publicFiles[val])
				}
				if len(expectedFiles) != len(result.Matched) {
					t.Fatalf("Received incorrect files:\nExpected:%+v\nReceived:%+v", expectedFiles, result.Matched)
				}
				for _, file := range result.Matched {
					var matchedTestFile testFile
					found := false
					for i := 0; i < len(expectedFiles) && !found; i++ {
						curr := expectedFiles[i]
						if file.File.ID == curr.file.GetID().String() {
							matchedTestFile = curr
							found = true
						}
					}
					if !found {
						t.Fatalf("Received unexpected file. Expected%+#v\nReceived%+v", expectedFiles, file)
					}
					expectedSentences := strings.FieldsFunc(matchedTestFile.content, func(r rune) bool {
						return r == '.' || r == '!' || r == '?'
					})
					expectedCount := len(expectedSentences)
					t.Logf("expectedSentences: %+#v", expectedSentences)
					if file.Count != expectedCount {
						t.Fatalf("Returned count: %d, expected: %d", file.Count, expectedCount)
					}
				}
			})
		}
	})
}
