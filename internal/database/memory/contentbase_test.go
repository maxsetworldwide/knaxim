package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func TestContent(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	cb := DB.Content()
	defer DB.Close(nil)
	t.Parallel()

	lines := []types.ContentLine{
		types.ContentLine{
			ID:       sid,
			Position: 0,
			Content:  []string{"this is the first line"},
		},
		types.ContentLine{
			ID:       sid,
			Position: 1,
			Content:  []string{"2nd line in content"},
		},
	}

	t.Log("Insert")
	err := cb.Insert(lines...)
	if err != nil {
		t.Fatalf("Failed to insert lines: %s", err)
	}

	t.Log("Len")
	l, err := cb.Len(sid)
	if err != nil {
		t.Fatalf("Failed to get length: %s", err)
	}
	if l != 2 {
		t.Fatalf("Incorrect length: %d", l)
	}

	t.Log("Slice")
	slice, err := cb.Slice(sid, 1, 2)
	if err != nil {
		t.Fatalf("Failed to get slice: %s", err.Error())
	}
	if slice[0].Position != 1 {
		t.Fatalf("Incorrect Position: %d", slice[0].Position)
	}

	t.Log("Regex")
	result, err := cb.RegexSearchFile("line", sid, 0, 2)
	if err != nil {
		t.Fatalf("Failed search: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("incorrect return: %v", result)
	}
}
