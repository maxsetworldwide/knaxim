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

package types

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

// GroupI is the group interface
type GroupI interface {
	Owner
	PermissionI
	SetName(string)
	GetMembers() []Owner
	AddMember(o Owner)
	RemoveMember(o Owner)
}

// Group is the basic Group implementation
type Group struct {
	Permission
	ID   OwnerID `json:"id" bson:"id"`
	Name string  `json:"name" bson:"name"`
	Max  int64   `json:"max,omitempty" bson:"max,omitempty"`
}

// NewGroup build Group with a particular name and owner
func NewGroup(name string, o Owner) *Group {
	ng := new(Group)
	ng.Own = o
	ng.Name = name

	var nid OwnerID
	nid.Type = 'g'
	nid.UserDefined = name2userdefined(name)
	// nid.Stamp = newstamp()

	ng.ID = nid

	return ng
}

// MaxFiles returns the maximum number of files that an owner can own
func (g *Group) MaxFiles() int64 {
	return g.Max
}

// GetID implements GroupI
func (g *Group) GetID() OwnerID {
	return g.ID
}

// GetName implements GroupI
func (g *Group) GetName() string {
	return g.Name
}

// SetName implements GroupI
func (g *Group) SetName(s string) {
	g.Name = s
}

// Match returns true if owner is equal to group, matches
// group's owner, or matches any of the group's members
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

// Equal returns true if the Owner has Equal ID values
func (g *Group) Equal(o Owner) bool {
	if o == nil {
		return false
	}
	return g.ID.Equal(o.GetID())
}

// Copy Group
func (g *Group) Copy() Owner {
	if g == nil {
		return nil
	}
	ng := new(Group)
	*ng = *g
	ng.Permission = *(g.CopyPerm(nil).(*Permission))
	return ng
}

// GetMembers implements GroupI
func (g *Group) GetMembers() []Owner {
	result := make([]Owner, 0, len(g.GetPerm("%member%")))
	return append(result, g.GetPerm("%member%")...)
}

// AddMember implements GroupI
func (g *Group) AddMember(o Owner) {
	if o == nil {
		return
	}
	if g.Equal(o) {
		return
	}
	g.SetPerm(o, "%member%", true)
}

// RemoveMember implements GroupI
func (g *Group) RemoveMember(o Owner) {
	if o == nil {
		return
	}
	g.SetPerm(o, "%member%", false)
}

// MarshalJSON builds json representation
func (g *Group) MarshalJSON() ([]byte, error) {
	vals := g.toMap()
	vals["id"] = g.ID
	vals["name"] = g.Name
	return json.Marshal(vals)
}

// MarshalBSON builds bson representation
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

// UnmarshalJSON decodes to Group
// group still needs Populate called inorder to load owner and members
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

// UnmarshalBSON decodes to Group
// group still needs Populate to be called inorder to load owner and members
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
