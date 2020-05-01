package memory

import (
	"bytes"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

var contentString = "View version of a file for the memory database"

var inputVS = &types.ViewStore{
	ID: types.StoreID{
		Hash:  98765,
		Stamp: 4321,
	},
	Content: []byte(contentString),
}

func TestViewbase(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	vb := DB.View()
	defer DB.Close(nil)
	t.Parallel()

	t.Log("View Insert")
	err := vb.Insert(inputVS)
	if err != nil {
		t.Fatalf("error inserting view: %s\n", err)
	}
	t.Log("View Get")
	result, err := vb.Get(inputVS.ID)
	if err != nil {
		t.Fatalf("error getting view: %s\n", err)
	}
	if !inputVS.ID.Equal(result.ID) || !bytes.Equal(inputVS.Content, result.Content) {
		t.Fatalf("Did not get correct view store:\ngot: %+#v\nexpected: %+#v\n", result, inputVS)
	}
}
