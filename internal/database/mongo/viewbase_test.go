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
)

func TestViewbase(t *testing.T) {
	t.Parallel()
	var vb *Viewbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestView"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatalf("Unable to init database: %s", err)
		}
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		mdb, err := db.Connect(methodtesting)
		if err != nil {
			t.Fatalf("Unable to connect to database: %s", err.Error())
		}
		defer mdb.Close(methodtesting)
		vb = mdb.View().(*Viewbase)
	}
	{
		contentString := "View version of a file FFFFF*****"
		inputVS := &types.ViewStore{
			ID: types.StoreID{
				Hash:  12345,
				Stamp: 678,
			},
			Content: []byte(contentString),
		}
		t.Run("Insert", func(t *testing.T) {
			err := vb.Insert(inputVS)
			if err != nil {
				t.Fatalf("error inserting: %s", err)
			}
		})
		t.Run("Get", func(t *testing.T) {
			result, err := vb.Get(inputVS.ID)
			if err != nil {
				t.Fatalf("error getting: %s", err)
			}
			if !inputVS.ID.Equal(result.ID) || !bytes.Equal(inputVS.Content, result.Content) {
				t.Errorf("Did not get correct view store:\ngot: %+#v\nexpected: %+#v\n", result, inputVS)
			}
		})
	}
}
