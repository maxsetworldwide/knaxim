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
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	fid := FileID{
		StoreID: StoreID{
			Hash:  15,
			Stamp: 16,
		},
		Stamp: []byte("test"),
	}
	jsonbytes, err := json.Marshal(fid)
	if err != nil {
		t.Fatal("Unable to encode FileID: ", err)
	}
	var unmarshaled FileID
	err = json.Unmarshal(jsonbytes, &unmarshaled)
	if err != nil {
		t.Fatal("Unable to decode FileID: ", err)
	}
	if !fid.Equal(unmarshaled) {
		t.Fatalf("fid mismatched: %+#v", unmarshaled)
	}
}
