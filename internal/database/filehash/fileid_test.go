package filehash

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	fid := FileID{
		StoreID: StoreID{
			Hash:  15,
			Stamp: 16,
		},
		Stamp: []byte("test"),
	}
	jsonbytes, err := json.Marshal(fid)
	if err != nil {
		t.Fatal("Unable to encode FileID: ", err)
	}
	var unmarshaled FileID
	err = json.Unmarshal(jsonbytes, &unmarshaled)
	if err != nil {
		t.Fatal("Unable to decode FileID: ", err)
	}
	if !fid.Equal(unmarshaled) {
		t.Fatalf("fid mismatched: %+#v", unmarshaled)
	}
}
