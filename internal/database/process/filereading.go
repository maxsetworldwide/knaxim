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

package process

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// InjestFile builds a file and file store from data and adds to database
func InjestFile(ctx context.Context, file types.FileI, contenttype string, stream io.Reader, dbconfig database.Database) (fs *types.FileStore, err error) {
	defer func() {
		if r := recover(); r != nil {
			fs = nil
			switch v := r.(type) {
			case srverror.Error:
				err = v
			case error:
				err = srverror.New(v, 500, "Error P1", "unable to process input")
			default:
				err = srverror.New(fmt.Errorf("Error Injecting File: %+#v", v), 500, "Error P2")
			}
		}
	}()
	fs, err = types.NewFileStore(stream)
	if err != nil {
		panic(err)
	}
	fs.ContentType = contenttype
	db, err := dbconfig.Connect(ctx)
	if err != nil {
		panic(srverror.New(err, 500, "Error P3", "Unable to connect to the database"))
	}
	defer db.Close(ctx)

	{
		ownerbase := db.Owner()
		currentspace, err := ownerbase.GetSpace(file.GetOwner().GetID())
		if err != nil {
			panic(err)
		}
		totalspace, err := ownerbase.GetTotalSpace(file.GetOwner().GetID())
		if err != nil {
			panic(err)
		}
		if currentspace+fs.FileSize > totalspace {
			panic(srverror.Basic(462, "No Space, Delete Files and empty trash to free space"))
		}
	}
	{
		sb := db.Store()
		matches, err := sb.MatchHash(fs.ID.Hash)
		if err != nil {
			panic(err)
		}
		var matched bool
		for _, m := range matches {
			if bytes.Equal(fs.Content, m.Content) {
				fs = m
				matched = true
				break
			}
		}
		if !matched {
			fs.ID, err = sb.Reserve(fs.ID)
			if err != nil {
				panic(err)
			}
			err = sb.Insert(fs)
			if err != nil {
				panic(err)
			}
		}
	}
	{
		fb := db.File()
		tempID, err := fb.Reserve(types.NewFileID(fs.ID))
		if err != nil {
			panic(err)
		}
		file.SetID(tempID)
		err = fb.Insert(file)
		if err != nil {
			panic(err)
		}
	}
	return fs, nil
}
