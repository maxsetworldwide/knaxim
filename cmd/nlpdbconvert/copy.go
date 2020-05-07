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
