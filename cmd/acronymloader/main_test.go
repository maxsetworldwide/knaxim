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

package main

import (
	"context"
	"os"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/mongo"
)

func TestCli(t *testing.T) {
	db := new(mongo.Acronymbase)
	timeoutctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	db.URI = *databaseURI
	db.DBName = *databaseName
	db.CollNames = make(map[string]string)
	db.CollNames["acronym"] = *collectionName
	err := db.Init(timeoutctx, true)
	if err != nil {
		t.Fatal("unable to init database ", err)
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal("unable to create Pipe ", err)
	}
	os.Stdin = pr
	pw.Write([]byte("ab,acronymbase"))
	pw.Close()
	main()
}
