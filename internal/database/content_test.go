package database

import (
	"io"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/filehash"
)

func TestContent(t *testing.T) {
	t.Parallel()
	sid := filehash.StoreID{
		Hash:  10,
		Stamp: 10,
	}
	lines := []ContentLine{
		ContentLine{
			ID:       sid,
			Position: 0,
			Content:  []string{"a"},
		},
		ContentLine{
			ID:       sid,
			Position: 2,
			Content:  []string{"c"},
		},
		ContentLine{
			ID:       sid,
			Position: 1,
			Content:  []string{"b"},
		},
	}

	rdr, err := NewContentReader(lines)
	if err != nil {
		t.Fatalf("Failed to create Content Reader: %s", err)
	}

	sb := new(strings.Builder)
	if _, err := io.Copy(sb, rdr); err != nil {
		t.Fatalf("Unable to Read: %s", err)
	}

	if s := sb.String(); s != "abc" {
		t.Fatalf("Incorrect resulting string: %s", s)
	}
}
