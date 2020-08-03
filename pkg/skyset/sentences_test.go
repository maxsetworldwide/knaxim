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
