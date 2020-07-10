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

package mongo

import (
	"context"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func TestContenbase(t *testing.T) {
	t.Parallel()
	var cb *Contentbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestContent"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to Init database", err)
		}
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		mdb, err := db.Connect(methodtesting)
		if err != nil {
			t.Fatalf("Unable to connect to database: %s", err.Error())
		}
		defer mdb.Close(methodtesting)
		cb = mdb.Content().(*Contentbase)
	}
	var fileids = []types.StoreID{
		types.StoreID{
			Hash:  7777,
			Stamp: 32621,
		},
		types.StoreID{
			Hash:  841602,
			Stamp: 28720,
		},
	}
	var fileStores = []*types.FileStore{
		&types.FileStore{
			ID:          fileids[0],
			Content:     []byte("asdfasdf"),
			ContentType: "test",
			FileSize:    420,
		},
		&types.FileStore{
			ID:          fileids[1],
			Content:     []byte("fdafdsa"),
			ContentType: "test",
			FileSize:    240,
		},
	}
	{
		sb := cb.Store()
		for _, fs := range fileStores {
			_, err := sb.Reserve(fs.ID)
			if err != nil {
				t.Fatalf("unable to Reserve file store id: %s", err)
			}
			err = sb.Insert(fs)
			if err != nil {
				t.Fatalf("unable to Insert file store: %s", err)
			}
		}
	}
	var data = []types.ContentLine{
		types.ContentLine{
			ID:       fileids[0],
			Position: 0,
			Content:  []string{"This is the first sentence."},
		},
		types.ContentLine{
			ID:       fileids[0],
			Position: 1,
			Content:  []string{"Another Sentence right here."},
		},
		types.ContentLine{
			ID:       fileids[0],
			Position: 2,
			Content:  []string{"More Sentences."},
		},
		types.ContentLine{
			ID:       fileids[1],
			Position: 0,
			Content:  []string{"This is another document."},
		},
		types.ContentLine{
			ID:       fileids[1],
			Position: 1,
			Content:  []string{"It only has 2 sentences"},
		},
	}
	t.Run("Insert", func(t *testing.T) {
		err := cb.Insert(data...)
		if err != nil {
			t.Error("Unable to Insert", err)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Len", func(t *testing.T) {
		l, err := cb.Len(fileids[0])
		if err != nil {
			t.Fatal("Err getting Length", err)
		}
		if l != 3 {
			t.Fatal("Err incorrect length", l)
		}
	})
	t.Run("Slice", func(t *testing.T) {
		result, err := cb.Slice(fileids[0], 1, 3)
		if err != nil {
			t.Fatal("Err getting slice", err)
		}
		if len(result) != 2 {
			t.Fatal("Slice returned incorrect amount", result)
		}
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("Equality check paniced", r)
			}
		}()
		for _, r := range result {
			original := data[r.Position]
			if !original.ID.Equal(r.ID) || original.Content[0] != r.Content[0] {
				t.Error("result mismatched original", r, original)
			}
		}
	})
	t.Run("Regex", func(t *testing.T) {
		result, err := cb.RegexSearchFile("only", fileids[1], 0, 2)
		if err != nil {
			t.Fatal("Err doing search", err)
		}
		if len(result) != 1 {
			t.Fatal("incorrect number of matches", result)
		}
		if !result[0].ID.Equal(data[4].ID) || result[0].Position != data[4].Position || result[0].Content[0] != data[4].Content[0] {
			t.Error("mismatched result", result[0])
		}
	})
}
