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
	"context"
	"testing"
	"time"
)

func TestAcronym(t *testing.T) {
	t.Parallel()
	var ab *Acronymbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestAcronym"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to Init database", err)
		}
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		mdb, err := db.Connect(methodtesting)
		if err != nil {
			t.Fatalf("Unable to connect to database: %s", err.Error())
		}
		defer mdb.Close(methodtesting)
		ab = mdb.Acronym().(*Acronymbase)
	}
	t.Run("Put", func(t *testing.T) {
		err := ab.Put("ab", "Acronymbase")
		if err != nil {
			t.Fatal("Unable to add acronym, ", err)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get", func(t *testing.T) {
		result, err := ab.Get("ab")
		if err != nil {
			t.Fatal("Unable to get acronym, ", err)
		}
		if len(result) != 1 || result[0] != "Acronymbase" {
			t.Fatal("incorrect result: ", result)
		}
	})
}
