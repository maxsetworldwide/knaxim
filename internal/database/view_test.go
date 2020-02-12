package database

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestViewStore(t *testing.T) {
	buf := new(bytes.Buffer)
	contentString := "This is the view content! It's like the file store content, but it should be a PDF version of the file."
	buf.WriteString(contentString)

	vs, err := NewViewStore(buf)
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
