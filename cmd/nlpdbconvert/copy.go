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
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collsToCopy = []string{
	"acronym",
	"chunk",
	"file",
	"group",
	"reset",
	"store",
	"user",
}

func copyColls(ctx context.Context, client *mongo.Client, src, dest string) error {
	srcDB := client.Database(src)
	destDB := client.Database(dest)

	for _, coll := range collsToCopy {
		srcColl := srcDB.Collection(coll)
		destColl := destDB.Collection(coll)
		cursor, err := srcColl.Find(ctx, bson.D{})
		if err != nil {
			return err
		}
		var result []interface{}
		if err := cursor.All(ctx, &result); err == nil {
			if len(result) > 0 {
				_, err = destColl.InsertMany(ctx, result)
				if err != nil {
					return err
				}
			} else if !*quiet {
				fmt.Printf("No data in %s collection\n", coll)
			}
		} else if err != mongo.ErrNoDocuments {
			return err
		}
	}
	return nil
}
