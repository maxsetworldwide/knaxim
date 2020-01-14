package database

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"git.maxset.io/server/knaxim/database/brand"

	"go.mongodb.org/mongo-driver/bson"
)

type OwnerID struct {
	Type        byte    `bson:"type"`
	UserDefined [3]byte `bson:"ud"`
	Stamp       []byte  `bson:"stamp"`
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

func newstamp() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(time.Now().UnixNano()))
	return buf.Bytes()
}

func (oid OwnerID) bytes() []byte {
	b := make([]byte, 0, 4+len(oid.Stamp))
	b = append(b, oid.Type)
	b = append(b, oid.UserDefined[:]...)
	b = append(b, oid.Stamp...)
	return b
}

func (oid OwnerID) String() string {
	out := new(strings.Builder)
	enc := base64.NewEncoder(base64.RawURLEncoding, out)
	enc.Write(oid.bytes())
	enc.Close()
	return out.String()
}

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
	var out OwnerID
	out.Type = b[0]
	for i, b := range b[1:4] {
		out.UserDefined[i] = b
	}
	out.Stamp = b[4:]
	return out, nil
}

func (oid OwnerID) Mutate() OwnerID {
	oid.Stamp = append(oid.Stamp, brand.Next())
	return oid
}

func (oid OwnerID) MarshalJSON() ([]byte, error) {
	content := oid.String()
	return json.Marshal(content)
}

func (oid *OwnerID) UnmarshalJSON(b []byte) error {
	oid = new(OwnerID)
	var oidstr *string
	err := json.Unmarshal(b, oidstr)
	if err != nil {
		return err
	}
	*oid, err = DecodeObjectIDString(*oidstr)
	return err
}

// func (oid OwnerID) MarshalBSON() ([]byte, error) {
// 	content := oid.String()
// 	return bson.Marshal(content)
// }

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

// TODO: in mongo add private type to handle loading owners from database
// construction type needed that mirrors Permission Object but with OwnerIDs
// instead of full references to Owners in order to handle circular references
// due to Groups

type Owner interface {
	GetID() OwnerID
	Match(o Owner) bool
	Equal(o Owner) bool
}

type publicowner struct {
	ID OwnerID `bson:"id"`
}

var Public = publicowner{
	ID: OwnerID{
		Type:        'p',
		UserDefined: [3]byte{'p', 'u', 'b'},
	},
}

func (p publicowner) GetID() OwnerID {
	return p.ID
}

func (p publicowner) Match(_ Owner) bool {
	return true
}

func (p publicowner) Equal(o Owner) bool {
	switch o.(type) {
	case publicowner:
		return true
	default:
		return false
	}
}

func publicfromjson(_ OwnerID, _ map[string]interface{}) (Owner, error) {
	return Public, nil
}

func publicfrombson(_ OwnerID, _ bson.M) (Owner, error) {
	return Public, nil
}

type PermissionI interface {
	json.Marshaler
	json.Unmarshaler
	bson.Marshaler
	bson.Unmarshaler
	GetOwner() Owner
	CheckPerm(Owner, string) bool
	GetPerm(string) []Owner
	PermTypes() []string
	SetPerm(Owner, string, bool)
	CopyPerm(Owner) PermissionI
	Populate(Ownerbase) error
}

type Permission struct {
	Own    Owner              `json:"own" bson:"own"`
	Perm   map[string][]Owner `json:"perm,omitempty" bson:"perm,omitempty"`
	ownid  OwnerID
	permid map[string][]OwnerID
}

func (p *Permission) GetOwner() Owner {
	return p.Own
}

func (p *Permission) CheckPerm(o Owner, s string) bool {
	for _, v := range p.Perm[s] {
		if v.Match(o) {
			return true
		}
	}
	return false
}

func (p *Permission) GetPerm(s string) []Owner {
	def := make([]Owner, len(p.Perm[s]))
	copy(def, p.Perm[s])
	return def
}

func (p *Permission) PermTypes() []string {
	out := make([]string, 0, len(p.Perm))
	for k, _ := range p.Perm {
		out = append(out, k)
	}
	return out
}

func (p *Permission) SetPerm(o Owner, s string, b bool) {
	if o != nil {
		if b {
			if p.Perm == nil {
				p.Perm = make(map[string][]Owner)
			}
			for _, v := range p.Perm[s] {
				if v.Equal(o) {
					return
				}
			}
			p.Perm[s] = append(p.Perm[s], o)
		} else {
			mi := -1
			for i, v := range p.Perm[s] {
				if v.Equal(o) {
					mi = i
					break
				}
			}
			if mi > -1 {
				p.Perm[s] = append(p.Perm[s][:mi], p.Perm[s][mi+1:]...)
			}
		}
	}
}

func (p *Permission) CopyPerm(newowner Owner) PermissionI {
	n := new(Permission)
	if newowner == nil {
		n.Own = p.Own
	} else {
		n.Own = newowner
	}
	if p.Perm != nil {
		n.Perm = make(map[string][]Owner)
		for k, v := range p.Perm {
			n.Perm[k] = make([]Owner, len(v))
			copy(n.Perm[k], v)
		}
	}
	return n
}

func (p *Permission) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["own"] = p.Own.GetID()
	if p.Perm != nil {
		perm := make(map[string][]OwnerID)
		for key, vals := range p.Perm {
			for _, o := range vals {
				perm[key] = append(perm[key], o.GetID())
			}
		}
		m["perm"] = perm
	}
	return m
}

func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.toMap())
}

func (p *Permission) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.toMap())
}

type pReference struct {
	Own  OwnerID              `bson:"own" json:"own"`
	Perm map[string][]OwnerID `bson:"perm,omitempty" json:"perm,omitempty"`
}

func (p *Permission) UnmarshalJSON(b []byte) error {
	form := new(pReference)
	err := json.Unmarshal(b, form)
	if err != nil {
		return err
	}
	p.ownid = form.Own
	p.permid = form.Perm
	return nil
}

func (p *Permission) UnmarshalBSON(b []byte) error {
	form := new(pReference)
	err := bson.Unmarshal(b, form)
	if err != nil {
		return err
	}
	p.ownid = form.Own
	p.permid = form.Perm
	return nil
}

func (p *Permission) Populate(ub Ownerbase) error {
	var err error
	p.Own, err = ub.Get(p.ownid)
	if err != nil {
		return err
	}
	for key, val := range p.permid {
		templist := make([]Owner, len(val))
		for i, id := range val {
			templist[i], err = ub.Get(id)
			if err != nil {
				return err
			}
		}
		if p.Perm == nil {
			p.Perm = make(map[string][]Owner)
		}
		p.Perm[key] = templist
	}
	return nil
}
