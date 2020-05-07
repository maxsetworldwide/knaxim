package skyset

import (
	"encoding/json"
	"testing"
)

var sentences = []string{
	"This is the first sentence for testing.",
	"I have a lovely bunch of coconuts.",
	"How much wood can a wood chuck chuck, if a wood chuck could chuck wood.",
	"I love cheese.",
	"The cheese, milk, and meat are missing from the fridge",
}

// TestSentences initiates BuildPhrases to ensure that it doesn't panic
// correct results should be examined by hand
func TestSentences(t *testing.T) {
	for _, sent := range sentences {
		phrases := BuildPhrases(sent)
		j, err := json.MarshalIndent(phrases, "", "\t")
		if err != nil {
			t.Logf("unable to json encode: %s\n", err.Error())
		}
		t.Logf("%s => %s\n", sent, string(j))
	}
}
