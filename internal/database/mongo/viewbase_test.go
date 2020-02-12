package mongo

import (
	"bytes"
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestViewbase(t *testing.T) {
	t.Parallel()
	var vb *Viewbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestView"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatalf("Unable to init database: %s", err)
		}
		vb = db.View(context.Background()).(*Viewbase)
	}
	{
		contentString := "View version of a file FFFFF*****"
		contentSize := int64(len(contentString) + 1)
		inputVS := &database.ViewStore{
			ID: filehash.StoreID{
				Hash:  12345,
				Stamp: 678,
			},
			Content:  []byte(contentString),
			FileSize: contentSize,
		}
		t.Run("Insert", func(t *testing.T) {
			err := vb.Insert(inputVS)
			if err != nil {
				t.Fatalf("error inserting: %s", err)
			}
		})
		t.Run("Get", func(t *testing.T) {
			result, err := vb.Get(inputVS.ID)
			if err != nil {
				t.Fatalf("error getting: %s", err)
			}
			if !inputVS.ID.Equal(result.ID) ||
				!bytes.Equal(inputVS.Content, result.Content) ||
				inputVS.FileSize != result.FileSize {
				t.Errorf("Did not get correct file store:\ngot: %+#v\nexpected: %+#v\n", result, inputVS)
			}
		})
	}
}
