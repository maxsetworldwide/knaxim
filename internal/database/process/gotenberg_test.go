package process

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func extractContent(t *testing.T, content *bytes.Buffer) (out []string, err error) {
	tikapath := os.Getenv("TIKA_PATH")
	if len(tikapath) == 0 {
		tikapath = "http://localhost:9998"
	}
	t.Logf("tikapath = %s", tikapath)

	extractor := NewContentExtractor(nil, tikapath)
	testctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	extracted, err := extractor.ExtractText(testctx, content)
	if err != nil {
		return
	}
	// out = lines.Content
	for _, curr := range extracted {
		for _, line := range curr.Content {
			out = append(out, line)
		}
	}
	return
}

func TestTextConversion(t *testing.T) {
	url := os.Getenv("GOTENBERG_PATH")
	if len(url) == 0 {
		url = "http://localhost:3000"
	}
	client := NewFileConverter(url)
	testContent := "Test content!!"
	var inBytes bytes.Buffer
	inBytes.Write([]byte(testContent))
	out, err := client.ConvertOffice("test.txt", &inBytes)
	if err != nil {
		t.Fatal("Conversion fail:", err)
	}
	lines, err := extractContent(t, out)
	if err != nil {
		t.Fatal("Extraction fail:", err)
	}
	t.Logf("lines out:%+#v", lines)
	found := false
	for i := 0; i < len(lines) && !found; i++ {
		found = strings.Index(lines[i], testContent) != -1
	}
	if !found {
		t.Fatalf("Test text did not appear in converted pdf:\nLooking for '%s', but file content was:%+#v", testContent, lines)
	}
}
