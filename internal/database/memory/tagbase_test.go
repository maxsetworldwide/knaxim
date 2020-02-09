package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/tag"
)

func TestTag(t *testing.T) {
	defer testingComplete.Done()
	tb := DB.Tag(nil)
	defer tb.Close(nil)
	t.Parallel()

	filetag := tag.Tag{
		Word: "test",
		Type: tag.USER,
		Data: tag.Data{
			tag.USER: map[string]string{
				"hello": "world",
			},
		},
	}

	t.Log("Upsert File")
	err := tb.UpsertFile(fid, filetag)
	if err != nil {
		t.Fatalf("failed to UpsertFile: %s", err)
	}

	storetag := tag.Tag{
		Word: "another",
		Type: tag.CONTENT,
		Data: tag.Data{
			tag.CONTENT: map[string]string{
				"data": "base",
			},
		},
	}

	t.Log("Upsert Store")
	err = tb.UpsertStore(sid, storetag)
	if err != nil {
		t.Fatalf("failed to upsert store tag: %s", err)
	}

	t.Log("FileTags")
	matched, err := tb.FileTags(fid)
	if err != nil {
		t.Fatalf("failed to get File Tags: %s", err)
	}
	if len(matched) != 2 || matched["test"] == nil || matched["another"] == nil {
		t.Fatalf("incorrect matches: %v", matched)
	}

	t.Log("GetFiles")
	fids, sids, err := tb.GetFiles([]tag.Tag{
		tag.Tag{
			Word: "test",
			Type: tag.USER,
		},
	})
	if err != nil {
		t.Fatalf("failed to GetFiles: %s", err)
	}
	if len(fids) != 1 || !fids[0].Equal(fid) || len(sids) != 1 || !sids[0].Equal(sid) {
		t.Fatalf("incorrect return from GetFiles: %v, %v", fids, sids)
	}

	t.Log("Search Data")
	matches, err := tb.SearchData(tag.USER, tag.Data{
		tag.USER: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		t.Fatalf("unable to SearchData: %s", err)
	}
	if len(matches) != 1 || matches[0].Word != "test" {
		t.Fatalf("incorrect return of SearchData: %v", matches)
	}
}
