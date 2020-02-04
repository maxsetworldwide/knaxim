package database

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

	cookie_oid, err := GetCookieUid(testrequest)
	if err != nil {
		t.Fatalf("unable to get oid: %s", err)
	}
	if !cookie_oid.Equal(user.GetID()) {
		t.Fatalf("mismatched cookie id: %v", cookie_oid)
	}
}
