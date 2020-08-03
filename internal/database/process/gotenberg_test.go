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

/*
this test requires Gotenberg and Tika to be running, and the TIKA_PATH and
GOTENBERG_PATH env variables to both be set to the correct URLs of the services
*/

package process

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func extractContent(t *testing.T, content []byte) (out []string, err error) {
	tikapath := os.Getenv("TIKA_PATH")
	if len(tikapath) == 0 {
		tikapath = "http://localhost:9998"
	}
	t.Logf("tikapath = %s", tikapath)

	extractor := NewContentExtractor(nil, tikapath)
	testctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	extracted, err := extractor.ExtractText(testctx, bytes.NewReader(content))
	if err != nil {
		return
	}
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
	inBytes := []byte(testContent)
	out, err := client.ConvertOffice("test.txt", inBytes)
	if err != nil {
		t.Fatal("Conversion fail:", err)
	}
	lines, err := extractContent(t, out)
	if err != nil {
		t.Fatal("Extraction fail:", err)
	}
	found := false
	for i := 0; i < len(lines) && !found; i++ {
		found = strings.Index(lines[i], testContent) != -1
	}
	if !found {
		t.Fatalf("Test text did not appear in converted pdf:\nLooking for '%s', but file content was:%+#v", testContent, lines)
	}
}
