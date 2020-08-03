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

package decode

import (
	"bytes"
	"fmt"

	"github.com/jdkato/prose/tokenize"
)

var stokenizer = tokenize.NewPunktSentenceTokenizer()

var maxSentLen = 5 << 10

// SentenceSplitter uses Punkt Sentence Tokenizer to scan for sentences in text
func SentenceSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(bytes.TrimSpace(data)) == 0 {
		return 0, nil, nil
	}
	sentences := stokenizer.Tokenize(string(data))
	if len(sentences) > 1 || (len(sentences) == 1 && atEOF) {
		offset := len(sentences[0])
		return offset, data[:offset], nil
	}
	if len(sentences) < 1 {
		return 0, nil, fmt.Errorf("Empty return from tokenizer, data = %s", string(data))
	}
	if len(data) > maxSentLen {
		for i := len(data) - 1; i >= 0; i-- {
			d := data[i]
			if d == '.' || d == '\n' || d == '?' || d == '!' {
				return i + 1, data[:i+1], nil
			}
		}
		for i := len(data) - 1; i >= 0; i-- {
			d := data[i]
			if d < 'A' || (d > 'Z' && d < 'a') || d > 'z' {
				return i + 1, data[:i+1], nil
			}
		}
	}
	return 0, nil, nil
}
