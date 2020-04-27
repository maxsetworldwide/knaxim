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

func TestName(t *testing.T) {
	name := "the_File.txt"
	tags, err := BuildNameTags(name)
	if err != nil {
		t.Fatalf("unable to build name tags: %s", err.Error())
	}
	if len(tags) != 4 && tags[0].Word == name {
		t.Fatalf("incorrect result: %v", tags)
	}
}
