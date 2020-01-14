package main

import "strings"

func splitSearch(search ...string) []string {
	out := make([]string, 0, len(search))
	for _, find := range search {
		find = strings.ReplaceAll(find, "\"", " ")
		for _, nospace := range strings.Split(find, " ") {
			for _, word := range strings.Split(nospace, "%20") {
				if len(word) > 0 {
					out = append(out, word)
				}
			}
		}
	}
	return out
}

func buildSearchRegex(search ...string) string {
	splits := make([]string, 0, len(search))
	for _, find := range search {
		for paren, phrase := range strings.Split(find, "\"") {
			if paren%2 == 0 {
				for _, nospace := range strings.Split(phrase, " ") {
					for _, word := range strings.Split(nospace, "%20") {
						if len(word) > 0 {
							splits = append(splits, word)
						}
					}
				}
			} else {
				if len(phrase) > 0 {
					splits = append(splits, phrase)
				}
			}
		}
	}
	return "(" + strings.Join(splits, ")|(") + ")"
}
