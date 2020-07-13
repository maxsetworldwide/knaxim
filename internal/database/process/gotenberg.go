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
	"io"
	"regexp"
	"strings"

	gotenberg "github.com/thecodingmachine/gotenberg-go-client/v7"
)

// ExtConst is an enum type indicating how a particular file
// type is to be processed
type ExtConst int

// Values of ExtConst
const (
	_ ExtConst = iota
	// Convertable to PDF
	OFFICE
	// URL to a website to convert to PDF
	URL
	// Already a PDF
	PDF
)

// ExtMap mapps Content-Type to ExtConst
var ExtMap = map[string]ExtConst{
	"application/pdf": PDF,
	"text/plain":      OFFICE,
	"application/rtf": OFFICE,
	"application/vnd.oasis.opendocument.text-template": OFFICE,
	"application/msword": OFFICE,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   OFFICE,
	"application/vnd.oasis.opendocument.text":                                   OFFICE,
	"application/vnd.ms-excel":                                                  OFFICE,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         OFFICE,
	"application/vnd.oasis.opendocument.spreadsheet":                            OFFICE,
	"application/vnd.ms-powerpoint":                                             OFFICE,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": OFFICE,
	"application/vnd.oasis.opendocument.presentation":                           OFFICE,
	"text/html": URL,
}

// MapContentType converts the content type header value to associated ExtConst
func MapContentType(ct string) ExtConst {
	return ExtMap[strings.TrimSpace(strings.Split(ct, ";")[0])]
}

var extRegex = regexp.MustCompile(`\.[^.]*$`)

// IdentifyFileAction generates the ExtConst based on the file header
func IdentifyFileAction(name string, ctype string) ExtConst {
	//check file name
	ext := strings.ToLower(extRegex.FindString(strings.TrimSpace(name)))
	if len(ext) > 0 {
		ext = ext[1:]
		switch ext {
		case "pdf":
			return PDF
		case "txt":
			fallthrough
		case "odt":
			fallthrough
		case "ods":
			fallthrough
		case "odp":
			fallthrough
		case "ppt":
			fallthrough
		case "pptx":
			fallthrough
		case "xls":
			fallthrough
		case "xlsx":
			fallthrough
		case "doc":
			fallthrough
		case "docx":
			fallthrough
		case "rtf":
			return OFFICE
		case "html":
			fallthrough
		case "htm":
			return URL
		}
	}
	return MapContentType(ctype)
}

// FileConverter is n connector to gotenberg to convert a wide
// variety of files into pdfs
type FileConverter gotenberg.Client

// NewFileConverter builds new connector to gotenberg server
// at url
func NewFileConverter(url string) *FileConverter {
	return (*FileConverter)(&gotenberg.Client{
		Hostname: url,
	})
}

// ConvertOffice converts many common file types into PDFs
func (fc *FileConverter) ConvertOffice(inputName string, in []byte) (out []byte, err error) {
	index, err := gotenberg.NewDocumentFromBytes(inputName, in)
	if err != nil {
		return
	}
	req := gotenberg.NewOfficeRequest(index)
	res, err := (*gotenberg.Client)(fc).Post(req)
	if err != nil {
		return
	}
	out = make([]byte, res.ContentLength)
	io.ReadFull(res.Body, out)

	return
}

// ConvertURL produces PDF version of webpage at address
func (fc *FileConverter) ConvertURL(url string) (out []byte, err error) {
	req := gotenberg.NewURLRequest(url)
	req.Margins(gotenberg.NoMargins)
	res, err := (*gotenberg.Client)(fc).Post(req)
	if err != nil {
		return
	}
	out = make([]byte, res.ContentLength)
	io.ReadFull(res.Body, out)

	return
}
