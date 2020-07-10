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
package main

import (
	"context"
	"os"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/mongo"
)

func TestCli(t *testing.T) {
	db := new(mongo.Acronymbase)
	timeoutctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	db.URI = *databaseURI
	db.DBName = *databaseName
	db.CollNames = make(map[string]string)
	db.CollNames["acronym"] = *collectionName
	err := db.Init(timeoutctx, true)
	if err != nil {
		t.Fatal("unable to init database ", err)
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal("unable to create Pipe ", err)
	}
	os.Stdin = pr
	pw.Write([]byte("ab,acronymbase"))
	pw.Close()
	main()
}
