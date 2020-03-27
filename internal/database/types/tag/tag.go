package tag

import (
	"fmt"
	"math"

	"git.maxset.io/web/knaxim/internal/database/types"
	"go.mongodb.org/mongo-driver/bson"
)

// Type is a byte flag for different kinds of tags
type Type uint32

// Types of tags, bitwise or them together to make compound tags
const (
	CONTENT Type = 1 << iota
	TOPIC
	ACTION
	RESOURCE
	PROCESS
)

const (
	// USER indicates a custom tag created by a user, ie a folder
	USER Type = (1 << 24) << iota
)

// ALLTYPES is the compound of all possible tag types
const ALLTYPES = Type(math.MaxUint32)

// ALLSYNTH is the combination of TOPIC, ACTION, RESOURCE, and PROCESS
const ALLSYNTH = Type(60)

func (t Type) String() string {
	switch t {
	case CONTENT:
		return "content"
	case TOPIC:
		return "topic"
	case ACTION:
		return "action"
	case PROCESS:
		return "process"
	case RESOURCE:
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

// DecodeType converts String representation back to Type
func DecodeType(s string) (Type, error) {
	switch s {
	case "content":
		return CONTENT, nil
	case "topic":
		return TOPIC, nil
	case "action":
		return ACTION, nil
	case "process":
		return PROCESS, nil
	case "resource":
		return RESOURCE, nil
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

// Tag is a combination of a word and a type to attach to a file or store
// Data is to store additional data about a particular tag
type Tag struct {
	Word string `bson:"word" json:"word"`
	Type Type   `bson:"type" json:"type"`
	Data Data   `bson:"data" json:"data"`
}

// Update adds or replaces data values in the passed tag
func (t Tag) Update(oth Tag) Tag {
	newt := Tag{
		Word: t.Word,
		Type: t.Type | oth.Type,
		Data: t.Data.Copy(),
	}
	for tk, mapping := range oth.Data {
		if newt.Data[tk] == nil {
			newt.Data[tk] = make(map[string]interface{})
		}
		for k, v := range mapping {
			newt.Data[tk][k] = v
		}
	}
	return newt
}

// Data maps a tag type to an abitrary collection of string to string mappings
type Data map[Type]map[string]interface{}

// Copy generates a new Data object with the same values as the original
func (d Data) Copy() Data {
	if d == nil {
		return make(Data)
	}
	newd := make(Data)
	for tk, mapping := range d {
		newd[tk] = make(map[string]interface{})
		for k, v := range mapping {
			newd[tk][k] = v
		}
	}
	return newd
}

// Contains returns true if all data fields in the provided data match the fields in the original
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

// MarshalBSON converts Data into a bson representation
func (d Data) MarshalBSON() ([]byte, error) {
	form := make(map[string]map[string]interface{})
	for typ, fields := range d {
		form[typ.String()] = fields
	}
	return bson.Marshal(form)
}

// UnmarshalBSON converts bson representation back into Data
func (d *Data) UnmarshalBSON(b []byte) error {
	*d = make(Data)
	var form map[string]map[string]interface{}
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

// StoreTag is a Tag of
type StoreTag struct {
	Tag   `bson:",inline"`
	Store types.StoreID `bson:"sid"`
}

type FileTag struct {
	Tag   `bson:",inline"`
	File  types.FileID  `bson:"file"`
	Owner types.OwnerID `bson:"owner"`
}
