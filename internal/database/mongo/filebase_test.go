package mongo

import (
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
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
		fb = db.File(methodtesting).(*Filebase)
	}
	var fileids = []filehash.FileID{
		filehash.FileID{
			StoreID: filehash.StoreID{
				Hash:  901,
				Stamp: 613,
			},
			Stamp: []byte{34},
		},
		filehash.FileID{
			StoreID: filehash.StoreID{
				Hash:  604,
				Stamp: 1834,
			},
			Stamp: []byte{16},
		},
	}
	var ownerids = []database.OwnerID{
		database.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try"),
		},
		database.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'d', 'e', 'v'},
			Stamp:       []byte("try2"),
		},
		database.OwnerID{
			Type:        'g',
			UserDefined: [3]byte{'t', 'g', 'r'},
			Stamp:       []byte("test"),
		},
	}
	var owners = []database.Owner{
		&database.User{
			ID:   ownerids[0],
			Name: "devon",
			Pass: database.UserCredential{
				Salt: []byte("thisisthesalt"),
				Hash: []byte("thisisthehash"),
			},
			Email:            "devon@test.test",
			CookieSig:        []byte("thisisthecookiesig"),
			CookieInactivity: time.Now().Add(time.Hour * 2),
			CookieTimeout:    time.Now().Add(time.Hour * 24),
		},
		&database.User{
			ID:   ownerids[1],
			Name: "developer",
			Pass: database.UserCredential{
				Salt: []byte("thisisthesalt2"),
				Hash: []byte("thisisthehash2"),
			},
			Email:            "develop@test.test",
			CookieSig:        []byte("thisisthecookiesig2"),
			CookieInactivity: time.Now().Add(time.Hour * 2),
			CookieTimeout:    time.Now().Add(time.Hour * 24),
		},
		&database.Group{
			ID:   ownerids[2],
			Name: "testGroup",
		},
	}
	owners[2].(*database.Group).Own = owners[0]
	owners[2].(*database.Group).AddMember(owners[1])
	var files = []database.FileI{
		&database.File{
			Permission: database.Permission{
				Own: owners[0],
			},
			ID:   fileids[0],
			Name: "First.txt",
		},
		&database.File{
			Permission: database.Permission{
				Own: owners[1],
				Perm: map[string][]database.Owner{
					"view": []database.Owner{owners[2]},
				},
			},
			ID:   fileids[1],
			Name: "Second.txt",
		},
	}
	ob := fb.Owner(fb.GetContext())
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
		matched, err := fb.MatchStore(ownerids[0], []filehash.StoreID{fileids[0].StoreID}, "view")
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
		newfile.SetPerm(database.Public, "view", true)
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
	t.Run("Remove", func(t *testing.T) {
		err := fb.Remove(files[0].GetID())
		if err != nil {
			t.Fatal("failed to remove", err)
		}
		_, err = fb.Get(files[0].GetID())
		if err != database.ErrNotFound {
			t.Fatal("file not removed", err)
		}
	})
}
