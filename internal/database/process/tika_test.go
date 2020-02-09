package process

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSentSplit(t *testing.T) {
	teststr := "This is the First Sentence. And here is another. What are you talking about?"

	scn := bufio.NewScanner(strings.NewReader(teststr))
	scn.Split(SentenceSplitter)

	if !scn.Scan() {
		t.Fatal("unable to get first sentence")
	}
	if line := scn.Text(); line != "This is the First Sentence." {
		t.Fatalf("incorrect first line: %s", line)
	}
	if !scn.Scan() {
		t.Fatal("unable to get second sentence")
	}
	if line := scn.Text(); line != " And here is another." {
		t.Fatalf("incorrect second line: %s", line)
	}
	if !scn.Scan() {
		t.Fatal("unable to get third sentence")
	}
	if line := scn.Text(); line != " What are you talking about?" {
		t.Fatalf("incorrect third line: %s", line)
	}
	if scn.Scan() {
		t.Fatal("did not finish scanning")
	}
}

func TestContentExtractor(t *testing.T) {
	testtxt := "This is the test content of a text file. This is to test tika content extraction. If this test is failing, then make sure that there is a tika instance running and that TIKA_PATH is correct. And one more sentence for good luck."
	tikapath := os.Getenv("TIKA_PATH")
	if len(tikapath) == 0 {
		tikapath = "http://localhost:9998"
	}
	t.Logf("tikapath = %s", tikapath)

	contex := NewContentExtractor(nil, tikapath)

	testctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	lines, err := contex.ExtractText(testctx, strings.NewReader(testtxt))
	if err != nil {
		t.Fatalf("Failed to ExtractText: %s", err)
	}
	if len(lines) != 4 {
		t.Fatalf("incorrect lines: %+v", lines)
	}

	testcsv := `Name,Profession,Birth month
Devon,Programmer,September
Drew,Paramedic,August`
	testcsv += "\nPiper\tSnuggler\tDecember"

	lines, err = contex.ExtractCSV(testctx, strings.NewReader(testcsv))
	if err != nil {
		t.Fatalf("Failed to ExtractCSV: %s", err)
	}
	if len(lines) != 4 || len(lines[0].Content) != 3 {
		t.Fatalf("incorrect lines: %v", lines)
	}
}
