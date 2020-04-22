package tag

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	message := "the cheese"
	tags, err := ExtractContentTags(strings.NewReader(message))
	if err != nil {
		t.Fatalf("unable to extract content tags: %s", err.Error())
	}
	if len(tags) != 2 {
		t.Fatalf("incorrect result: %v", tags)
	}
}
