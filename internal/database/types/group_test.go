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

import "testing"

func TestGroup(t *testing.T) {
	testuser := NewUser("testuser", "testtest", "test@test.test")
	testuser2 := NewUser("user2", "testest", "test@test.test")
	group := NewGroup("testgroup", testuser)
	group.AddMember(testuser2)

	if !group.Match(testuser) {
		t.Fatal("owner does not match group")
	}
	if !group.Match(testuser2) {
		t.Fatalf("member does not match group")
	}

	{
		members := group.GetMembers()
		if len(members) != 1 || !members[0].Equal(testuser2) {
			t.Fatal("incorrect member list")
		}
	}
	{
		gjson, err := group.MarshalJSON()
		if err != nil {
			t.Fatalf("unable to MarshalJSON group: %s", err)
		}
		ng := new(Group)
		err = ng.UnmarshalJSON(gjson)
		if err != nil {
			t.Fatalf("unable to UnmarshalJSON group: %s", err)
		}
		if ng.GetName() != group.GetName() {
			t.Fatal("incorrect decoded group")
		}
	}
	{
		gbson, err := group.MarshalBSON()
		if err != nil {
			t.Fatalf("unable to MarshalBSON group: %s", err)
		}
		ng := new(Group)
		err = ng.UnmarshalBSON(gbson)
		if err != nil {
			t.Fatalf("unable to UnmarshalBSON group: %s", err)
		}
		if ng.GetName() != group.GetName() {
			t.Fatal("incorrect decoded group")
		}
	}
}
