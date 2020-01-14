package database

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"git.maxset.io/server/knaxim/database/filehash"
	"git.maxset.io/server/knaxim/srverror"

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

type ContentLine struct {
	ID filehash.StoreID `bson:"id"`
	//PageNum  int              `bson:"pagenum"`
	Position int      `bson:"position"`
	Content  []string `bson:"content"`
}

func (ce *ContentExtractor) ExtractText(ctx context.Context, filecontent io.Reader) ([]ContentLine, error) {

	out := make([]ContentLine, 0, 128)

	text, err := ((*tika.Client)(ce)).Parse(ctx, filecontent)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error 012", "content parse error")
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(SentenceSplitter)
	count := 0
	for scanner.Scan() {
		out = append(out, ContentLine{
			Position: count,
			Content:  []string{scanner.Text()},
		})
		count++
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	//Generate ContentLines based on meta
	return out, nil
}

var xlsNewSheet = regexp.MustCompile("^Sheet[[:digit:]]+")

var csvRow = regexp.MustCompile("^(\t.*)+$")

func (ce *ContentExtractor) ExtractCSV(ctx context.Context, filecontent io.Reader) ([]ContentLine, error) {

	out := make([]ContentLine, 0, 128)

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
			out = append(out, ContentLine{
				//PageNum:  page,
				Position: count,
				Content:  strings.Split(line, "\t")[1:],
			})
			count++
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func NewContentReader(lines []ContentLine) (result io.Reader, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = fmt.Errorf("Building Content Reader %v", v)
			}
		}
	}()
	out := make([]ContentLine, len(lines))
	copy(out, lines)
	for i := range out {
		for target := out[i].Position; target != i; target = out[i].Position {
			if target == out[target].Position {
				panic("double position")
			}
			out[i], out[target] = out[target], out[i]
		}
	}
	linereaders := make([]io.Reader, 0, len(out))
	for _, line := range out {
		linereaders = append(linereaders, strings.NewReader(strings.Join(line.Content, ", ")))
	}
	return io.MultiReader(linereaders...), nil
}
