package process

import "git.maxset.io/web/knaxim/pkg/skyset"

type nlpaggregate struct {
	sentence uint
	data     map[skyset.Synth]map[string]nlpaggregatedata
}

type nlpaggregatedata struct {
	count uint
	first uint
}

var includepos = map[skyset.Synth]map[skyset.PennPOS]bool{
	skyset.TOPIC: map[skyset.PennPOS]bool{
		skyset.NN:   true,
		skyset.NNP:  true,
		skyset.NNS:  true,
		skyset.NNPS: true,
	},
	skyset.ACTION: map[skyset.PennPOS]bool{
		skyset.VB:  true,
		skyset.VBD: true,
		skyset.VBG: true,
		skyset.VBN: true,
		skyset.VBP: true,
		skyset.VBZ: true,
	},
	skyset.RESOURCE: map[skyset.PennPOS]bool{
		skyset.NN:   true,
		skyset.NNS:  true,
		skyset.NNP:  true,
		skyset.NNPS: true,
	},
	skyset.PROCESS: map[skyset.PennPOS]bool{
		skyset.VB:  true,
		skyset.VBD: true,
		skyset.VBG: true,
		skyset.VBN: true,
		skyset.VBP: true,
		skyset.VBZ: true,
	},
}

func (nlp nlpaggregate) add(phr []skyset.Phrase) {
	nlp.sentence++
	if nlp.data == nil {
		nlp.data = make(map[skyset.Synth]map[string]nlpaggregatedata)
	}
	for _, p := range phr {
		for _, t := range p.Tokens {
			if includepos[p.Synth] != nil && includepos[p.Synth][t.Pos] {
				if nlp.data[p.Synth] == nil {
					nlp.data[p.Synth] = map[string]nlpaggregatedata{
						t.Text: nlpaggregatedata{
							first: nlp.sentence,
						},
					}
				}
				temp := nlp.data[p.Synth][t.Text]
				temp.count++
				if temp.first == 0 {
					temp.first = nlp.sentence
				}
				nlp.data[p.Synth][t.Text] = temp
			}
		}
	}
}

type nlpdatalist []struct {
	word  string
	first uint
	count uint
}

func (n nlpdatalist) Len() int {
	return len(n)
}

func (n nlpdatalist) Less(i, j int) bool {
	if n[i].count > n[j].count {
		return true
	}
	if n[i].count == n[j].count {
		return n[i].first < n[j].first
	}
	return false
}

func (n nlpdatalist) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
