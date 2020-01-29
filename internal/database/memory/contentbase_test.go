package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestContent(t *testing.T) {
	t.Parallel()
	cb := DB.Content(nil)
	defer cb.Close(nil)

	sid := filehash.StoreID{
		Hash:  10,
		Stamp: 10,
	}

	lines := []database.ContentLine{
		database.ContentLine{
			ID:       sid,
			Position: 0,
			Content:  []string{"this is the first line"},
		},
		database.ContentLine{
			ID:       sid,
			Position: 1,
			Content:  []string{"2nd line in content"},
		},
	}

	err := cb.Insert(lines...)
	if err != nil {
		t.Fatalf("Failed to insert lines: %s", err)
	}
	l, err := cb.Len(sid)
	if err != nil {
		t.Fatalf("Failed to get length: %s", err)
	}
	if l != 2 {
		t.Fatalf("Incorrect length: %d", l)
	}
	slice, err := cb.Slice(sid, 1, 2)
	if err != nil {
		t.Fatalf("Failed to get slice: %s", err)
	}
	if slice[0].Position != 1 {
		t.Fatalf("Incorrect Position: %d", slice[0].Position)
	}
	result, err := cb.RegexSearchFile("line", sid, 0, 2)
	if err != nil {
		t.Fatalf("Failed search: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("incorrect return: %v", result)
	}
}
