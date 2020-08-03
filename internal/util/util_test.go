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

import (
	"testing"
)

type SplitTest struct {
	in  []string
	out []string
}

var splitTests = []SplitTest{
	{[]string{}, []string{}},
	{[]string{""}, []string{}},
	{[]string{"a"}, []string{"a"}},
	{[]string{"a", "b"}, []string{"a", "b"}},
	{[]string{"b c"}, []string{"b", "c"}},
	{[]string{"a", "", "b"}, []string{"a", "b"}},
	{[]string{"a", ":b :"}, []string{"a", ":b", ":"}},
	{[]string{"a", "\"a"}, []string{"a", "a"}},
	{[]string{"\"a k p f\""}, []string{"a", "k", "p", "f"}},
	{[]string{"\"a", "k p", "\"e"}, []string{"a", "k", "p", "e"}},
	{[]string{"a     p"}, []string{"a", "p"}},
	{[]string{"     a   "}, []string{"a"}},
	{[]string{"a", "b 'c d'", "e"}, []string{"a", "b", "'c", "d'", "e"}},
	{[]string{"'a", "b'"}, []string{"'a", "b'"}},
	{[]string{"AB", "cd"}, []string{"AB", "cd"}},
	{[]string{"A b", "C d"}, []string{"A", "b", "C", "d"}},
}

// check for slice equality
// order matters in this check, however order may not be important in the
// actual SplitSearch method.
func stringSliceEq(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSplitSearch(t *testing.T) {
	for _, test := range splitTests {
		result := SplitSearch(test.in...)
		if !stringSliceEq(test.out, result) {
			t.Errorf("Fail: input: %v, got %v, expected %v", test.in, result, test.out)
		}
	}
}

type RegexTest struct {
	in  []string
	out string
}

var regexTests = []RegexTest{
	{[]string{}, "()"},
	{[]string{""}, "()"},
	{[]string{"a b"}, "(a)|(b)"},
	{[]string{"a    b"}, "(a)|(b)"},
	{[]string{"\"a b\""}, "(a b)"},
	{[]string{"\"a   b\""}, "(a   b)"},
	{[]string{"a", "b"}, "(a)|(b)"},
	{[]string{"a b", "c e d"}, "(a)|(b)|(c)|(e)|(d)"},
	{[]string{"a b", "\"c e d\""}, "(a)|(b)|(c e d)"},
	{[]string{"a b", "\"c e\" d"}, "(a)|(b)|(c e)|(d)"},
	{[]string{"a b", "\"c e\" \"d"}, "(a)|(b)|(c e)|(d)"},
}

func TestBuildSearchRegex(t *testing.T) {
	for _, test := range regexTests {
		result := BuildSearchRegex(test.in...)
		if result != test.out {
			t.Errorf("Fail: input: %v, got %v, expected %v", test.in, result, test.out)
		}
	}
}
