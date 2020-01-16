package filehash

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/brand"
)

type FileID struct {
	StoreID
	Stamp []byte `bson:"stamp"`
}

func NewFileID(st StoreID) FileID {
	var n FileID
	n.StoreID = st
	n.Stamp = []byte{brand.Next()}
	return n
}

func (f FileID) String() string {
	build := new(strings.Builder)
	encoder := base64.NewEncoder(base64.RawURLEncoding, build)

	binary.Write(encoder, binary.LittleEndian, f.StoreID.Hash)
	binary.Write(encoder, binary.LittleEndian, f.StoreID.Stamp)
	encoder.Write(f.Stamp)
	encoder.Close()

	return build.String()
}

func (f FileID) Mutate() FileID {
	f.Stamp = append(f.Stamp, brand.Next())
	return f
}

func DecodeFileID(h string) (fid FileID, err error) {
	defer func() {
		if r := recover(); r != nil {
			fid = FileID{}
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("Failed to decode FileID: %v", r)
			}
		}
	}()
	buf, err := base64.RawURLEncoding.DecodeString(h)
	if err != nil {
		return FileID{}, err
	}
	var n FileID
	n.StoreID = decodeStoreID(buf[:6])
	n.Stamp = buf[6:]
	return n, nil
}

func (f FileID) Equal(oth FileID) bool {
	return f.StoreID.Equal(oth.StoreID) && bytes.Equal(f.Stamp, oth.Stamp)
}

func (f FileID) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *FileID) UnmarshalJSON(b []byte) error {
	var fstr string
	err := json.Unmarshal(b, &fstr)
	if err != nil {
		return err
	}
	*f, err = DecodeFileID(fstr)
	return err
}
