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

package util

import "strings"

// SplitSearch divides the search strings into individual words for
// searching
func SplitSearch(search ...string) []string {
	out := make([]string, 0, len(search))
	for _, find := range search {
		find = strings.ReplaceAll(find, "\"", " ")
		for _, word := range strings.Split(find, " ") {
			if len(word) > 0 {
				out = append(out, word)
			}
		}
	}
	return out
}

// BuildSearchRegex generates a regular expression from the search terms
func BuildSearchRegex(search ...string) string {
	splits := make([]string, 0, len(search))
	for _, find := range search {
		for paren, phrase := range strings.Split(find, "\"") {
			if paren%2 == 0 {
				for _, word := range strings.Split(phrase, " ") {
					if len(word) > 0 {
						splits = append(splits, word)
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
