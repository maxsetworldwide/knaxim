package mongo

import (
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestContenbase(t *testing.T) {
	t.Parallel()
	var cb *Contentbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestContent"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to Init database", err)
		}
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		cb = db.Content(methodtesting).(*Contentbase)
	}
	var fileids = []filehash.StoreID{
		filehash.StoreID{
			Hash:  7777,
			Stamp: 32621,
		},
		filehash.StoreID{
			Hash:  841602,
			Stamp: 28720,
		},
	}
	var fileStores = []*database.FileStore{
		&database.FileStore{
			ID:          fileids[0],
			Content:     []byte("asdfasdf"),
			ContentType: "test",
			FileSize:    420,
		},
		&database.FileStore{
			ID:          fileids[1],
			Content:     []byte("fdafdsa"),
			ContentType: "test",
			FileSize:    240,
		},
	}
	{
		sb := cb.Store(nil)
		for _, fs := range fileStores {
			_, err := sb.Reserve(fs.ID)
			if err != nil {
				t.Fatalf("unable to Reserve file store id: %s", err)
			}
			err = sb.Insert(fs)
			if err != nil {
				t.Fatalf("unable to Insert file store: %s", err)
			}
		}
	}
	var data = []database.ContentLine{
		database.ContentLine{
			ID:       fileids[0],
			Position: 0,
			Content:  []string{"This is the first sentence."},
		},
		database.ContentLine{
			ID:       fileids[0],
			Position: 1,
			Content:  []string{"Another Sentence right here."},
		},
		database.ContentLine{
			ID:       fileids[0],
			Position: 2,
			Content:  []string{"More Sentences."},
		},
		database.ContentLine{
			ID:       fileids[1],
			Position: 0,
			Content:  []string{"This is another document."},
		},
		database.ContentLine{
			ID:       fileids[1],
			Position: 1,
			Content:  []string{"It only has 2 sentences"},
		},
	}
	t.Run("Insert", func(t *testing.T) {
		err := cb.Insert(data...)
		if err != nil {
			t.Error("Unable to Insert", err)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Len", func(t *testing.T) {
		l, err := cb.Len(fileids[0])
		if err != nil {
			t.Fatal("Err getting Length", err)
		}
		if l != 3 {
			t.Fatal("Err incorrect length", l)
		}
	})
	t.Run("Slice", func(t *testing.T) {
		result, err := cb.Slice(fileids[0], 1, 3)
		if err != nil {
			t.Fatal("Err getting slice", err)
		}
		if len(result) != 2 {
			t.Fatal("Slice returned incorrect amount", result)
		}
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("Equality check paniced", r)
			}
		}()
		for _, r := range result {
			original := data[r.Position]
			if !original.ID.Equal(r.ID) || original.Content[0] != r.Content[0] {
				t.Error("result mismatched original", r, original)
			}
		}
	})
	t.Run("Regex", func(t *testing.T) {
		result, err := cb.RegexSearchFile("only", fileids[1], 0, 2)
		if err != nil {
			t.Fatal("Err doing search", err)
		}
		if len(result) != 1 {
			t.Fatal("incorrect number of matches", result)
		}
		if !result[0].ID.Equal(data[4].ID) || result[0].Position != data[4].Position || result[0].Content[0] != data[4].Content[0] {
			t.Error("mismatched result", result[0])
		}
	})
}
