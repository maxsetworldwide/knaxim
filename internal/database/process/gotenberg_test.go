/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
