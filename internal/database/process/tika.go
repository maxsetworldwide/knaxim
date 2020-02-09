package process

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"github.com/google/go-tika/tika"
	"github.com/jdkato/prose/tokenize"
)

var stokenizer = tokenize.NewPunktSentenceTokenizer()

//TODO split on newline or any punc when token gets too large

func SentenceSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(bytes.TrimSpace(data)) == 0 {
		return 0, nil, nil
	}
	//TODO: extract Tokenize's rules to use here in order to optimize for Scanner Operations.
	sentences := stokenizer.Tokenize(string(data))
	if len(sentences) > 1 || (len(sentences) == 1 && atEOF) {
		offset := len(sentences[0])
		return offset, data[:offset], nil
	}
	if len(sentences) < 1 {
		return 0, nil, srverror.New(fmt.Errorf("Empty return from tokenizer, data = %s", string(data)), 500, "Database Error 011")
	}
	if len(data) > 512 {
		breakpoint := bytes.LastIndexByte(data, '\n') //last new line
		if breakpoint < 0 {
			breakpoint = bytes.LastIndexAny(data, " \t\r\f\n.!?;:)}]>,'\"`")
			if breakpoint < 0 {
				if len(data) < 2048 {
					return 0, nil, nil
				} else {
					return 2048, data[:2048], nil
				}
			}
		}
		return breakpoint + 1, data[:breakpoint+1], nil
	}

	return 0, nil, nil
}

type ContentExtractor tika.Client

func NewContentExtractor(httpClient *http.Client, urlString string) *ContentExtractor {
	return (*ContentExtractor)(tika.NewClient(httpClient, urlString))
}

func (ce *ContentExtractor) ExtractText(ctx context.Context, filecontent io.Reader) ([]database.ContentLine, error) {

	out := make([]database.ContentLine, 0, 128)

	text, err := ((*tika.Client)(ce)).Parse(ctx, filecontent)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error 012", "content parse error")
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(SentenceSplitter)
	count := 0
	for scanner.Scan() {
		out = append(out, database.ContentLine{
			Position: count,
			Content:  []string{scanner.Text()},
		})
		count++
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	//Generate database.ContentLines based on meta
	return out, nil
}

var xlsNewSheet = regexp.MustCompile("^Sheet[[:digit:]]+")

var csvRow = regexp.MustCompile("^.*([,\t].*)*$")
var csvSep = regexp.MustCompile("[,\t]")

// var csvRow = regexp.MustCompile("^.*(,.*)*$")

func (ce *ContentExtractor) ExtractCSV(ctx context.Context, filecontent io.Reader) ([]database.ContentLine, error) {

	out := make([]database.ContentLine, 0, 128)

	text, err := ((*tika.Client)(ce)).Parse(ctx, filecontent)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	count := 0
	page := 0
	for scanner.Scan() {
		line := scanner.Text()
		if xlsNewSheet.MatchString(line) {
			page++
		} else if csvRow.MatchString(line) {
			out = append(out, database.ContentLine{
				//PageNum:  page,
				Position: count,
				Content:  csvSep.Split(line, -1),
			})
			count++
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
