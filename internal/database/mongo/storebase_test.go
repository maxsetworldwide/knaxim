package mongo

import (
	"bytes"
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestStorebase(t *testing.T) {
	t.Parallel()
	var sb *Storebase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestStore"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to init database", err)
		}
		sb = db.Store(context.Background()).(*Storebase)
	}
	{
		input := filehash.StoreID{
			Hash:  24098,
			Stamp: 123,
		}
		t.Run("Reserve=basic", func(t *testing.T) {
			out, err := sb.Reserve(input)
			if err != nil {
				t.Error("error reserving", err)
			}
			if out.Hash != input.Hash || out.Stamp != input.Stamp {
				t.Error("return mismatch", out)
			}
		})
		t.Run("Reserve=mutate", func(t *testing.T) {
			out, err := sb.Reserve(input)
			if err != nil {
				t.Error("error reserving", err)
			}
			if out.Hash != input.Hash || out.Stamp != input.Mutate().Stamp {
				t.Error("return mismatch", out)
			}
		})
	}
	{
		input := &database.FileStore{
			ID: filehash.StoreID{
				Hash:  24098,
				Stamp: 123,
			},
			Content:     []byte("Here is a test file.#$%!1234t##g1"),
			ContentType: "text",
			FileSize:    33,
		}
		t.Run("Insert", func(t *testing.T) {
			err := sb.Insert(input)
			if err != nil {
				t.Fatal("error inserting", err)
			}
		})
		t.Run("Get", func(t *testing.T) {
			out, err := sb.Get(input.ID)
			if err != nil {
				t.Fatal("error getting", err)
			}
			if !input.ID.Equal(out.ID) ||
				!bytes.Equal(input.Content, out.Content) ||
				input.ContentType != out.ContentType ||
				input.FileSize != out.FileSize {
				t.Error("did not get correct file store", out)
			}
		})
		t.Run("MatchHash", func(t *testing.T) {
			out, err := sb.MatchHash(input.ID.Hash)
			if err != nil {
				t.Fatal("error match hash", err)
			}
			if !input.ID.Equal(out[0].ID) ||
				!bytes.Equal(input.Content, out[0].Content) ||
				input.ContentType != out[0].ContentType ||
				input.FileSize != out[0].FileSize {
				t.Error("did not get correct file store", out[0])
			}
		})
		t.Run("Update", func(t *testing.T) {
			input.Perr = &database.ProcessingError{
				Status:  420,
				Message: "Hey, You see this",
			}
			err := sb.UpdateMeta(input)
			if err != nil {
				t.Fatalf("unable to UpdateMeta: %s", err)
			}
		})
	}
}
