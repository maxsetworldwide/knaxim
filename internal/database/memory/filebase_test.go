package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

var fid = filehash.FileID{
	StoreID: sid,
	Stamp:   []byte{'a', 'b'},
}

func TestFiles(t *testing.T) {
	t.Parallel()
	defer testingComplete.Done()
	fb := DB.File(nil)
	defer fb.Close(nil)

	storedid, err := fb.Reserve(fid)
	if err != nil {
		t.Fatalf("unable to reserve id: %s", err)
	}
	if !storedid.Equal(fid) {
		t.Fatalf("incorrect storeid reserved: %v", storedid)
	}

	file := &database.File{
		Permission: database.Permission{
			Own: test1,
		},
		ID:   fid,
		Name: "TestFile",
	}

	err = fb.Insert(file)
	if err != nil {
		t.Fatalf("failed to insert file: %s", err)
	}

	gotten, err := fb.Get(fid)
	if err != nil {
		t.Fatalf("failed to find file: %s", err)
	}
	if !gotten.GetID().Equal(fid) {
		t.Fatalf("incorrect gotten file: %s", err)
	}

	allgot, err := fb.GetAll(fid)
	if err != nil {
		t.Fatalf("failed to find files: %s", err)
	}
	if len(allgot) != 1 || !allgot[0].GetID().Equal(fid) {
		t.Fatalf("incorrect GetAll: %v", allgot)
	}

	file.SetPerm(test2, "view", true)
	err = fb.Update(file)
	if err != nil {
		t.Fatalf("failed to update file: %s", err)
	}

	owned, err := fb.GetOwned(test1.GetID())
	if err != nil {
		t.Fatalf("failed to GetOwned: %s", err)
	}
	if len(owned) != 1 || !owned[0].GetID().Equal(fid) {
		t.Fatalf("incorrect return from owned: %v", owned)
	}

	shared, err := fb.GetPermKey(test2.GetID(), "view")
	if err != nil {
		t.Fatalf("failed to GetPermKey: %s", err)
	}
	if len(shared) != 1 || !shared[0].GetID().Equal(fid) {
		t.Fatalf("incorrect return from shared: %v", shared)
	}

	matched, err := fb.MatchStore(test1.GetID(), []filehash.StoreID{sid})
	if err != nil {
		t.Fatalf("unable to match sid: %s", err)
	}
	if len(matched) != 1 || !matched[0].GetID().Equal(fid) {
		t.Fatalf("incorrect returned matched: %v", matched)
	}

	err = fb.Remove(fid)
	if err != nil {
		t.Fatalf("failed to remove: %s", err)
	}

	expectnil, err := fb.Get(fid)
	if err == nil || expectnil != nil {
		t.Fatalf("found removed file: %v, %s", expectnil, err)
	}
}
