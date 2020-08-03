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

package mongo

import (
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func TestOwnerbase(t *testing.T) {
	t.Parallel()
	var ob *Ownerbase
	db := new(Database)
	*db = *configuration.DB
	db.DBName = "TestOwners"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := db.Init(ctx, true); err != nil {
		t.Fatal("Inable to init database", err)
	}
	methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	mdb, err := db.Connect(methodtesting)
	if err != nil {
		t.Fatalf("Unable to connect to database: %s", err.Error())
	}
	defer mdb.Close(methodtesting)
	ob = mdb.Owner().(*Ownerbase)
	var ownerids = []types.OwnerID{
		types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try"),
		},
		types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try"),
		},
		types.OwnerID{
			Type:        'g',
			UserDefined: [3]byte{'t', 'g', 'r'},
			Stamp:       []byte("test"),
		},
	}
	var data = []types.Owner{
		&types.User{
			ID:   ownerids[0],
			Name: "devon",
			Pass: types.UserCredential{
				Salt: []byte("thisisthesalt"),
				Hash: []byte("thisisthehash"),
			},
			Email:            "devon@test.test",
			CookieSig:        []byte("thisisthecookiesig"),
			CookieInactivity: time.Now().Add(time.Hour * 2),
			CookieTimeout:    time.Now().Add(time.Hour * 24),
		},
		&types.User{
			ID:   ownerids[1],
			Name: "developer",
			Pass: types.UserCredential{
				Salt: []byte("thisisthesalt2"),
				Hash: []byte("thisisthehash2"),
			},
			Email:            "develop@test.test",
			CookieSig:        []byte("thisisthecookiesig2"),
			CookieInactivity: time.Now().Add(time.Hour * 2),
			CookieTimeout:    time.Now().Add(time.Hour * 24),
		},
		&types.Group{
			ID:   ownerids[2],
			Name: "testGroup",
		},
	}
	data[2].(*types.Group).Own = data[0]
	data[2].(*types.Group).AddMember(data[1])
	t.Run("Reserve", func(t *testing.T) {
		for i, ele := range data {
			var err error
			switch v := ele.(type) {
			case *types.User:
				v.ID, err = ob.Reserve(v.ID, v.Name)
			case *types.Group:
				v.ID, err = ob.Reserve(v.ID, v.Name)
			}
			if err != nil {
				t.Fatal("unable to reserve id", err)
			}
			t.Logf("Index: %d;ID: %v", i, ele.GetID())
			data[i] = ele
			ownerids[i] = ele.GetID()
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Insert", func(t *testing.T) {
		for _, owner := range data {
			err := ob.Insert(owner)
			if err != nil {
				t.Fatal("unable to insert", owner, err)
			}
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get=User", func(t *testing.T) {
		result, err := ob.Get(ownerids[0])
		if err != nil {
			t.Fatal(err)
		}
		if !result.Equal(data[0]) {
			t.Error("incorrect returned user", result)
		}
	})
	t.Run("Get=Group", func(t *testing.T) {
		result, err := ob.Get(ownerids[2])
		if err != nil {
			t.Fatal(err)
		}
		if !result.Equal(data[2]) {
			t.Error("incorrect returned group", result)
		}
	})
	ob.trackOwners = newTrackOwners()
	t.Run("FindUserName", func(t *testing.T) {
		result, err := ob.FindUserName("devon")
		if err != nil {
			t.Fatal(err)
		}
		if !data[0].Equal(result) {
			t.Error("incorrect returned document")
		}
	})
	t.Run("FindGroupName", func(t *testing.T) {
		result, err := ob.FindGroupName("testGroup")
		if err != nil {
			t.Fatal(err)
		}
		if !data[2].Equal(result) {
			t.Error("incorrect returned document", result)
		}
	})
	ob.trackOwners = newTrackOwners()
	t.Run("GetGroups=Owner", func(t *testing.T) {
		owned, member, err := ob.GetGroups(ownerids[0])
		if err != nil {
			t.Fatal("failed to get groups", err)
		}
		if len(owned) != 1 && len(member) != 0 {
			t.Fatal("incorrect returns", owned, member)
		}
		if !data[2].Equal(owned[0]) {
			t.Fatal("group not correct", owned[0])
		}
	})
	t.Run("GetGroups=Member", func(t *testing.T) {
		owned, member, err := ob.GetGroups(ownerids[1])
		if err != nil {
			t.Fatal("failed to get groups", err)
		}
		if len(owned) != 0 && len(member) != 1 {
			t.Fatal("incorrect returns", owned, member)
		}
		if !data[2].Equal(member[0]) {
			t.Fatal("group not correct", member[0])
		}
	})

	t.Run("Update", func(t *testing.T) {
		newowner := data[2].(*types.Group)
		newowner.SetName("updatedgroup")
		err := ob.Update(newowner)
		if err != nil {
			t.Fatal("failed to update", err)
		}
		current, err := ob.Get(newowner.GetID())
		if err != nil {
			t.Fatal("failed to get changed", err)
		}
		if current.(types.GroupI).GetName() != newowner.GetName() {
			t.Fatalf("update had no effect: %+#v", current)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		key, err := ob.GetResetKey(data[0].GetID())
		if err != nil {
			t.Fatal("unable to get reset key: ", err)
		}
		t.Log("key: ", key)
		oid, err := ob.CheckResetKey(key)
		if err != nil {
			t.Fatal("unable to check reset key")
		}
		if !oid.Equal(data[0].GetID()) {
			t.Fatal("incorrect returned Owner ID: ", oid)
		}
	})
}
