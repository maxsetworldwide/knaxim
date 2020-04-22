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
