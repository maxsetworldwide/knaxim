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
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/brand"
)

// FileID uniquely identifies a file and its associated file store.
type FileID struct {
	StoreID
	Stamp []byte `bson:"stamp"`
}

// NewFileID generates a new file id for a particular file store.
func NewFileID(st StoreID) FileID {
	var n FileID
	n.StoreID = st
	n.Stamp = []byte{brand.Next()}
	return n
}

// String returns a base64 encoding of the FileID as a string
func (f FileID) String() string {
	build := new(strings.Builder)
	encoder := base64.NewEncoder(base64.RawURLEncoding, build)

	binary.Write(encoder, binary.LittleEndian, f.StoreID.Hash)
	binary.Write(encoder, binary.LittleEndian, f.StoreID.Stamp)
	encoder.Write(f.Stamp)
	encoder.Close()

	return build.String()
}

// Mutate returns a new FileID that is associated with the same StoreID
func (f FileID) Mutate() FileID {
	f.Stamp = append(f.Stamp, brand.Next())
	return f
}

// DecodeFileID produces a FileID from a base64 encoded string, inverse of String().
func DecodeFileID(h string) (fid FileID, err error) {
	defer func() {
		if r := recover(); r != nil {
			fid = FileID{}
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("Failed to decode FileID: %v", r)
			}
		}
	}()
	buf, err := base64.RawURLEncoding.DecodeString(h)
	if err != nil {
		return FileID{}, err
	}
	var n FileID
	n.StoreID = decodeStoreID(buf[:6])
	n.Stamp = buf[6:]
	return n, nil
}

// Equal is true if the provided FileID is the same value as the FileID
func (f FileID) Equal(oth FileID) bool {
	return f.StoreID.Equal(oth.StoreID) && bytes.Equal(f.Stamp, oth.Stamp)
}

// MarshalJSON returns the string representation of the FileID in json format
func (f FileID) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON decodes a json string value into a FileID, using DecodeFileID
func (f *FileID) UnmarshalJSON(b []byte) error {
	var fstr string
	err := json.Unmarshal(b, &fstr)
	if err != nil {
		return err
	}
	*f, err = DecodeFileID(fstr)
	return err
}
