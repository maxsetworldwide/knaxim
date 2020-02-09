package tag

import (
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson"
)

type Type uint32

const (
	CONTENT Type = 1 << iota
	SYNTH_TOPIC
	SYNTH_ACTION
	SYNTH_RESOURCE
	SYNTH_PROCESS
)

const (
	USER Type = (1 << 24) << iota
)

const ALLTYPES = Type(math.MaxUint32)
const ALLSYNTH = Type(60)

func (t Type) String() string {
	switch t {
	case CONTENT:
		return "content"
	case SYNTH_TOPIC:
		return "topic"
	case SYNTH_ACTION:
		return "action"
	case SYNTH_PROCESS:
		return "process"
	case SYNTH_RESOURCE:
		return "resource"
	case USER:
		return "user"
	case ALLTYPES:
		return "alltypes"
	case ALLSYNTH:
		return "allsynth"
	default:
		return fmt.Sprintf("%X", uint32(t))
	}
}

func DecodeType(s string) (Type, error) {
	switch s {
	case "content":
		return CONTENT, nil
	case "topic":
		return SYNTH_TOPIC, nil
	case "action":
		return SYNTH_ACTION, nil
	case "process":
		return SYNTH_PROCESS, nil
	case "resource":
		return SYNTH_RESOURCE, nil
	case "user":
		return USER, nil
	case "alltypes":
		return ALLTYPES, nil
	case "allsynth":
		return ALLSYNTH, nil
	default:
		var out Type
		_, err := fmt.Sscanf(s, "%X", &out)
		if err != nil {
			return 0, err
		}
		return out, nil
	}
}

type Tag struct {
	Word string `bson:"word" json:"word"`
	Type Type   `bson:"type" json:"type"`
	Data Data   `bson:"data" json:"data"`
}

func (t Tag) Update(oth Tag) Tag {
	newt := Tag{
		Word: t.Word,
		Type: t.Type | oth.Type,
		Data: t.Data.Copy(),
	}
	for tk, mapping := range oth.Data {
		if newt.Data[tk] == nil {
			newt.Data[tk] = make(map[string]string)
		}
		for k, v := range mapping {
			newt.Data[tk][k] = v
		}
	}
	return newt
}

type Data map[Type]map[string]string

func (d Data) Copy() Data {
	if d == nil {
		return make(Data)
	}
	newd := make(Data)
	for tk, mapping := range d {
		newd[tk] = make(map[string]string)
		for k, v := range mapping {
			newd[tk][k] = v
		}
	}
	return newd
}

func (d Data) Contains(oth Data) bool {
	for typ, mapping := range oth {
		for k, v := range mapping {
			if d[typ] == nil || d[typ][k] != v {
				return false
			}
		}
	}
	return true
}

func (d Data) MarshalBSON() ([]byte, error) {
	form := make(map[string]map[string]string)
	for typ, fields := range d {
		form[typ.String()] = fields
	}
	return bson.Marshal(form)
}

func (d *Data) UnmarshalBSON(b []byte) error {
	*d = make(Data)
	var form map[string]map[string]string
	err := bson.Unmarshal(b, &form)
	if err != nil {
		return err
	}
	for tstr, fields := range form {
		t, err := DecodeType(tstr)
		if err != nil {
			return err
		}
		(*d)[t] = fields
	}
	return nil
}
