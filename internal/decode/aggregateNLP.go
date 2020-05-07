package decode

import (
	"sort"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/skyset"
)

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

var ignoreToBe = map[string]bool{
	"is":    true,
	"was":   true,
	"am":    true,
	"were":  true,
	"are":   true,
	"been":  true,
	"being": true,
	"be":    true,
}

func (nlp *nlpaggregate) add(phr []skyset.Phrase) {
	nlp.sentence++
	if nlp.data == nil {
		nlp.data = make(map[skyset.Synth]map[string]nlpaggregatedata)
	}
	for _, p := range phr {
		for _, t := range p.Tokens {
			if includepos[p.Synth] != nil &&
				includepos[p.Synth][t.Pos] &&
				((p.Synth != skyset.ACTION && p.Synth != skyset.PROCESS) || !ignoreToBe[strings.ToLower(t.Text)]) {
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

type nlpdatalistelement struct {
	word  string
	first uint
	count uint
}

type nlpdatalist []nlpdatalistelement

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

func (nlp *nlpaggregate) report() map[skyset.Synth]nlpdatalist {
	out := make(map[skyset.Synth]nlpdatalist)
	for syn, data := range nlp.data {
		var temp nlpdatalist
		for word, info := range data {
			temp = append(temp, nlpdatalistelement{
				word:  word,
				first: info.first,
				count: info.count,
			})
		}
		sort.Sort(temp)
		out[syn] = temp
	}
	return out
}

func (n nlpdatalist) tags(typ tag.Type) (tags []tag.Tag) {
	for i, data := range n {
		if i >= 50 {
			break
		}
		tags = append(tags, tag.Tag{
			Word: data.word,
			Type: typ,
			Data: tag.Data{
				typ: map[string]interface{}{
					"first":        data.first,
					"count":        data.count,
					"significance": i,
				},
			},
		})
	}
	return
}
