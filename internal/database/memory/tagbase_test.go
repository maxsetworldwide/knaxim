package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

func TestTag(t *testing.T) {
	defer testingComplete.Done()
	tb := DB.Tag(nil)
	defer tb.Close(nil)
	t.Parallel()
	fileids := []types.FileID{
		types.FileID{
			StoreID: types.StoreID{
				Hash:  123456,
				Stamp: 789,
			},
			Stamp: []byte("abc"),
		},
		types.FileID{
			StoreID: types.StoreID{
				Hash:  4568196,
				Stamp: 4897,
			},
			Stamp: []byte("test"),
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
	tags := []tag.FileTag{
		tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "test",
				Type: tag.USER | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[1],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "test2",
				Type: tag.TOPIC | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[1],
			Tag: tag.Tag{
				Word: "test3",
				Type: tag.ACTION | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[1],
			Owner: ownerids[1],
			Tag: tag.Tag{
				Word: "test4",
				Type: tag.ALLFILE,
			},
		},
	}
	t.Logf("fileids: %v", fileids)
	t.Logf("ownerids: %v", ownerids)
	t.Log("Upsert")
	err := tb.Upsert(tags...)
	if err != nil {
		t.Fatalf("Unable to Upsert tags: %s", err.Error())
	}

	t.Log("Get")
	{
		gtags, err := tb.Get(fileids[0], ownerids[0])
		if err != nil {
			t.Fatalf("Unable to Get tags: %s", err.Error())
		}
		if len(gtags) != 2 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("GetType")
	{
		gtags, err := tb.GetType(fileids[0], ownerids[0], tag.USER)
		if err != nil {
			t.Fatalf("Unable to GetType: %s", err.Error())
		}
		if len(gtags) != 1 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("GetAll")
	{
		gtags, err := tb.GetAll(tag.USER, ownerids[0])
		if err != nil {
			t.Fatalf("Unable to GetAll: %s", err.Error())
		}
		if len(gtags) != 1 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("SearchFiles")
	{
		stags, err := tb.SearchFiles(fileids, tag.FileTag{
			Tag: tag.Tag{
				Word: "test",
				Type: tag.ALLTYPES,
				Data: tag.Data{
					tag.SEARCH: map[string]interface{}{
						"regex": true,
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("Unable to Search Files: %s", err.Error())
		}
		if len(stags) != 2 {
			t.Fatalf("Incorrect Return: %v", stags)
		}
	}
}
