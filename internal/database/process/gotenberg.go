package process

import (
	"io"

	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

type ExtConst int

const (
	OFFICE ExtConst = iota
	URL
	PDF
)

var ExtMap = map[string]ExtConst{
	".pdf":  PDF,
	".txt":  OFFICE,
	".rtf":  OFFICE,
	".fodt": OFFICE,
	".doc":  OFFICE,
	".docx": OFFICE,
	".odt":  OFFICE,
	".xls":  OFFICE,
	".xlsx": OFFICE,
	".ods":  OFFICE,
	".ppt":  OFFICE,
	".pptx": OFFICE,
	".odp":  OFFICE,
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
