package types

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

// StoreID is used to uniquely identify a FileStore.
// Contains a hash of the file's content used to quickly identify copies of the
// same file.
type StoreID struct {
	Hash  uint32 `bson:"hash"`
	Stamp uint16 `bson:"stamp"`
}

// ToNum converts StoreID to a int64 representation
func (sid StoreID) ToNum() int64 {
	return int64(sid.Stamp)<<32 | int64(sid.Hash)
}

// NewStoreID builds a StoreID from a reader of the file's contents
func NewStoreID(r io.Reader) (StoreID, error) {
	return NewStoreIDComplete(r, adler32.New(), uint16(time.Now().UnixNano()))
}

// NewStoreIDComplete build a StoreID from a reader of the file's contents,
// a hash scheme and stamp value
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

// String returns a base64 encoding as a string
func (sid StoreID) String() string {
	build := new(strings.Builder)
	encoder := base64.NewEncoder(base64.URLEncoding, build)

	binary.Write(encoder, binary.LittleEndian, sid.Hash)
	binary.Write(encoder, binary.LittleEndian, sid.Stamp)
	encoder.Close()

	return build.String()
}

// Mutate returns a new StoreID with the same hash
func (sid StoreID) Mutate() StoreID {
	sid.Stamp++
	return sid
}

// DecodeStoreID returns a StoreID from base64 encoded string, inverse of String()
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

// Equal returns true if the other StoreID has the same value
func (sid StoreID) Equal(oid StoreID) bool {
	return sid.Hash == oid.Hash && sid.Stamp == oid.Stamp
}
