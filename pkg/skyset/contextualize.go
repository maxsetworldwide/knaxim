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

package skyset

func contextualize(groups []grouping) []Phrase {
	out := make([]Phrase, 0, len(groups))
	history := []Synth{CONDITION}
	for _, g := range groups {
		out = append(out, Phrase{
			Synth:  nextSynth(g.group, &history),
			Tokens: g.tokens,
		})
	}
	return out
}

func nextSynth(grp group, history *[]Synth) (out Synth) {
	last := len(*history) - 1
	end := func() {
		if len(*history) > 1 {
			*history = (*history)[:last]
		} else {
			(*history)[0] = CONDITION
		}
	}

	switch grp {
	case NOUN:
		switch (*history)[last] {
		case CONDITION:
			out = TOPIC
			(*history)[last] = TOPIC
		case ACTION:
			out = RESOURCE
			(*history)[last] = RESOURCE
		default:
			end()
			out = nextSynth(grp, history)
		}
	case VERB:
		switch (*history)[last] {
		case CONDITION:
			fallthrough
		case TOPIC:
			out = ACTION
			(*history)[last] = ACTION
		case RESOURCE:
			out = PROCESS
			(*history)[last] = PROCESS
		default:
			end()
			out = nextSynth(grp, history)
		}
	case QUAL:
		out = CONDITION
		if (*history)[last] != CONDITION {
			*history = append(*history, CONDITION)
		}
	case CXN:
		out = CONNECTION
		switch (*history)[last] {
		case RESOURCE:
			(*history)[last] = ACTION
		case PROCESS:
			(*history)[last] = RESOURCE
		default:
			(*history)[last] = CONDITION
		}
	default:
		out = UNKNOWN
		*history = []Synth{CONDITION}
	}
	return
}
