package database

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestViewStore(t *testing.T) {
	contentString := "This is the view content! It's like the file store content, but it should be a PDF version of the file."
	inBytes := []byte(contentString)

	mockStoreID := filehash.StoreID{
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
