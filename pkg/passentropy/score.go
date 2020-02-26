package passentropy

// This package estimates the amount of entropy within a particular
// password. It is meant to be used as a looser check of password
// complexity instead of requiring a specfic length and set of
// characters. This allows potential users to use a diverse set of
// methods to create safe passwords without being restricted to a
// particular set of rules

import (
	"math"
)

// These Constants are the expected scores for passwords of differing
// length and complexity. For example, Char6Cap1num1 is the expected
// score of a 6 character password with 1 capital and 1 number, while
// Char16Cap2num2oth2 is a 16 Characer, 2 capital, 2 numbers, and 2
// other characters (non-character and non-number, puctuation for
// example)
const (
	Char6Cap1num1 float64 = 4.82842712474619

	Char8             float64 = 4.756828460010884
	Char8Cap1         float64 = 5.303517070658851
	Char8Cap1num1     float64 = 5.833658625477635
	Char8Cap1num1oth1 float64 = 6.3437015248821105

	Char16             float64 = 8
	Char16Cap2         float64 = 8.919416985907818
	Char16Cap2num2     float64 = 9.81100525195611
	Char16Cap2num2oth2 float64 = 10.66879174342578
)

// Score estimates the complexity of a password
func Score(pass string) float64 {
	return newScoreBuilder(pass).calc()
}

type scoreBuilder struct {
	lowercase float64
	uppercase float64
	number    float64
	other     float64
}

func newScoreBuilder(word string) scoreBuilder {
	var sb scoreBuilder
	for _, c := range word {
		switch {
		case 'a' <= c && c <= 'z':
			sb.lowercase += 1.0
		case 'A' <= c && c <= 'Z':
			sb.uppercase += 1.0
		case '0' <= c && c <= '9':
			sb.number += 1.0
		default:
			sb.other += 1.0
		}
	}
	return sb
}

func (sb scoreBuilder) calc() float64 {
	return math.Sqrt(sb.lowercase*math.Sqrt(sb.lowercase)) + math.Sqrt(sb.uppercase*math.Sqrt(sb.uppercase)) + math.Sqrt(sb.number*math.Sqrt(sb.number)) + math.Sqrt(sb.other*math.Sqrt(sb.other))
}
