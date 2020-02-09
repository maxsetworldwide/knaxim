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
	for i, _ := range a {
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
