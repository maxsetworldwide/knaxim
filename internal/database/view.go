// basing off of database/file.go

package database

import (
	"bytes"
	"compress/gzip"
	"io"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

type ViewStore struct {
	ID      filehash.StoreID `json:"id" bson:"id"`
	Content []byte           `json:"content" bson:"-"`
}

func NewViewStore(r io.Reader) (*ViewStore, error) {
	store := new(ViewStore)

	pipeRead, pipeWrite := io.Pipe()
	contentBuf := new(bytes.Buffer)
	gzWrite := gzip.NewWriter(contentBuf)

	go func() {
		writeall := io.MultiWriter(pipeWrite, gzWrite)
		defer pipeWrite.Close()
		var err error
		if _, err = io.Copy(writeall, r); err != nil {
			pipeWrite.CloseWithError(err)
		}
	}()

	var err error
	store.ID, err = filehash.NewStoreID(pipeRead)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error V1")
	}
	if err = gzWrite.Close(); err != nil {
		return nil, srverror.New(err, 500, "Database Error V2")
	}

	store.Content = contentBuf.Bytes()
	return store, nil
}

func (vs *ViewStore) Reader() (io.Reader, error) {
	buf := bytes.NewReader(vs.Content)
	out, err := gzip.NewReader(buf)
	if err != nil {
		srverror.New(err, 500, "Database Error V3", "file reading error")
	}
	return out, err
}
