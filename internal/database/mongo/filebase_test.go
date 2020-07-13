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

package mongo

import (
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

func TestFilebase(t *testing.T) {
	t.Parallel()
	var fb *Filebase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestFiles"
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
		fb = mdb.File().(*Filebase)
	}
	var fileids = []types.FileID{
		types.FileID{
			StoreID: types.StoreID{
				Hash:  901,
				Stamp: 613,
			},
			Stamp: []byte{34},
		},
		types.FileID{
			StoreID: types.StoreID{
				Hash:  604,
				Stamp: 1834,
			},
			Stamp: []byte{16},
		},
	}
	var ownerids = []types.OwnerID{
		types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try"),
		},
		types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try2"),
		},
		types.OwnerID{
			Type:        'g',
			UserDefined: [3]byte{'t', 'g', 'r'},
			Stamp:       []byte("test"),
		},
	}
	var owners = []types.Owner{
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
	owners[2].(*types.Group).Own = owners[0]
	owners[2].(*types.Group).AddMember(owners[1])
	var files = []types.FileI{
		&types.File{
			Permission: types.Permission{
				Own: owners[0],
			},
			ID:   fileids[0],
			Name: "First.txt",
		},
		&types.File{
			Permission: types.Permission{
				Own: owners[1],
				Perm: map[string][]types.Owner{
					"view": []types.Owner{owners[2]},
				},
			},
			ID:   fileids[1],
			Name: "Second.txt",
		},
	}
	ob := fb.Owner()
	for i, oid := range ownerids {
		tempoid, err := ob.Reserve(oid, owners[i].(interface{ GetName() string }).GetName())
		if err != nil {
			t.Fatal("unable to reserve owner id", i, err)
		}
		if !oid.Equal(tempoid) {
			t.Fatalf("incorrect owner id inserted\nindex: %d\n%+#v", i, tempoid)
		}
	}
	for _, owner := range owners {
		err := ob.Insert(owner)
		if err != nil {
			t.Fatal("failed insert owners", err)
		}
	}
	t.Run("Reserve", func(t *testing.T) {
		for _, fid := range fileids {
			tempfid, err := fb.Reserve(fid)
			if err != nil {
				t.Fatal("unable to reserve file id", err)
			}
			if !fid.Equal(tempfid) {
				t.Fatal("incorrect id reserved", tempfid)
			}
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Insert", func(t *testing.T) {
		for _, file := range files {
			err := fb.Insert(file)
			if err != nil {
				t.Fatal("Unable to insert file", file, err)
			}
		}
	})
	t.Run("Get", func(t *testing.T) {
		for i, fid := range fileids {
			returned, err := fb.Get(fid)
			if err != nil {
				t.Fatal("unable to get fid", i, err)
			}
			if returned.GetName() != files[i].GetName() {
				t.Fatal("mismatched files", i, returned)
			}
			t.Logf("Found: %+#v", returned)
		}
	})
	t.Run("GetAll", func(t *testing.T) {
		returned, err := fb.GetAll(fileids...)
		if err != nil {
			t.Fatal("unable to get all", err)
		}
		t.Logf("Found: %+#v", returned)
		if len(returned) != len(fileids) {
			t.Fatal("did not find all expected")
		}
	})
	t.Run("GetOwned", func(t *testing.T) {
		matched, err := fb.GetOwned(ownerids[0])
		if err != nil {
			t.Fatal("unable to find file", err)
		}
		if len(matched) != 1 {
			t.Fatal("incorrect returned matched", matched)
		}
		if matched[0].GetName() != files[0].GetName() {
			t.Fatal("mismatched files", matched[0])
		}
	})
	t.Run("GetPermKey", func(t *testing.T) {
		matched, err := fb.GetPermKey(ownerids[2], "view")
		if err != nil {
			t.Fatal("unable to find file", err)
		}
		if len(matched) != 1 {
			t.Fatal("incorrect returned matched", matched)
		}
		if matched[0].GetName() != files[1].GetName() {
			t.Fatal("mismatched files", matched[0])
		}
	})
	t.Run("MatchStore", func(t *testing.T) {
		matched, err := fb.MatchStore(ownerids[0], []types.StoreID{fileids[0].StoreID}, "view")
		if err != nil {
			t.Fatal("failed to match store", err)
		}
		if len(matched) != 1 {
			t.Fatal("incorrect returned matched", matched)
		}
		if matched[0].GetName() != files[0].GetName() {
			t.Fatal("mismatched files", matched[0])
		}
	})
	t.Run("Update", func(t *testing.T) {
		newfile := files[0]
		newfile.SetName("update.txt")
		newfile.SetPerm(types.Public, "view", true)
		err := fb.Update(newfile)
		if err != nil {
			t.Fatal("failed to update", err)
		}
		state, err := fb.Get(newfile.GetID())
		if err != nil {
			t.Fatal("failed to get updated", err)
		}
		if state.GetName() != newfile.GetName() {
			t.Fatal("failed to update file", state)
		}
	})
	t.Run("Count", func(t *testing.T) {
		count, err := fb.Count(ownerids[0])
		if err != nil {
			t.Fatal("unable to count files: ", err.Error())
		}
		if count != 1 {
			t.Fatalf("incorrect count: %d", count)
		}
	})
	t.Run("Remove", func(t *testing.T) {
		err := fb.Remove(files[0].GetID())
		if err != nil {
			t.Fatal("failed to remove", err)
		}
		_, err = fb.Get(files[0].GetID())
		if se, ok := err.(srverror.Error); !ok || se.Status() != errors.ErrNotFound.Status() {
			t.Fatal("file not removed", err)
		}
	})
}
