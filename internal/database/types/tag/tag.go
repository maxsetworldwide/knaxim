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
	// SEARCH indicates that there are additional filter parameters within the Data to apply when searching tags
	SEARCH Type = (1 << 16) << iota
)

const (
	// USER indicates a custom tag created by a user, ie a folder
	USER Type = (1 << 24) << iota
	// DATE is to record a mapping of a date to a file made by a user.
	DATE
	// NAME is to record tags based on the name of a file
	NAME
)

// ALLTYPES is the compound of all possible tag types
const ALLTYPES = Type(math.MaxUint32)

// ALLSTORE are all the types of tags that are associated with a FileStore
const ALLSTORE = CONTENT | TOPIC | ACTION | RESOURCE | PROCESS

// ALLFILE are all the types of tags that are associated with a File
const ALLFILE = USER | DATE | NAME

// ALLSYNTH is the combination of TOPIC, ACTION, RESOURCE, and PROCESS
const ALLSYNTH = TOPIC | ACTION | RESOURCE | PROCESS

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
	case SEARCH:
		return "search"
	case USER:
		return "user"
	case DATE:
		return "date"
	case NAME:
		return "name"
	case ALLTYPES:
		return "alltypes"
	case ALLSYNTH:
		return "allsynth"
	case ALLSTORE:
		return "allstore"
	case ALLFILE:
		return "allfile"
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
	case "search":
		return SEARCH, nil
	case "user":
		return USER, nil
	case "date":
		return DATE, nil
	case "name":
		return NAME, nil
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
	if t.Word == "" {
		t.Word = oth.Word
	}
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

// FilterType returns a new instance of tag data that is a subset of keys defined by the type
func (d Data) FilterType(t Type) Data {
	out := make(Data)
	var dataPresent bool
	for k, v := range d {
		if k&t > 0 && k&^t == 0 {
			out[k] = v
			dataPresent = true
		}
	}
	if !dataPresent {
		return nil
	}
	return out
}

// StoreTag is a Tag tied to a FileStore
type StoreTag struct {
	Tag   `bson:",inline"`
	Store types.StoreID `bson:"store"`
}

// Update updates the StoreTag based on provided other StoreTag
func (st StoreTag) Update(oth StoreTag) StoreTag {
	id := st.Store
	if id.Equal(types.StoreID{}) {
		id = oth.Store
	}
	return StoreTag{
		Store: id,
		Tag:   st.Tag.Update(oth.Tag),
	}
}

// FileTag is a Tag tied to a File
type FileTag struct {
	Tag   `bson:",inline"`
	File  types.FileID  `bson:"file"`
	Owner types.OwnerID `bson:"owner"`
}

// Update creates a combination of the two tags
func (ft FileTag) Update(oth FileTag) FileTag {
	id := ft.File
	if id.Equal(types.FileID{}) {
		id = oth.File
	}
	owner := ft.Owner
	if owner.Equal(types.OwnerID{}) {
		owner = oth.Owner
	}
	return FileTag{
		File:  id,
		Owner: owner,
		Tag:   ft.Tag.Update(oth.Tag),
	}
}

// StoreTag builds StoreTag from FileTag
func (ft FileTag) StoreTag() StoreTag {
	return StoreTag{
		Tag: Tag{
			Word: ft.Tag.Word,
			Type: ft.Tag.Type & ALLSTORE,
			Data: ft.Tag.Data.FilterType(ALLSTORE),
		},
		Store: ft.File.StoreID,
	}
}

// Pure filters out elements that are represented with StoreTag
func (ft FileTag) Pure() FileTag {
	return FileTag{
		Tag: Tag{
			Word: ft.Tag.Word,
			Type: ft.Tag.Type & ALLFILE,
			Data: ft.Tag.Data.FilterType(ALLFILE),
		},
		File:  ft.File,
		Owner: ft.Owner,
	}
}
