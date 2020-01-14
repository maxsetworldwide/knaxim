package filehash

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"strings"
	"time"
)

type StoreID struct {
	Hash  uint32 `bson:"hash"`
	Stamp uint16 `bson:"stamp"`
}

func (sid StoreID) ToNum() int64 {
	return int64(sid.Stamp)<<32 | int64(sid.Hash)
}

func NewStoreID(r io.Reader) (StoreID, error) {
	return NewStoreIDComplete(r, adler32.New(), uint16(time.Now().UnixNano()))
}

func NewStoreIDComplete(r io.Reader, h hash.Hash32, s uint16) (StoreID, error) {
	var st StoreID
	st.Stamp = s

	h.Reset()
	if _, err := io.Copy(h, r); err != nil {
		return StoreID{}, err
	}

	st.Hash = h.Sum32()

	return st, nil
}

func (s StoreID) String() string {
	build := new(strings.Builder)
	encoder := base64.NewEncoder(base64.URLEncoding, build)

	binary.Write(encoder, binary.LittleEndian, s.Hash)
	binary.Write(encoder, binary.LittleEndian, s.Stamp)
	encoder.Close()

	return build.String()
}

func (s StoreID) Mutate() StoreID {
	s.Stamp++
	return s
}

func DecodeStoreID(h string) (sid StoreID, err error) {
	defer func() {
		if r := recover(); r != nil {
			sid = StoreID{}
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("Unable to decode storeid: %v", r)
			}
		}
	}()
	buf, err := base64.URLEncoding.DecodeString(h)
	if err != nil {
		return StoreID{}, err
	}
	return decodeStoreID(buf), nil
}

func decodeStoreID(b []byte) StoreID {
	var n StoreID
	n.Hash = binary.LittleEndian.Uint32(b[:4])
	n.Stamp = binary.LittleEndian.Uint16(b[4:])
	return n
}

func (sid StoreID) Equal(oid StoreID) bool {
	return sid.Hash == oid.Hash && sid.Stamp == oid.Stamp
}
