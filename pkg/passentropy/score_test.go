package passentropy

import "testing"

func TestScore(t *testing.T) {
	cases := map[string]float64{
		"": 0,

		"A1aaaa": Char6Cap1num1,

		"aaaaaaaa": Char8,
		"Aaaaaaaa": Char8Cap1,
		"A1aaaaaa": Char8Cap1num1,
		"A1!aaaaa": Char8Cap1num1oth1,

		"aaaaaaaaaaaaaaaa": Char16,
		"AAaaaaaaaaaaaaaa": Char16Cap2,
		"AA11aaaaaaaaaaaa": Char16Cap2num2,
		"AA11!!aaaaaaaaaa": Char16Cap2num2oth2,
	}

	for pass, expect := range cases {
		s := Score(pass)
		if s != expect {
			t.Logf("incorrect score for `%s`: %g, expected: %g", pass, s, expect)
			t.Fail()
		}
	}
}
