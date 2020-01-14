package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAcronym(t *testing.T) {
	cookies := testlogin(t, 0)
	LogBuffer(t)
	t.Run("Get=ab", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/acronym/ab", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		res, err := server.Client().Do(req)
		LogBuffer(t)
		if err != nil {
			t.Fatal("Client Error: ", err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("non success status code: %+#v", res)
		}
		var matches struct {
			Matched []string `json:"matched"`
		}
		err = json.NewDecoder(res.Body).Decode(&matches)
		if err != nil {
			t.Fatal("Unable to decode response ", err)
		}
		if len(matches.Matched) != len(testingacronyms["ab"]) {
			t.Fatal("incorrect return: ", matches.Matched, testingacronyms["ab"])
		}
	LOOP:
		for _, match := range matches.Matched {
			for _, original := range testingacronyms["ab"] {
				if match == original {
					continue LOOP
				}
			}
			t.Fatal("Incorrect return, ", match)
		}
	})
}
