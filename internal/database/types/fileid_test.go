/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package types

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
