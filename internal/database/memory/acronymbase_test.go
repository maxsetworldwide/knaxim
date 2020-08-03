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

package memory

import "testing"

func TestAcronym(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	ab := DB.Acronym()
	defer DB.Close(nil)
	t.Parallel()

	t.Log("Acronym Put")
	err := ab.Put("t", "test")
	if err != nil {
		t.Fatalf("Unable to put acronym: %s", err)
	}

	t.Log("Acronym Get")
	matches, err := ab.Get("t")
	if err != nil {
		t.Fatalf("Unable to get acronym: %s", err)
	}
	if len(matches) != 1 || matches[0] != "test" {
		t.Fatalf("incorrect matches: %v", matches)
	}
}
