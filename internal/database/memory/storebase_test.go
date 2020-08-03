// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

var sid = types.StoreID{
	Hash:  42,
	Stamp: 42,
}

func fillstores(db *Database) {
	db.Stores[sid.String()] = &types.FileStore{
		ID:          sid,
		Content:     []byte("placeholder"),
		ContentType: "test",
		FileSize:    42,
	}
}

func TestStore(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	sb := DB.Store()
	defer DB.Close(nil)
	t.Parallel()

	var sid = types.StoreID{
		Hash:  10,
		Stamp: 10,
	}

	fs := &types.FileStore{
		ID:          sid,
		Content:     []byte("this is the content"),
		ContentType: "testfile",
		FileSize:    1234,
	}

	t.Log("Reserve")
	savedid, err := sb.Reserve(sid)
	if err != nil {
		t.Fatalf("Unable to reserve id: %s", err)
	}
	if !savedid.Equal(sid) {
		t.Fatalf("wrong id saved: %v", savedid)
	}

	t.Log("Insert")
	err = sb.Insert(fs)
	if err != nil {
		t.Fatalf("unable to insert: %s", err)
	}

	t.Log("Get")
	gotten, err := sb.Get(sid)
	if err != nil {
		t.Fatalf("unable to get filestore: %s", err)
	}
	if !gotten.ID.Equal(sid) {
		t.Fatalf("incorrect gotten filestore: %v", gotten)
	}

	t.Log("MatchHash")
	matched, err := sb.MatchHash(10)
	if err != nil {
		t.Fatalf("unable to match hash: %s", err)
	}
	if len(matched) != 1 || !matched[0].ID.Equal(sid) {
		t.Fatalf("incorrect matches: %v", matched)
	}

	t.Log("UpdateMeta")
	fs.Perr = &errors.Processing{
		Status:  420,
		Message: "hello",
	}
	err = sb.UpdateMeta(fs)
	if err != nil {
		t.Fatalf("Failed to UpdateMeta: %s", err)
	}
	if fs2, _ := sb.Get(sid); fs2.Perr == nil || fs2.Perr.Status != 420 {
		t.Fatalf("file store not updated: %+v", fs2)
	}
}
