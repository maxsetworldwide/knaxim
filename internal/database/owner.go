package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/brand"
	"go.mongodb.org/mongo-driver/bson"
)

// OwnerID used to uniquely identify owners
type OwnerID struct {
	Type        byte    `bson:"type"`
	UserDefined [3]byte `bson:"ud"`
	Stamp       []byte  `bson:"stamp,omitempty"`
}

func name2userdefined(name string) [3]byte {
	var out [3]byte
	if len(name) >= 3 {
		for i, b := range []byte(name[:3]) {
			out[i] = b
		}
	} else {
		nameprefix := []byte(name)
		for len(nameprefix) < 3 {
			nameprefix = append(nameprefix, '*')
		}
		for i, b := range nameprefix {
			out[i] = b
		}
	}
	return out
}

// func newstamp() []byte {
// 	buf := new(bytes.Buffer)
// 	binary.Write(buf, binary.BigEndian, uint16(time.Now().UnixNano()))
// 	return buf.Bytes()
// }

func (oid OwnerID) bytes() []byte {
	b := make([]byte, 0, 4+len(oid.Stamp))
	b = append(b, oid.Type)
	b = append(b, oid.UserDefined[:]...)
	b = append(b, oid.Stamp...)
	return b
}

// String returns string representation of OwnerID
func (oid OwnerID) String() string {
	out := new(strings.Builder)
	enc := base64.NewEncoder(base64.RawURLEncoding, out)
	enc.Write(oid.bytes())
	enc.Close()
	return out.String()
}

// DecodeObjectIDString inverse operation of String()
func DecodeObjectIDString(s string) (oid OwnerID, err error) {
	defer func() {
		if r := recover(); r != nil {
			oid = OwnerID{}
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("Unable to Decode Object ID: %v", r)
			}
		}
	}()
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return OwnerID{}, err
	}
	if len(b) < 4 {
		return OwnerID{}, fmt.Errorf("Unable to Decode Object ID, too short")
	}
	var out OwnerID
	out.Type = b[0]
	for i, b := range b[1:4] {
		out.UserDefined[i] = b
	}
	out.Stamp = b[4:]
	return out, nil
}

// Mutate creates new OwnerID of the same type
func (oid OwnerID) Mutate() OwnerID {
	oid.Stamp = append(oid.Stamp, brand.Next())
	return oid
}

// MarshalJSON build json representation of OwnerID
func (oid OwnerID) MarshalJSON() ([]byte, error) {
	content := oid.String()
	return json.Marshal(content)
}

// UnmarshalJSON builds OwnerID from json
func (oid *OwnerID) UnmarshalJSON(b []byte) error {
	var oidstr string
	err := json.Unmarshal(b, &oidstr)
	if err != nil {
		return err
	}
	*oid, err = DecodeObjectIDString(oidstr)
	return err
}

// Equal is true if all fields are equal
func (oid OwnerID) Equal(oth OwnerID) bool {
	if oid.Type != oth.Type {
		return false
	}
	if oid.UserDefined != oth.UserDefined {
		return false
	}
	if len(oid.Stamp) != len(oth.Stamp) {
		return false
	}
	for i, s := range oid.Stamp {
		if oth.Stamp[i] != s {
			return false
		}
	}
	return true
}

// Owner is the interface for types that can own things
type Owner interface {
	GetID() OwnerID
	GetName() string
	Match(o Owner) bool
	Equal(o Owner) bool
	Copy() Owner
}

type publicowner struct {
}

// Public is the owner that matches all other owners
var Public publicowner

//GetID implements Owner
func (p publicowner) GetID() OwnerID {
	return OwnerID{
		Type:        'p',
		UserDefined: [3]byte{'a', 'l', 'l'},
	}
}

// GetName implements Owner
func (p publicowner) GetName() string {
	return "Public"
}

// Match implements Owner, always true
func (p publicowner) Match(_ Owner) bool {
	return true
}

// Equal implements Owner
func (p publicowner) Equal(o Owner) bool {
	switch o.(type) {
	case publicowner:
		return true
	default:
		return false
	}
}

// Copy implements Owner, Public is an empty value so doesn't actually allocate a new object
func (p publicowner) Copy() Owner {
	return Public
}

func publicfromjson(_ OwnerID, _ map[string]interface{}) (Owner, error) {
	return Public, nil
}

func publicfrombson(_ OwnerID, _ bson.M) (Owner, error) {
	return Public, nil
}
