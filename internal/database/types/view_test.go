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
)

func TestViewStore(t *testing.T) {
	contentString := "This is the view content! It's like the file store content, but it should be a PDF version of the file."
	inBytes := []byte(contentString)

	mockStoreID := StoreID{
		Hash:  12345,
		Stamp: 6789,
	}
	vs, err := NewViewStore(mockStoreID, bytes.NewReader(inBytes))
	if err != nil {
		t.Fatalf("error creating viewstore: %s", err)
	}

	rdr, err := vs.Reader()
	if err != nil {
		t.Fatalf("unable to create reader from viewstore: %s", err)
	}

	sb := new(strings.Builder)
	if _, err := io.Copy(sb, rdr); err != nil {
		t.Fatalf("unable to copy from reader: %s", err)
	}

	if s := sb.String(); s != contentString {
		t.Fatalf("incorrect read string: expected '%s', got '%s'", contentString, s)
	}
}
