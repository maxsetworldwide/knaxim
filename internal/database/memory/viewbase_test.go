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

package memory

import (
	"bytes"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
)

var contentString = "View version of a file for the memory database"

var inputVS = &types.ViewStore{
	ID: types.StoreID{
		Hash:  98765,
		Stamp: 4321,
	},
	Content: []byte(contentString),
}

func TestViewbase(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	vb := DB.View()
	defer DB.Close(nil)
	t.Parallel()

	t.Log("View Insert")
	err := vb.Insert(inputVS)
	if err != nil {
		t.Fatalf("error inserting view: %s\n", err)
	}
	t.Log("View Get")
	result, err := vb.Get(inputVS.ID)
	if err != nil {
		t.Fatalf("error getting view: %s\n", err)
	}
	if !inputVS.ID.Equal(result.ID) || !bytes.Equal(inputVS.Content, result.Content) {
		t.Fatalf("Did not get correct view store:\ngot: %+#v\nexpected: %+#v\n", result, inputVS)
	}
}
