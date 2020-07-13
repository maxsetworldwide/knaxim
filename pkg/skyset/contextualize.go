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
