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
)

var fid = types.FileID{
	StoreID: sid,
	Stamp:   []byte{'a', 'b'},
}

func TestFiles(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	fb := DB.File()
	defer DB.Close(nil)
	t.Parallel()

	storedid, err := fb.Reserve(fid)
	if err != nil {
		t.Fatalf("unable to reserve id: %s", err)
	}
	if !storedid.Equal(fid) {
		t.Fatalf("incorrect storeid reserved: %v", storedid)
	}

	file := &types.File{
		Permission: types.Permission{
			Own: test1,
		},
		ID:   fid,
		Name: "TestFile",
	}

	t.Log("Insert")
	err = fb.Insert(file)
	if err != nil {
		t.Fatalf("failed to insert file: %s", err)
	}

	t.Log("Get")
	gotten, err := fb.Get(fid)
	if err != nil {
		t.Fatalf("failed to find file: %s", err)
	}
	if !gotten.GetID().Equal(fid) {
		t.Fatalf("incorrect gotten file: %s", err)
	}

	t.Log("GetAll")
	allgot, err := fb.GetAll(fid)
	if err != nil {
		t.Fatalf("failed to find files: %s", err)
	}
	if len(allgot) != 1 || !allgot[0].GetID().Equal(fid) {
		t.Fatalf("incorrect GetAll: %v", allgot)
	}

	t.Log("Update")
	file.SetPerm(test2, "view", true)
	err = fb.Update(file)
	if err != nil {
		t.Fatalf("failed to update file: %s", err)
	}

	t.Log("Get Owned")
	owned, err := fb.GetOwned(test1.GetID())
	if err != nil {
		t.Fatalf("failed to GetOwned: %s", err)
	}
	if len(owned) != 1 || !owned[0].GetID().Equal(fid) {
		t.Fatalf("incorrect return from owned: %v", owned)
	}

	t.Log("Get PermKey")
	shared, err := fb.GetPermKey(test2.GetID(), "view")
	if err != nil {
		t.Fatalf("failed to GetPermKey: %s", err)
	}
	if len(shared) != 1 || !shared[0].GetID().Equal(fid) {
		t.Fatalf("incorrect return from shared: %v", shared)
	}

	t.Log("MatchStore")
	matched, err := fb.MatchStore(test1.GetID(), []types.StoreID{sid})
	if err != nil {
		t.Fatalf("unable to match sid: %s", err)
	}
	if len(matched) != 1 || !matched[0].GetID().Equal(fid) {
		t.Fatalf("incorrect returned matched: %v", matched)
	}

	t.Log("Count")
	count, err := fb.Count(test1.GetID())
	if err != nil {
		t.Fatalf("unable to count files: %s", err.Error())
	}
	if count != 1 {
		t.Fatalf("incorrect file count: %d", count)
	}

	t.Log("Remove")
	err = fb.Remove(fid)
	if err != nil {
		t.Fatalf("failed to remove: %s", err)
	}

	expectnil, err := fb.Get(fid)
	if err == nil || expectnil != nil {
		t.Fatalf("found removed file: %v, %s", expectnil, err)
	}
	t.Log("Remove End")
}
