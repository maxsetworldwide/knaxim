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
