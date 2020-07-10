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

package process_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/memory"
	. "git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/types"
)

func TestInjustFile(t *testing.T) {
	var db = &memory.Database{}

	var content = `This is an example text.
  How do you know this is correct. Well you got this text back.
  Good Bye.`

	var testOwner = &types.User{
		ID: types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'a', 'b', 'c'},
			Stamp:       []byte("test"),
		},
		Name: "testuser",
	}
	initctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err := db.Init(initctx, true)
	if err != nil {
		t.Fatal("unable to init database", err)
	}
	injestctx, cancel2 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel2()
	file := &types.File{
		Name: "test.txt",
	}
	_, err = db.Owner().Reserve(testOwner.GetID(), testOwner.GetName())
	if err != nil {
		t.Fatal("unable to reserve testOwner:", err)
	}
	err = db.Owner().Insert(testOwner)
	if err != nil {
		t.Fatalf("Failed to insert test Owner")
	}
	file.Own = testOwner
	_, err = InjestFile(injestctx, file, "content/txt", strings.NewReader(content), db)
	if err != nil {
		t.Fatal("injest failed", err)
	}

	// generateContentTags(injestctx, fs, db)
}
