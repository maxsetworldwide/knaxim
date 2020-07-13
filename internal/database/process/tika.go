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

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"github.com/google/go-tika/tika"
	"github.com/jdkato/prose/tokenize"
)

var stokenizer = tokenize.NewPunktSentenceTokenizer()

//TODO split on newline or any punc when token gets too large

// SentenceSplitter implements the scanner split func type for splitting strings into sentences
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
		return 0, nil, srverror.New(fmt.Errorf("Empty return from tokenizer, data = %s", string(data)), 500, "Error 011")
	}
	if len(data) > 512 {
		breakpoint := bytes.LastIndexByte(data, '\n') //last new line
		if breakpoint < 0 {
			breakpoint = bytes.LastIndexAny(data, " \t\r\f\n.!?;:)}]>,'\"`")
			if breakpoint < 0 {
				if len(data) < 2048 {
					return 0, nil, nil
				}
				return 2048, data[:2048], nil
			}
		}
		return breakpoint + 1, data[:breakpoint+1], nil
	}

	return 0, nil, nil
}

// ContentExtractor is a connector to a tika server and build content lines from file streams
type ContentExtractor tika.Client

// NewContentExtractor connects to a tika server at a given address
func NewContentExtractor(httpClient *http.Client, urlString string) *ContentExtractor {
	return (*ContentExtractor)(tika.NewClient(httpClient, urlString))
}

// ExtractText process a file stream assuming it is some kind of text file
func (ce *ContentExtractor) ExtractText(ctx context.Context, filecontent io.Reader) ([]types.ContentLine, error) {

	out := make([]types.ContentLine, 0, 128)

	text, err := ((*tika.Client)(ce)).Parse(ctx, filecontent)
	if err != nil {
		return nil, srverror.New(err, 500, "Error 012", "content parse error")
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(SentenceSplitter)
	count := 0
	for scanner.Scan() {
		out = append(out, types.ContentLine{
			Position: count,
			Content:  []string{scanner.Text()},
		})
		count++
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	//Generate types.ContentLines based on meta
	return out, nil
}

var xlsNewSheet = regexp.MustCompile("^Sheet[[:digit:]]+")

var csvRow = regexp.MustCompile("^.*([,\t].*)*$")
var csvSep = regexp.MustCompile("[,\t]")

// var csvRow = regexp.MustCompile("^.*(,.*)*$")

// ExtractCSV process a byte stream assuming it is a type of comma or tab separated values
func (ce *ContentExtractor) ExtractCSV(ctx context.Context, filecontent io.Reader) ([]types.ContentLine, error) {

	out := make([]types.ContentLine, 0, 128)

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
			out = append(out, types.ContentLine{
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
