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

package process_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/memory"
	. "git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/types"
)

func TestInjustFile(t *testing.T) {
	var db = &memory.Database{}

	var content = `This is an example text.
  How do you know this is correct. Well you got this text back.
  Good Bye.`

	var testOwner = &types.User{
		ID: types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'a', 'b', 'c'},
			Stamp:       []byte("test"),
		},
		Name: "testuser",
	}
	initctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err := db.Init(initctx, true)
	if err != nil {
		t.Fatal("unable to init database", err)
	}
	injestctx, cancel2 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel2()
	file := &types.File{
		Name: "test.txt",
	}
	_, err = db.Owner().Reserve(testOwner.GetID(), testOwner.GetName())
	if err != nil {
		t.Fatal("unable to reserve testOwner:", err)
	}
	err = db.Owner().Insert(testOwner)
	if err != nil {
		t.Fatalf("Failed to insert test Owner")
	}
	file.Own = testOwner
	_, err = InjestFile(injestctx, file, "content/txt", strings.NewReader(content), db)
	if err != nil {
		t.Fatal("injest failed", err)
	}

	// generateContentTags(injestctx, fs, db)
}
