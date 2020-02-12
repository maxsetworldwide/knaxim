package process

import (
	"bytes"

	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

type FileConverter gotenberg.Client

func NewFileConverter(url string) *FileConverter {
	return (*FileConverter)(&gotenberg.Client{
		Hostname: url,
	})
}

func (fc *FileConverter) ConvertOffice(inputName string, in *bytes.Buffer) (out *bytes.Buffer, err error) {
	out = new(bytes.Buffer)
	index, err := gotenberg.NewDocumentFromBytes(inputName, in.Bytes())
	if err != nil {
		return
	}
	req := gotenberg.NewOfficeRequest(index)
	res, err := (*gotenberg.Client)(fc).Post(req)
	if err != nil {
		return
	}
	out.ReadFrom(res.Body)
	return
}
