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
	"bytes"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

var contentString = "View version of a file for the memory database"

var inputVS = &types.ViewStore{
	ID: types.StoreID{
		Hash:  98765,
		Stamp: 4321,
	},
	Content: []byte(contentString),
}

func TestViewbase(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	vb := DB.View()
	defer DB.Close(nil)
	t.Parallel()

	t.Log("View Insert")
	err := vb.Insert(inputVS)
	if err != nil {
		t.Fatalf("error inserting view: %s\n", err)
	}
	t.Log("View Get")
	result, err := vb.Get(inputVS.ID)
	if err != nil {
		t.Fatalf("error getting view: %s\n", err)
	}
	if !inputVS.ID.Equal(result.ID) || !bytes.Equal(inputVS.Content, result.Content) {
		t.Fatalf("Did not get correct view store:\ngot: %+#v\nexpected: %+#v\n", result, inputVS)
	}
}
