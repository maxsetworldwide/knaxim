package database

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

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
