package decode

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/types"
)

var tikapath = os.Getenv("TIKA_PATH")

var gotenburgpath = os.Getenv("GOTENBURG_PATH")

func init() {
	if len(tikapath) == 0 {
		tikapath = "http://localhost:9998"
	}
	if len(gotenburgpath) == 0 {
		gotenburgpath = "http://localhost:3000"
	}
}

func TestRead(t *testing.T) {
	testctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var db database.Database
	db = &memory.Database{}
	db.Init(testctx, true)
	fs, err := types.NewFileStore(strings.NewReader("This is the test file. I want this to pass."))
	if err != nil {
		t.Fatalf("unable to create File Store: %s", err.Error())
	}
	fs.ContentType = "text/plain"
	sb := db.Store()
	if fs.ID, err = sb.Reserve(fs.ID); err != nil {
		t.Fatalf("unable to reserve id for filestore: %s", err.Error())
	}
	if err := sb.Insert(fs); err != nil {
		t.Fatalf("unable to insert filestore: %s", err.Error())
	}
	sb.Close(testctx)
	lock := make(chan struct{}, 1)
	lock <- struct{}{}
	testctx = context.WithValue(testctx, PROCESSING, lock)
	testctx = context.WithValue(testctx, TIMEOUT, time.Minute)
	Read(testctx, cancel, "test.txt", fs, db, tikapath, gotenburgpath)
	databasejson, err := json.MarshalIndent(db, "", "\t")
	if err != nil {
		t.Fatalf("unable to produce database output: %s", err.Error())
	}
	t.Logf("state of db: %s\n", string(databasejson))
	select {
	case <-lock:
		t.Logf("lock released")
	default:
		t.Errorf("lock held")
	}
}
