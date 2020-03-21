package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.maxset.io/web/knaxim/internal/config"
)

type acronymEntry struct {
	key string
	val string
}

var acronymEntries = []acronymEntry{
	{"af", "Air Force"},
	{"ae", "acronym entry"},
	{"entry", "outry"},
	{"af", "Air Force Part 2: Electric Boogaloo"},
}

func setupAcronym(t *testing.T) {
	t.Helper()
	AttachAcronym(testRouter.PathPrefix("/acronym").Subrouter())
	cookies = testlogin(t, 0, false)
	ab := config.DB.Acronym(nil)
	for _, acr := range acronymEntries {
		if err := ab.Put(acr.key, acr.val); err != nil {
			t.Errorf("Acronym database put error: %s", err)
		}
	}
	ab.Close(nil)
}

func sendAcronymRequest(t *testing.T, query string) acronymResult {
	url := "/api/acronym/" + query
	request, _ := http.NewRequest("GET", url, nil)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Expected acronym response code to be 200 or 404, instead is %d.", response.Code)
	}
	var matches acronymResult
	err := json.NewDecoder(response.Result().Body).Decode(&matches)
	if err != nil {
		t.Errorf("JSON Decode error, possibly received empty results despite response code being 200:\n%s", err)
	}
	return matches
}

type acronymResult struct {
	Matched []string `json:"matched"`
}

func TestAcronymAPI(t *testing.T) {
	setupAcronym(t)
	for _, acr := range acronymEntries {
		result := sendAcronymRequest(t, acr.key)
		if !sliceContains(result.Matched, acr.val) {
			t.Errorf("Fail: expected acronym was not returned:\nexpected %s\nreceived %v", acr.val, result.Matched)
		}
	}
	// test for non existent values
	nonExistentResult := sendAcronymRequest(t, "this shouldn't exist")
	if len(nonExistentResult.Matched) != 0 {
		t.Errorf("Fail: Expected no results from non-existent query, instead got %v", nonExistentResult.Matched)
	}
}
