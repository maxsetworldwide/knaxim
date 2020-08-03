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

package mongo

import (
	"bytes"
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

func TestStorebase(t *testing.T) {
	t.Parallel()
	var sb *Storebase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestStore"
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
		sb = mdb.Store().(*Storebase)
	}
	{
		input := types.StoreID{
			Hash:  24098,
			Stamp: 123,
		}
		t.Run("Reserve=basic", func(t *testing.T) {
			out, err := sb.Reserve(input)
			if err != nil {
				t.Error("error reserving", err)
			}
			if out.Hash != input.Hash || out.Stamp != input.Stamp {
				t.Error("return mismatch", out)
			}
		})
		t.Run("Reserve=mutate", func(t *testing.T) {
			out, err := sb.Reserve(input)
			if err != nil {
				t.Error("error reserving", err)
			}
			if out.Hash != input.Hash || out.Stamp != input.Mutate().Stamp {
				t.Error("return mismatch", out)
			}
		})
	}
	{
		input := &types.FileStore{
			ID: types.StoreID{
				Hash:  24098,
				Stamp: 123,
			},
			Content:     []byte("Here is a test file.#$%!1234t##g1"),
			ContentType: "text",
			FileSize:    33,
		}
		t.Run("Insert", func(t *testing.T) {
			err := sb.Insert(input)
			if err != nil {
				t.Fatal("error inserting", err)
			}
		})
		t.Run("Get", func(t *testing.T) {
			out, err := sb.Get(input.ID)
			if err != nil {
				t.Fatal("error getting", err)
			}
			if !input.ID.Equal(out.ID) ||
				!bytes.Equal(input.Content, out.Content) ||
				input.ContentType != out.ContentType ||
				input.FileSize != out.FileSize {
				t.Error("did not get correct file store", out)
			}
		})
		t.Run("MatchHash", func(t *testing.T) {
			out, err := sb.MatchHash(input.ID.Hash)
			if err != nil {
				t.Fatal("error match hash", err)
			}
			if !input.ID.Equal(out[0].ID) ||
				!bytes.Equal(input.Content, out[0].Content) ||
				input.ContentType != out[0].ContentType ||
				input.FileSize != out[0].FileSize {
				t.Error("did not get correct file store", out[0])
			}
		})
		t.Run("Update", func(t *testing.T) {
			input.Perr = &errors.Processing{
				Status:  420,
				Message: "Hey, You see this",
			}
			err := sb.UpdateMeta(input)
			if err != nil {
				t.Fatalf("unable to UpdateMeta: %s", err)
			}
		})
	}
}
