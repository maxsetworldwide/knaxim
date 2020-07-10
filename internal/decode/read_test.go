/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
