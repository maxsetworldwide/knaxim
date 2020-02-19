package process

import (
	"io"
	"strings"

	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

type ExtConst int

const (
	_ ExtConst = iota
	OFFICE
	URL
	PDF
)

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

func MapContentType(ct string) ExtConst {
	return ExtMap[strings.TrimSpace(strings.Split(ct, ";")[0])]
}

type FileConverter gotenberg.Client

func NewFileConverter(url string) *FileConverter {
	return (*FileConverter)(&gotenberg.Client{
		Hostname: url,
	})
}

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
