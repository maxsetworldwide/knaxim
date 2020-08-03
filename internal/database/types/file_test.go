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

package types

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

func TestFileStore(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString("This is the file content")

	fs, err := NewFileStore(buf)
	if err != nil {
		t.Fatalf("error creating filestore: %s", err)
	}

	rdr, err := fs.Reader()
	if err != nil {
		t.Fatalf("unable to create reader: %s", err)
	}
	sb := new(strings.Builder)
	if _, err := io.Copy(sb, rdr); err != nil {
		t.Fatalf("unable to copy reader: %s", err)
	}

	if s := sb.String(); s != "This is the file content" {
		t.Fatalf("uncorrect read string: %s", s)
	}

	fs.Perr = &errors.Processing{
		Status:  444,
		Message: "all lasers",
	}
	fs2 := fs.Copy()
	if fs.ContentType != fs2.ContentType || fs.FileSize != fs2.FileSize || fs.Perr.Status != fs2.Perr.Status {
		t.Fatal("Failed to copy file store")
	}
}

func TestFile(t *testing.T) {
	fid := FileID{
		StoreID: StoreID{
			Hash:  10,
			Stamp: 10,
		},
		Stamp: []byte("test"),
	}
	file := &File{
		Permission: Permission{
			Own: NewUser("test", "test", "test"),
		},
		ID:   fid,
		Name: "testfile",
	}
	if !file.GetID().Equal(fid) {
		t.Fatalf("failed to get id")
	}
	fjson, err := file.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to MarshalJSON: %s", err)
	}
	fbson, err := file.MarshalBSON()
	if err != nil {
		t.Fatalf("failed to MarshalBSON: %s", err)
	}

	{
		filecopy := new(File)
		err := filecopy.UnmarshalJSON(fjson)
		if err != nil {
			t.Log(string(fjson))
			t.Fatalf("unable to decode file json: %s", err)
		}
		if filecopy.GetName() != "testfile" {
			t.Fatalf("incorrect file object from json: %v", filecopy)
		}
	}

	{
		filecopy := new(File)
		err := filecopy.UnmarshalBSON(fbson)
		if err != nil {
			t.Log(string(fbson))
			t.Fatalf("unable to decode file bson: %s", err)
		}
		if filecopy.GetName() != "testfile" {
			t.Fatalf("incorrect file object from bson: %v", filecopy)
		}
	}

	decoder := new(FileDecoder)
	err = decoder.UnmarshalJSON(fjson)
	if err != nil {
		t.Log(string(fjson))
		t.Fatalf("unable to decode file json: %s", err)
	}
	if !decoder.File().GetID().Equal(fid) {
		t.Fatalf("incorrect id decoded from json")
	}

	*decoder = FileDecoder{}
	err = decoder.UnmarshalBSON(fbson)
	if err != nil {
		t.Fatalf("unable to decode file bson")
	}
	if !decoder.File().GetID().Equal(fid) {
		t.Fatalf("incorrect id decoded from bson")
	}
}

func TestWeb(t *testing.T) {
	fid := FileID{
		StoreID: StoreID{
			Hash:  10,
			Stamp: 10,
		},
		Stamp: []byte("test"),
	}
	file := WebFile{
		File: File{
			Permission: Permission{
				Own: NewUser("test", "test", "test"),
			},
			ID:   fid,
			Name: "testfile",
		},
		URL: "test.test.test",
	}

	if !file.GetID().Equal(fid) {
		t.Fatalf("failed to get id")
	}
	fjson, err := file.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to MarshalJSON: %s", err)
	}
	fbson, err := file.MarshalBSON()
	if err != nil {
		t.Fatalf("failed to MarshalBSON: %s", err)
	}

	{
		filecopy := new(WebFile)
		err := filecopy.UnmarshalJSON(fjson)
		if err != nil {
			t.Log(string(fjson))
			t.Fatalf("unable to decode file json: %s", err)
		}
		if filecopy.GetName() != "testfile" {
			t.Fatalf("incorrect file object from json: %v", filecopy)
		}
	}

	{
		filecopy := new(WebFile)
		err := filecopy.UnmarshalBSON(fbson)
		if err != nil {
			t.Log(string(fbson))
			t.Fatalf("unable to decode file bson: %s", err)
		}
		if filecopy.GetName() != "testfile" {
			t.Fatalf("incorrect file object from bson: %v", filecopy)
		}
	}

	decoder := new(FileDecoder)
	err = decoder.UnmarshalJSON(fjson)
	if err != nil {
		t.Log(string(fjson))
		t.Fatalf("unable to decode file json: %s", err)
	}
	if !decoder.File().GetID().Equal(fid) {
		t.Fatalf("incorrect id decoded from json")
	}

	*decoder = FileDecoder{}
	err = decoder.UnmarshalBSON(fbson)
	if err != nil {
		t.Fatalf("unable to decode file bson")
	}
	if !decoder.File().GetID().Equal(fid) {
		t.Fatalf("incorrect id decoded from bson")
	}
}
