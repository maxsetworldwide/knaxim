package skyset

import (
	"regexp"
)

var rules = map[group]*regexp.Regexp{
	CXN:  regexp.MustCompile("(:?[C1A]+)"),
	NOUN: regexp.MustCompile("(:?[DWTpy]?(([NnOo]?Q)|([RUrw]*[JKjEBG]))*[NnOoB]+Q?)|(:?X)|(:?[PpYy]+[NnOo]*)|(:?Q)"),
	VERB: regexp.MustCompile("(:?M?[RUrw2]*([VFGHIZ]B?)+([RUrw2]|([VFGHIZ2]B?))*)|(:?M)"),
	QUAL: regexp.MustCompile("(:?[2RUrwJKjEiDWTpy]*[2RUrwJKjEiDWT])|(:?[2RUrwJKjEi]+[NnOo][2i]+)"),
}

func matchRule(seq []byte) group {
	for g, regex := range rules {
		indexes := regex.FindIndex(seq)
		if indexes != nil && indexes[0] == 0 && indexes[1] >= len(seq) {
			return g
		}
	}
	return UNK
}

func assemble(tokens []Token) []grouping {
	sequence := tokens2bytes(tokens)
	cxnBreaks := rules[CXN].FindAllIndex(sequence, -1)
	cxnBreaks = append(cxnBreaks, []int{len(sequence), len(sequence)})
	begin := 0
	var cuts []cut
	for _, cxnbreak := range cxnBreaks {
		if cxnbreak[0]-begin > 0 {
			cuts = append(cuts, shiftcuts(cutSequence(sequence[begin:cxnbreak[0]]), begin)...)
		}
		if cxnbreak[1]-cxnbreak[0] > 0 {
			cuts = append(cuts, cut{
				group: CXN,
				start: cxnbreak[0],
				end:   cxnbreak[1],
			})
		}
		begin = cxnbreak[1]
	}
	cuts = combineLists(cuts)
	out := make([]grouping, 0, len(cuts))
	for _, c := range cuts {
		out = append(out, grouping{
			group:  c.group,
			tokens: tokens[c.start:c.end],
		})
	}
	return out
}

func tokens2bytes(tokens []Token) []byte {
	out := make([]byte, 0, len(tokens))
	for _, t := range tokens {
		out = append(out, t.Pos.Byte())
	}
	return out
}

type cut struct {
	group group
	start int
	end   int
}

func cutSequence(seq []byte) []cut {
	solutions := make([][]cut, len(seq)+1)
	scores := make([]float64, len(seq)+1)
	for i := range scores {
		scores[i] = float64(2 * len(seq))
	}

	solutions[0] = []cut{}
	scores[0] = 0
	for i := 1; i < len(solutions); i++ {
		for j := 0; j < i; j++ {
			if solutions[j] != nil && scores[j]+1 < scores[i] {
				if seqGroup := matchRule(seq[j:i]); seqGroup != UNK {
					scr := scores[j] + 1
					if seqGroup == QUAL {
						scr += float64(i-j) / float64(len(scores))
					}
					if scr < scores[i] {
						solutions[i] = make([]cut, 0, len(solutions[j])+1)
						solutions[i] = append(solutions[i], solutions[j]...)
						solutions[i] = append(
							solutions[i],
							cut{
								group: seqGroup,
								start: j,
								end:   i,
							},
						)
						scores[i] = scr
					}
				}
			}
		}
	}
	solution := solutions[len(seq)]
	if solution == nil {
		return []cut{cut{
			group: UNK,
			start: 0,
			end:   len(seq),
		}}
	}
	return solution
}

func shiftcuts(cuts []cut, n int) []cut {
	var out []cut
	for _, c := range cuts {
		c.start += n
		c.end += n
		out = append(out, c)
	}
	return out
}

func combineLists(cuts []cut) []cut {
	combine := func(cuts []cut, listtype group) []cut {
		for i := 0; i < len(cuts)-4; i++ {
			if cuts[i].group == listtype && cuts[i+1].group == CXN && cuts[i+2].group == listtype && cuts[i+3].group == CXN && cuts[i+4].group == listtype {
				lastCut := i + 4
				for lastCut < len(cuts)-2 {
					if cuts[lastCut+1].group == CXN && cuts[lastCut+2].group == listtype {
						lastCut = lastCut + 2
					} else {
						break
					}
				}
				newcuts := make([]cut, 0, len(cuts)+i-lastCut)
				newcuts = append(newcuts, cuts[:i]...)
				newcuts = append(newcuts, cut{
					group: listtype,
					start: cuts[i].start,
					end:   cuts[lastCut].end,
				})
				newcuts = append(newcuts, cuts[lastCut+1:]...)
				cuts = newcuts
			}
		}
		return cuts
	}
	cuts = combine(cuts, NOUN)
	cuts = combine(cuts, VERB)
	cuts = combine(cuts, QUAL)
	return cuts
}
