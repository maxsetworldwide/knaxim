package types

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"time"

	"git.maxset.io/web/knaxim/pkg/srverror"

	dberrs "git.maxset.io/web/knaxim/internal/database/types/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// FileStore represents a file's content
type FileStore struct {
	ID          StoreID            `json:"id" bson:"id"`
	Content     []byte             `json:"content" bson:"-"`
	ContentType string             `json:"ctype" bson:"ctype"`
	FileSize    int64              `json:"fsize" bson:"fsize"`
	Perr        *dberrs.Processing `json:"err,omitempty" bson:"perr,omitempty"`
}

// NewFileStore builds a FileStore from a reader of the file content
func NewFileStore(r io.Reader) (*FileStore, error) {
	n := new(FileStore)
	n.Perr = dberrs.FileLoadInProgress

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
	n.ID, err = NewStoreID(pout)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error F1")
	}
	if err = g.Close(); err != nil {
		return nil, srverror.New(err, 500, "Database Error F2")
	}

	n.Content = ContentBuf.Bytes()
	return n, nil
}

// Reader returns a reader of the file contents
func (fs *FileStore) Reader() (io.Reader, error) {
	contentbuffer := bytes.NewReader(fs.Content)
	out, err := gzip.NewReader(contentbuffer)
	if err != nil {
		err = srverror.New(err, 500, "Database Error F3", "file reading error")
	}
	return out, err
}

// Copy returns a new instance of a FileStore
func (fs *FileStore) Copy() *FileStore {
	c := make([]byte, len(fs.Content))
	copy(c, fs.Content)
	var perrcopy *dberrs.Processing
	if fs.Perr != nil {
		perrcopy = &dberrs.Processing{
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

// FileI is the interface type to represent a file
type FileI interface {
	PermissionI
	GetID() FileID
	SetID(FileID)
	GetName() string
	SetName(n string)
	GetDate() FileTime
	Copy() FileI
}

// FileTime is a store of the relevant times of a file
type FileTime struct {
	Upload time.Time `json:"upload" bson:"upload"`
}

// File is the meta data and permission values of a file
type File struct {
	Permission
	ID   FileID   `json:"id" bson:"id"`
	Name string   `json:"name" bson:"name"`
	Date FileTime `json:"date" bson:"date"`
}

// WebFile is an extention of a File with URL value
type WebFile struct {
	File
	URL string `json:"url" bson:"url"`
}

// GetID implements FileI
func (f *File) GetID() FileID {
	return f.ID
}

// SetID implements FileI
func (f *File) SetID(id FileID) {
	f.ID = id
}

// GetName implements FileI
func (f *File) GetName() string {
	return f.Name
}

// SetName implements FileI
func (f *File) SetName(n string) {
	f.Name = n
}

// GetDate implements FileI
func (f *File) GetDate() FileTime {
	return f.Date
}

// Copy build a new instance of the File
func (f *File) Copy() FileI {
	nf := new(File)
	*nf = *f
	nf.Permission = *(f.CopyPerm(nil).(*Permission))
	return nf
}

// Copy build a new instance of the WebFile
func (wp *WebFile) Copy() FileI {
	nf := new(WebFile)
	*nf = *wp
	nf.Permission = *(wp.CopyPerm(nil).(*Permission))
	return nf
}

// MarshalJSON converts file into json representation
func (f *File) MarshalJSON() ([]byte, error) {
	vals := f.toMap()
	vals["id"] = f.ID
	vals["name"] = f.Name
	vals["date"] = f.Date
	return json.Marshal(vals)
}

// MarshalBSON converts file into bson representation
func (f *File) MarshalBSON() ([]byte, error) {
	vals := f.toMap()
	vals["id"] = f.ID
	vals["name"] = f.Name
	vals["date"] = f.Date
	return bson.Marshal(vals)
}

type fForm struct {
	ID   FileID   `bson:"id" json:"id"`
	Name string   `json:"name" bson:"name"`
	URL  *string  `json:"url,omitempty" bson:"url,omitempty"`
	Date FileTime `json:"date" bson:"date"`
}

// UnmarshalJSON converts json to File
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

// UnmarshalBSON convert bson to file
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

// MarshalJSON converts webfile to json
func (wp *WebFile) MarshalJSON() ([]byte, error) {
	vals := wp.toMap()
	vals["id"] = wp.ID
	vals["name"] = wp.Name
	vals["url"] = wp.URL
	vals["date"] = wp.Date
	return json.Marshal(vals)
}

// MarshalBSON converts webfile to bson
func (wp *WebFile) MarshalBSON() ([]byte, error) {
	vals := wp.toMap()
	vals["id"] = wp.ID
	vals["name"] = wp.Name
	vals["url"] = wp.URL
	vals["date"] = wp.Date
	return bson.Marshal(vals)
}

// UnmarshalJSON converts json to File
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

// UnmarshalBSON converts bson to file
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

// FileDecoder Unmarshals as either a File or WebFile
type FileDecoder struct {
	F *File
	W *WebFile
}

// File returns FileI that was decoded
func (fd *FileDecoder) File() FileI {
	if fd.F == nil {
		return fd.W
	}
	return fd.F
}

// UnmarshalJSON decodes json to either File or WebFile
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
	}
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

// UnmarshalBSON decodes bson into WebFile or File
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
	}
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
