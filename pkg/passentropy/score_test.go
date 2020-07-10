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
