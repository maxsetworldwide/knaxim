package mongo

import (
	"context"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

func putStorePlaceholder(db *Database, sid types.StoreID) error {
	placeholder, err := types.NewFileStore(strings.NewReader(sid.String()))
	if err != nil {
		return err
	}
	placeholder.ID = sid
	placeholder.ContentType = "test"
	placeholder.FileSize = 42
	placeholder.Perr = nil
	sb := db.Store()
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
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		mdb, err := db.Connect(methodtesting)
		if err != nil {
			t.Fatalf("Unable to connect to database: %s", err.Error())
		}
		defer mdb.Close(methodtesting)
		tb = mdb.Tag().(*Tagbase)
	}
	fileids := []types.FileID{
		types.FileID{
			StoreID: types.StoreID{
				Hash:  9843,
				Stamp: 354,
			},
			Stamp: []byte("aa"),
		},
		types.FileID{
			StoreID: types.StoreID{
				Hash:  15468,
				Stamp: 98,
			},
			Stamp: []byte("bb"),
		},
	}
	ownerids := []types.OwnerID{
		types.OwnerID{
			Type:        't',
			UserDefined: [3]byte{'e', 's', 't'},
			Stamp:       []byte{'1'},
		},
		types.OwnerID{
			Type:        't',
			UserDefined: [3]byte{'e', 's', 't'},
			Stamp:       []byte{'2'},
		},
	}
	t.Run("Upsert", func(t *testing.T) {
		err := tb.Upsert([]tag.FileTag{
			tag.FileTag{
				File:  fileids[0],
				Owner: ownerids[0],
				Tag: tag.Tag{
					Word: "first",
					Type: tag.CONTENT | tag.USER,
					Data: tag.Data{
						tag.CONTENT: map[string]interface{}{"count": 1},
						tag.USER:    map[string]interface{}{"protect": true},
					},
				},
			},
			tag.FileTag{
				File:  fileids[0],
				Owner: ownerids[1],
				Tag: tag.Tag{
					Word: "second",
					Type: tag.TOPIC | tag.RESOURCE,
				},
			},
			tag.FileTag{
				File:  fileids[1],
				Owner: ownerids[0],
				Tag: tag.Tag{
					Word: "third",
					Type: tag.ACTION | tag.PROCESS,
				},
			},
			tag.FileTag{
				File:  fileids[1],
				Owner: ownerids[1],
				Tag: tag.Tag{
					Word: "fourth",
					Type: tag.ALLSTORE | tag.ALLFILE,
				},
			},
		}...)
		if err != nil {
			t.Errorf("unable to upsert tags: %s", err.Error())
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get", func(t *testing.T) {
		tags, err := tb.Get(fileids[0], ownerids[0])
		if err != nil {
			t.Fatalf("unable to get file tags: %s", err.Error())
		}
		if len(tags) != 2 {
			t.Fatalf("incorrect returned tags: %v", tags)
		}
	})
	t.Run("GetType", func(t *testing.T) {
		tags, err := tb.GetType(fileids[1], ownerids[1], tag.ACTION)
		if err != nil {
			t.Fatalf("unable to get file tags: %s", err.Error())
		}
		if len(tags) != 2 {
			t.Fatalf("incorrect returned tags: %v", tags)
		}
	})
	t.Run("GetAll", func(t *testing.T) {
		tags, err := tb.GetAll(tag.USER, ownerids[0])
		if err != nil {
			t.Fatalf("unable to get file tags: %s", err.Error())
		}
		if len(tags) != 1 {
			t.Fatalf("incorrect returned tags: %v", tags)
		}
	})
	t.Run("SearchFiles", func(t *testing.T) {
		ids, err := tb.SearchFiles(fileids, tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "first",
				Type: tag.ALLTYPES,
			},
		})
		if err != nil {
			t.Fatalf("unable to search file tags: %s", err.Error())
		}
		if len(ids) != 1 {
			t.Fatalf("incorrect returned ids: %v", ids)
		}
	})
	t.Run("SearchFiles Regex", func(t *testing.T) {
		ids, err := tb.SearchFiles(fileids, tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "fir",
				Type: tag.SEARCH | tag.CONTENT,
				Data: tag.Data{
					tag.SEARCH: map[string]interface{}{
						"regex": true,
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("unable to search file tags: %s", err.Error())
		}
		if len(ids) != 1 {
			t.Fatalf("incorrect returned ids: %v", ids)
		}
	})
}
