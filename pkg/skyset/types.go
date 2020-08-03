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

import (
	"fmt"
)

//group are phrase general types, intermediate to assigning synth value
type group uint8

type grouping struct {
	group  group
	tokens []Token
}

// Enumeration of group
const (
	UNK  group = iota
	NOUN group = iota
	VERB group = iota
	QUAL group = iota
	CXN  group = iota
)

//getGroup converts string representation to const
func getGroup(str string) group {
	switch str {
	case "Unknown":
		return UNK
	case "NOUN":
		return NOUN
	case "VERB":
		return VERB
	case "QUAL":
		return QUAL
	case "CXN":
		return CXN
	default:
		return UNK
	}
}

//String returns the string representation of Group
func (g group) String() string {
	switch g {
	case UNK:
		return "Unknown"
	case NOUN:
		return "NOUN"
	case VERB:
		return "VERB"
	case QUAL:
		return "QUAL"
	case CXN:
		return "CXN"
	default:
		return "Unknown"
	}
}

//MarshalJSON converts the Group to the string representation when marshalling into json
func (g group) MarshalJSON() ([]byte, error) {
	return []byte("\"" + g.String() + "\""), nil
}

//UnmarshalJSON decodes the Group type from json
func (g *group) UnmarshalJSON(input []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Unmarshal Group: %v", e)
		}
	}()
	*g = getGroup(string(input[1 : len(input)-1]))
	return nil
}

// Synth is a type of phrase, skyset category
type Synth int8

// Enumeration of Synth
const (
	UNKNOWN    Synth = iota
	TOPIC      Synth = iota
	ACTION     Synth = iota
	RESOURCE   Synth = iota
	PROCESS    Synth = iota
	CONDITION  Synth = iota
	CONNECTION Synth = iota
)

//GetSynth converts the string representation to the const value
func GetSynth(s string) Synth {
	switch s {
	case "TR":
		return TOPIC
	case "AS":
		return ACTION
	case "RP":
		return RESOURCE
	case "PC":
		return PROCESS
	case "CQ":
		return CONDITION
	case "CXN":
		return CONNECTION
	default:
		return UNKNOWN
	}
}

//String returns the string representation of Synth
func (s Synth) String() string {
	switch s {
	case UNKNOWN:
		return "Unknown"
	case TOPIC:
		return "TR"
	case ACTION:
		return "AS"
	case RESOURCE:
		return "RP"
	case PROCESS:
		return "PC"
	case CONDITION:
		return "CQ"
	case CONNECTION:
		return "CXN"
	default:
		return "Unknown"
	}
}

//MarshalJSON returns string representation in json form of Synth
func (s Synth) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

//UnmarshalJSON decodes string representation from json of Synth
func (s *Synth) UnmarshalJSON(js []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Unmarshal Synth: %v", e)
		}
	}()
	*s = GetSynth(string(js[1 : len(js)-1]))
	return nil
}

//PennPOS represents Penn Tree Part of Speech
type PennPOS uint8

// Enumeration of PennPOS
const (
	CC   PennPOS = iota
	CD   PennPOS = iota
	DT   PennPOS = iota
	EX   PennPOS = iota
	FW   PennPOS = iota
	IN   PennPOS = iota
	JJ   PennPOS = iota
	JJR  PennPOS = iota
	JJS  PennPOS = iota
	LS   PennPOS = iota
	MD   PennPOS = iota
	NN   PennPOS = iota
	NNS  PennPOS = iota
	NNP  PennPOS = iota
	NNPS PennPOS = iota
	PDT  PennPOS = iota
	POS  PennPOS = iota
	PRP  PennPOS = iota
	PRPS PennPOS = iota
	RB   PennPOS = iota
	RBR  PennPOS = iota
	RBS  PennPOS = iota
	RP   PennPOS = iota
	SYM  PennPOS = iota
	TO   PennPOS = iota
	UH   PennPOS = iota
	VB   PennPOS = iota
	VBD  PennPOS = iota
	VBG  PennPOS = iota
	VBN  PennPOS = iota
	VBP  PennPOS = iota
	VBZ  PennPOS = iota
	WDT  PennPOS = iota
	WP   PennPOS = iota
	WPS  PennPOS = iota
	WRB  PennPOS = iota
	PUNC PennPOS = iota
)

//GetPennPOS converts string representation to matching const
func GetPennPOS(s string) PennPOS {
	switch s {
	case "CC":
		return CC
	case "CD":
		return CD
	case "DT":
		return DT
	case "EX":
		return EX
	case "FW":
		return FW
	case "IN":
		return IN
	case "JJ":
		return JJ
	case "JJR":
		return JJR
	case "JJS":
		return JJS
	case "LS":
		return LS
	case "MD":
		return MD
	case "NN":
		return NN
	case "NNS":
		return NNS
	case "NNP":
		return NNP
	case "NNPS":
		return NNPS
	case "PDT":
		return PDT
	case "POS":
		return POS
	case "PRP":
		return PRP
	case "PRP$":
		return PRPS
	case "RB":
		return RB
	case "RBR":
		return RBR
	case "RBS":
		return RBS
	case "RP":
		return RP
	case "SYM":
		return SYM
	case "TO":
		return TO
	case "UH":
		return UH
	case "VB":
		return VB
	case "VBD":
		return VBD
	case "VBG":
		return VBG
	case "VBN":
		return VBN
	case "VBP":
		return VBP
	case "VBZ":
		return VBZ
	case "WDT":
		return WDT
	case "WP":
		return WP
	case "WP$":
		return WPS
	case "WRB":
		return WRB
	default:
		return PUNC
	}
}

//Byte converts const to single byte representation
func (p PennPOS) Byte() byte {
	switch p {
	case CC:
		return 'C'
	case CD:
		return 'E'
	case DT:
		return 'D'
	case EX:
		return 'X'
	case FW:
		return 'x'
	case IN:
		return 'i'
	case JJ:
		return 'J'
	case JJR:
		return 'K'
	case JJS:
		return 'j'
	case LS:
		return 'L'
	case MD:
		return 'M'
	case NN:
		return 'N'
	case NNS:
		return 'n'
	case NNP:
		return 'O'
	case NNPS:
		return 'o'
	case PDT:
		return 'T'
	case POS:
		return 'Q'
	case PRP:
		return 'P'
	case PRPS:
		return 'p'
	case RB:
		return 'R'
	case RBR:
		return 'U'
	case RBS:
		return 'r'
	case RP:
		return 'B'
	case SYM:
		return 'S'
	case TO:
		return '2'
	case UH:
		return 'A'
	case VB:
		return 'V'
	case VBD:
		return 'F'
	case VBG:
		return 'G'
	case VBN:
		return 'H'
	case VBP:
		return 'I'
	case VBZ:
		return 'Z'
	case WDT:
		return 'W'
	case WP:
		return 'Y'
	case WPS:
		return 'y'
	case WRB:
		return 'w'
	case PUNC:
		return '1'
	default:
		return 0
	}
}

//String returns string representation for PennPOS
func (p PennPOS) String() string {
	switch p {
	case CC:
		return "CC"
	case CD:
		return "CD"
	case DT:
		return "DT"
	case EX:
		return "EX"
	case FW:
		return "FW"
	case IN:
		return "IN"
	case JJ:
		return "JJ"
	case JJR:
		return "JJR"
	case JJS:
		return "JJS"
	case LS:
		return "LS"
	case MD:
		return "MD"
	case NN:
		return "NN"
	case NNS:
		return "NNS"
	case NNP:
		return "NNP"
	case NNPS:
		return "NNPS"
	case PDT:
		return "PDT"
	case POS:
		return "POS"
	case PRP:
		return "PRP"
	case PRPS:
		return "PRP$"
	case RB:
		return "RB"
	case RBR:
		return "RBR"
	case RBS:
		return "RBS"
	case RP:
		return "RP"
	case SYM:
		return "SYM"
	case TO:
		return "TO"
	case UH:
		return "UH"
	case VB:
		return "VB"
	case VBD:
		return "VBD"
	case VBG:
		return "VBG"
	case VBN:
		return "VBN"
	case VBP:
		return "VBP"
	case VBZ:
		return "VBZ"
	case WDT:
		return "WDT"
	case WP:
		return "WP"
	case WPS:
		return "WP$"
	case WRB:
		return "WRB"
	case PUNC:
		return "PUNC"
	default:
		return "XX"
	}
}

//MarshalJSON converts PennPOS into string representation in json encoding
func (p PennPOS) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.String() + "\""), nil
}

//UnmarshalJSON decodes PennPOS from json representation
func (p *PennPOS) UnmarshalJSON(js []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Unmarshal PennPOS: %v", e)
		}
	}()
	*p = GetPennPOS(string(js[1 : len(js)-1]))
	return nil
}

//Token represent a single word
type Token struct {
	Text string  `json:"word"`
	Pos  PennPOS `json:"pos"`
}

//Phrase is a series of words within one Synth
type Phrase struct {
	Synth  Synth   `json:"synth"`
	Tokens []Token `json:"tokens"`
}
