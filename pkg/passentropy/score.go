package passentropy

import (
	"math"
)

const Char6Cap1num1 float64 = 4.82842712474619

const Char8 float64 = 4.756828460010884
const Char8Cap1 float64 = 5.303517070658851
const Char8Cap1num1 float64 = 5.833658625477635
const Char8Cap1num1oth1 float64 = 6.3437015248821105

const Char16 float64 = 8
const Char16Cap2 float64 = 8.919416985907818
const Char16Cap2num2 float64 = 9.81100525195611
const Char16Cap2num2oth2 float64 = 10.66879174342578

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
