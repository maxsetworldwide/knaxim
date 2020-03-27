// basing off of database/file.go

package types

import (
	"bytes"
	"compress/gzip"
	"io"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

// ViewStore represents the PDF representation of a File
type ViewStore struct {
	ID      StoreID `json:"id" bson:"id"`
	Content []byte  `json:"content" bson:"-"`
}

// NewViewStore builds ViewStore from a StoreID and Reader of
// pdf representation
func NewViewStore(id StoreID, r io.Reader) (*ViewStore, error) {
	store := new(ViewStore)

	contentBuf := new(bytes.Buffer)
	gzWrite := gzip.NewWriter(contentBuf)

	var err error
	if _, err = io.Copy(gzWrite, r); err != nil {
		return nil, err
	}

	if err = gzWrite.Close(); err != nil {
		return nil, srverror.New(err, 500, "Database Error V2")
	}

	store.Content = contentBuf.Bytes()
	store.ID = id
	return store, nil
}

// Reader returns a reader of the content of a pdf representation of the file
func (vs *ViewStore) Reader() (io.Reader, error) {
	buf := bytes.NewReader(vs.Content)
	out, err := gzip.NewReader(buf)
	if err != nil {
		srverror.New(err, 500, "Database Error V3", "file reading error")
	}
	return out, err
}
