package mongo

import (
	"context"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
)

func putStorePlaceholder(db *Database, sid filehash.StoreID) error {
	placeholder, err := database.NewFileStore(strings.NewReader(sid.String()))
	if err != nil {
		return err
	}
	placeholder.ID = sid
	placeholder.ContentType = "test"
	placeholder.FileSize = 42
	placeholder.Perr = nil
	sb := db.Store(nil)
	_, err = sb.Reserve(sid)
	if err != nil {
		return err
	}
	return sb.Insert(placeholder)
}

func TestTagbase(t *testing.T) {
	t.Parallel()
	var tb *Tagbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestTag"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to init database", err)
		}
		tb = db.Tag(context.Background()).(*Tagbase)
	}
	fileids := []filehash.FileID{
		filehash.FileID{
			StoreID: filehash.StoreID{
				Hash:  9843,
				Stamp: 354,
			},
			Stamp: []byte("aa"),
		},
		filehash.FileID{
			StoreID: filehash.StoreID{
				Hash:  15468,
				Stamp: 98,
			},
			Stamp: []byte("bb"),
		},
	}
	for _, fid := range fileids {
		putStorePlaceholder(&tb.Database, fid.StoreID)
	}
	stags := []tag.Tag{
		tag.Tag{
			Word: "basic",
			Type: tag.CONTENT,
		},
		tag.Tag{
			Word: "multiple",
			Type: tag.CONTENT,
		},
	}
	ftags := []tag.Tag{
		tag.Tag{
			Word: "folder",
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					"userid": "test",
				},
			},
		},
		tag.Tag{
			Word: "multiple",
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					"userid": "test",
				},
			},
		},
	}
	t.Run("UpsertFile", func(t *testing.T) {
		err := tb.UpsertFile(fileids[0], ftags...)
		if err != nil {
			t.Fatal("failed to UpsertFile 1", err)
		}
		err = tb.UpsertFile(fileids[1], ftags...)
		if err != nil {
			t.Fatal("failed to UpsertFile 2", err)
		}
	})
	t.Run("UpsertStore", func(t *testing.T) {
		err := tb.UpsertStore(fileids[0].StoreID, stags...)
		if err != nil {
			t.Fatal("failed to UpsertStore 1", err)
		}
		err = tb.UpsertStore(fileids[1].StoreID, stags...)
		if err != nil {
			t.Fatal("failed to UpsertStore 2", err)
		}
	})
	t.Run("FileTags", func(t *testing.T) {
		result, err := tb.FileTags(fileids...)
		if err != nil {
			t.Fatal("failed to find tags", err)
		}
		for _, matches := range result {
			if len(matches) != len(ftags)+len(stags) {
				t.Fatal("incorrect returned tags", matches)
			}
		}
	})
	// t.Run("GetFiles", func(t *testing.T) {
	// 	files, stores, err := tb.GetFiles([]tag.Tag{tag.Tag{
	// 		Word: "multiple",
	// 		Type: tag.ALLTYPES,
	// 	}})
	// 	if err != nil {
	// 		t.Fatal("unable to get file references", err)
	// 	}
	// 	if len(files) != len(fileids) {
	// 		t.Fatalf("file tag not found:\n%+#v\n%+#v", files, stores)
	// 	}
	// 	if len(stores) != len(fileids) {
	// 		t.Fatalf("store tag not found:\n%+#v\n%+#v", files, stores)
	// 	}
	// })
	t.Run("GetFiles=users", func(t *testing.T) {
		files, stores, err := tb.GetFiles([]tag.Tag{tag.Tag{
			Word: "multiple",
			Type: tag.ALLTYPES,
			Data: tag.Data{
				tag.USER: map[string]string{
					"userid": "test",
				},
			},
		}})
		if err != nil {
			t.Fatal("unable to get file references", err)
		}
		if len(files) != len(fileids) {
			t.Fatalf("file tags not found:\n%+#v\n%+#v", files, stores)
		}
		if len(stores) != len(fileids) {
			t.Fatalf("store tags not found:\n%+#v\n%+#v", files, stores)
		}
	})
	t.Run("SearchData", func(t *testing.T) {
		tags, err := tb.SearchData(tag.USER, tag.Data{
			tag.USER: map[string]string{
				"userid": "test",
			},
		})
		if err != nil {
			t.Fatal("unable to Search Data", err)
		}
		if len(tags) != 4 {
			t.Fatal("incorrect return: ", tags)
		}
	})
}
