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

// basing off of database/file.go

package types

import (
	"bytes"
	"compress/gzip"
	"io"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

// ViewStore represents the PDF representation of a File
type ViewStore struct {
	ID      StoreID `json:"id" bson:"id"`
	Content []byte  `json:"content" bson:"-"`
}

// NewViewStore builds ViewStore from a StoreID and Reader of
// pdf representation
func NewViewStore(id StoreID, r io.Reader) (*ViewStore, error) {
	store := new(ViewStore)

	contentBuf := new(bytes.Buffer)
	gzWrite := gzip.NewWriter(contentBuf)

	var err error
	if _, err = io.Copy(gzWrite, r); err != nil {
		return nil, err
	}

	if err = gzWrite.Close(); err != nil {
		return nil, srverror.New(err, 500, "Error V1")
	}

	store.Content = contentBuf.Bytes()
	store.ID = id
	return store, nil
}

// Reader returns a reader of the content of a pdf representation of the file
func (vs *ViewStore) Reader() (io.Reader, error) {
	buf := bytes.NewReader(vs.Content)
	out, err := gzip.NewReader(buf)
	if err != nil {
		srverror.New(err, 500, "Error V2", "file reading error")
	}
	return out, err
}
