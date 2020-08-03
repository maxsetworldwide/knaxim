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
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := NewUser("testuser", "testtest", "test@test.test")

	if !user.Match(user) {
		t.Fatal("basic equality check failed")
	}

	if !user.GetLock().Valid(map[string]interface{}{
		"pass": "testtest",
	}) {
		t.Fatal("failed to unlock")
	}

	cookies := user.NewCookies(time.Now().Add(12*time.Hour), time.Now().Add(24*time.Hour))

	if user.GetID().String() != cookies[1].Value {
		t.Fatalf("incorrect cookie value")
	}

	testrequest := httptest.NewRequest("GET", "/test/test", &bytes.Buffer{})

	for _, c := range cookies {
		testrequest.AddCookie(c)
	}

	if !user.CheckCookie(testrequest) {
		t.Fatalf("Failed to validate cookies")
	}

	cookieOID, err := GetCookieUID(testrequest)
	if err != nil {
		t.Fatalf("unable to get oid: %s", err)
	}
	if !cookieOID.Equal(user.GetID()) {
		t.Fatalf("mismatched cookie id: %v", cookieOID)
	}
}
