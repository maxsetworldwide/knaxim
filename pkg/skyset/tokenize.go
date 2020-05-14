package skyset

import (
	"regexp"
	"strings"

	"github.com/jdkato/prose/tag"
	t "github.com/jdkato/prose/tokenize"
)

var wtokenizer = t.NewTreebankWordTokenizer()
var ptagger = tag.NewPerceptronTagger()

func tokenize(s string) []Token {
	wl := wtokenizer.Tokenize(s)
	wl = splitconjunctions(wl)
	tokens := make([]Token, 0, len(wl))
	for _, tok := range ptagger.Tag(wl) {
		var pos PennPOS
		if !strings.ContainsAny(tok.Text, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			pos = PUNC
		} else {
			pos = GetPennPOS(tok.Tag)
		}
		tokens = append(tokens, Token{
			Text: tok.Text,
			Pos:  pos,
		})
	}
	return tokens
}

var puncdetect = regexp.MustCompile("[[:^alpha:]]+")

func splitconjunctions(wl []string) []string {
	nl := make([]string, 0, len(wl)+2)
	for _, word := range wl {
		start := 0
		splits := puncdetect.FindAllStringIndex(word, 0)
		for _, match := range splits {
			nl = append(nl, word[start:match[0]])
			start = match[0]
		}
		nl = append(nl, word[start:])
	}
	return nl
}