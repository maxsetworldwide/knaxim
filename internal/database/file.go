package database

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"time"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
)

type ProcessingError struct {
	Status  int    `json:"status" bson:"s"`
	Message string `json:"msg" bson:"m"`
}

type FileStore struct {
	ID          filehash.StoreID `json:"id" bson:"id"`
	Content     []byte           `json:"content" bson:"-"`
	ContentType string           `json:"ctype" bson:"ctype"`
	FileSize    int64            `json:"fsize" bson:"fsize"`
	Perr        *ProcessingError `json:"err,omitempty" bson:"perr,omitempty"`
}

func NewFileStore(r io.Reader) (*FileStore, error) {
	n := new(FileStore)

	pout, pin := io.Pipe()
	ContentBuf := new(bytes.Buffer)
	g := gzip.NewWriter(ContentBuf)

	go func() {
		writeall := io.MultiWriter(pin, g)
		defer pin.Close()
		var err error
		if n.FileSize, err = io.Copy(writeall, r); err != nil {
			pin.CloseWithError(err)
		}
	}()

	var err error
	n.ID, err = filehash.NewStoreID(pout)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error F1")
	}
	if err = g.Close(); err != nil {
		return nil, srverror.New(err, 500, "Database Error F2")
	}

	n.Content = ContentBuf.Bytes()
	return n, nil
}

func (fs *FileStore) Reader() (io.Reader, error) {
	contentbuffer := bytes.NewReader(fs.Content)
	out, err := gzip.NewReader(contentbuffer)
	if err != nil {
		err = srverror.New(err, 500, "Database Error F3", "file reading error")
	}
	return out, err
}

func (fs *FileStore) Copy() *FileStore {
	c := make([]byte, len(fs.Content))
	copy(c, fs.Content)
	var perrcopy *ProcessingError
	if fs.Perr != nil {
		perrcopy = &ProcessingError{
			Status:  fs.Perr.Status,
			Message: fs.Perr.Message,
		}
	}
	return &FileStore{
		ID:          fs.ID,
		ContentType: fs.ContentType,
		FileSize:    fs.FileSize,
		Content:     c,
		Perr:        perrcopy,
	}
}

type FileI interface {
	PermissionI
	GetID() filehash.FileID
	SetID(filehash.FileID)
	GetName() string
	SetName(n string)
	Copy() FileI
}

type FileTime struct {
	Upload time.Time `json:"upload" bson:"upload"`
}

type File struct {
	Permission
	ID   filehash.FileID `json:"id" bson:"id"`
	Name string          `json:"name" bson:"name"`
	Date FileTime        `json:"date" bson:"date"`
}

type WebFile struct {
	File
	URL string `json:"url" bson:"url"`
}

func (f *File) GetID() filehash.FileID {
	return f.ID
}

func (f *File) SetID(id filehash.FileID) {
	f.ID = id
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) SetName(n string) {
	f.Name = n
}

func (f *File) Copy() FileI {
	nf := new(File)
	*nf = *f
	nf.Permission = *(f.CopyPerm(nil).(*Permission))
	return nf
}

func (f *WebFile) Copy() FileI {
	nf := new(WebFile)
	*nf = *f
	nf.Permission = *(f.CopyPerm(nil).(*Permission))
	return nf
}

func (f *File) MarshalJSON() ([]byte, error) {
	vals := f.toMap()
	vals["id"] = f.ID
	vals["name"] = f.Name
	vals["date"] = f.Date
	return json.Marshal(vals)
}

func (f *File) MarshalBSON() ([]byte, error) {
	vals := f.toMap()
	vals["id"] = f.ID
	vals["name"] = f.Name
	vals["date"] = f.Date
	return bson.Marshal(vals)
}

type fForm struct {
	ID   filehash.FileID `bson:"id" json:"id"`
	Name string          `json:"name" bson:"name"`
	URL  *string         `json:"url,omitempty" bson:"url,omitempty"`
	Date FileTime        `json:"date" bson:"date"`
}

func (f *File) UnmarshalJSON(b []byte) error {
	err := f.Permission.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	form := new(fForm)
	err = json.Unmarshal(b, form)
	if err != nil {
		return err
	}
	f.ID = form.ID
	f.Name = form.Name
	f.Date = form.Date
	return nil
}

func (f *File) UnmarshalBSON(b []byte) error {
	err := f.Permission.UnmarshalBSON(b)
	if err != nil {
		return err
	}
	form := new(fForm)
	err = bson.Unmarshal(b, form)
	if err != nil {
		return err
	}
	f.ID = form.ID
	f.Name = form.Name
	f.Date = form.Date
	return nil
}

func (wp *WebFile) MarshalJSON() ([]byte, error) {
	vals := wp.toMap()
	vals["id"] = wp.ID
	vals["name"] = wp.Name
	vals["url"] = wp.URL
	vals["date"] = wp.Date
	return json.Marshal(vals)
}

func (wp *WebFile) MarshalBSON() ([]byte, error) {
	vals := wp.toMap()
	vals["id"] = wp.ID
	vals["name"] = wp.Name
	vals["url"] = wp.URL
	vals["date"] = wp.Date
	return bson.Marshal(vals)
}

func (wp *WebFile) UnmarshalJSON(b []byte) error {
	err := wp.Permission.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	form := new(fForm)
	err = json.Unmarshal(b, form)
	if err != nil {
		return err
	}
	wp.ID = form.ID
	wp.Name = form.Name
	wp.Date = form.Date
	if form.URL != nil {
		wp.URL = *form.URL
	}
	return nil
}

func (wp *WebFile) UnmarshalBSON(b []byte) error {
	err := wp.Permission.UnmarshalBSON(b)
	if err != nil {
		return err
	}
	form := new(fForm)
	err = bson.Unmarshal(b, form)
	if err != nil {
		return err
	}
	wp.ID = form.ID
	wp.Name = form.Name
	wp.Date = form.Date
	if form.URL != nil {
		wp.URL = *form.URL
	}
	return nil
}

type FileDecoder struct {
	F *File
	W *WebFile
}

func (fd *FileDecoder) File() FileI {
	if fd.F == nil {
		return fd.W
	}
	return fd.F
}

func (fd *FileDecoder) UnmarshalJSON(b []byte) error {
	form := new(fForm)
	err := json.Unmarshal(b, form)
	if err != nil {
		return err
	}
	if form.URL == nil {
		fd.F = new(File)
		err = fd.F.Permission.UnmarshalJSON(b)
		if err != nil {
			return err
		}
		fd.F.ID = form.ID
		fd.F.Name = form.Name
		fd.F.Date = form.Date
		return nil
	} else {
		fd.W = new(WebFile)
		err = fd.W.Permission.UnmarshalJSON(b)
		if err != nil {
			return err
		}
		fd.W.ID = form.ID
		fd.W.Name = form.Name
		fd.W.URL = *form.URL
		fd.W.Date = form.Date
		return nil
	}
}

func (fd *FileDecoder) UnmarshalBSON(b []byte) error {
	form := new(fForm)
	err := bson.Unmarshal(b, form)
	if err != nil {
		return err
	}
	if form.URL == nil {
		fd.F = new(File)
		err = fd.F.Permission.UnmarshalBSON(b)
		if err != nil {
			return err
		}
		fd.F.ID = form.ID
		fd.F.Name = form.Name
		fd.F.Date = form.Date
		return nil
	} else {
		fd.W = new(WebFile)
		err = fd.W.Permission.UnmarshalBSON(b)
		if err != nil {
			return err
		}
		fd.W.ID = form.ID
		fd.W.Name = form.Name
		fd.W.URL = *form.URL
		fd.W.Date = form.Date
		return nil
	}
}
