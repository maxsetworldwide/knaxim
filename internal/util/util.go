/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
