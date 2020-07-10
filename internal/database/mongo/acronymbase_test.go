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
)

func TestAcronym(t *testing.T) {
	t.Parallel()
	var ab *Acronymbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestAcronym"
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
		ab = mdb.Acronym().(*Acronymbase)
	}
	t.Run("Put", func(t *testing.T) {
		err := ab.Put("ab", "Acronymbase")
		if err != nil {
			t.Fatal("Unable to add acronym, ", err)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get", func(t *testing.T) {
		result, err := ab.Get("ab")
		if err != nil {
			t.Fatal("Unable to get acronym, ", err)
		}
		if len(result) != 1 || result[0] != "Acronymbase" {
			t.Fatal("incorrect result: ", result)
		}
	})
}
