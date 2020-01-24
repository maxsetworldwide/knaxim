package database

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type GroupI interface {
	Owner
	PermissionI
	SetName(string)
	GetMembers() []Owner
	AddMember(o Owner)
	RemoveMember(o Owner)
}

type Group struct {
	Permission
	ID   OwnerID `json:"id" bson:"id"`
	Name string  `json:"name" bson:"name"`
}

func NewGroup(name string, o Owner) *Group {
	ng := new(Group)
	ng.Own = o
	ng.Name = name

	var nid OwnerID
	nid.Type = 'g'
	nid.UserDefined = name2userdefined(name)
	nid.Stamp = newstamp()

	ng.ID = nid

	return ng
}

func (g *Group) GetID() OwnerID {
	return g.ID
}

func (g *Group) GetName() string {
	return g.Name
}

func (g *Group) SetName(s string) {
	g.Name = s
}

func (g *Group) Match(o Owner) bool {
	if g.Equal(o) {
		return true
	}
	if g.GetOwner().Match(o) {
		return true
	}
	for _, m := range g.GetPerm("%member%") {
		if m.Match(o) {
			return true
		}
	}
	return false
}

func (g *Group) Equal(o Owner) bool {
	if o == nil {
		return false
	}
	return g.ID.Equal(o.GetID())
}

func (g *Group) Copy() Owner {
	if g == nil {
		return nil
	}
	ng := new(Group)
	*ng = *g
	ng.Permission = *(g.CopyPerm(nil).(*Permission))
	return ng
}

func (g *Group) GetMembers() []Owner {
	result := make([]Owner, 0, len(g.GetPerm("%member%")))
	return append(result, g.GetPerm("%member%")...)
}

func (g *Group) AddMember(o Owner) {
	if o == nil {
		return
	}
	g.SetPerm(o, "%member%", true)
}

func (g *Group) RemoveMember(o Owner) {
	if o == nil {
		return
	}
	g.SetPerm(o, "%member%", false)
}

func (g *Group) MarshalJSON() ([]byte, error) {
	vals := g.toMap()
	vals["id"] = g.ID
	vals["name"] = g.Name
	return json.Marshal(vals)
}

func (g *Group) MarshalBSON() ([]byte, error) {
	vals := g.toMap()
	vals["id"] = g.ID
	vals["name"] = g.Name
	return bson.Marshal(vals)
}

type gForm struct {
	ID   OwnerID `json:"id" bson:"id"`
	Name string  `json:"name" bson:"name"`
}

func (g *Group) UnmarshalJSON(b []byte) error {
	err := g.Permission.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	form := new(gForm)
	err = json.Unmarshal(b, form)
	if err != nil {
		return err
	}
	g.ID = form.ID
	g.Name = form.Name
	return nil
}

func (g *Group) UnmarshalBSON(b []byte) error {
	err := g.Permission.UnmarshalBSON(b)
	if err != nil {
		return err
	}
	form := new(gForm)
	err = bson.Unmarshal(b, form)
	if err != nil {
		return err
	}
	g.ID = form.ID
	g.Name = form.Name
	return nil
}
