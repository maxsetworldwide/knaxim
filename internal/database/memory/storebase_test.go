package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

var sid = filehash.StoreID{
	Hash:  10,
	Stamp: 10,
}

func TestStore(t *testing.T) {
	t.Parallel()
	defer testingComplete.Done()
	sb := DB.Store(nil)
	defer sb.Close(nil)

	fs := &database.FileStore{
		ID:          sid,
		Content:     []byte("this is the content"),
		ContentType: "testfile",
		FileSize:    1234,
	}

	savedid, err := sb.Reserve(sid)
	if err != nil {
		t.Fatalf("Unable to reserve id: %s", err)
	}
	if !savedid.Equal(sid) {
		t.Fatalf("wrong id saved: %v", savedid)
	}
	err = sb.Insert(fs)
	if err != nil {
		t.Fatalf("unable to insert: %s", err)
	}

	gotten, err := sb.Get(sid)
	if err != nil {
		t.Fatalf("unable to get filestore: %s", err)
	}
	if !gotten.ID.Equal(sid) {
		t.Fatalf("incorrect gotten filestore: %v", gotten)
	}

	matched, err := sb.MatchHash(10)
	if err != nil {
		t.Fatalf("unable to match hash: %s", err)
	}
	if len(matched) != 1 || !matched[0].ID.Equal(sid) {
		t.Fatalf("incorrect matches: %v", matched)
	}
}