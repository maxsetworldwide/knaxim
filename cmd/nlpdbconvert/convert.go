package main

import (
	"context"
	"fmt"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userTagType = tag.Type(uint32(1 << 24))

func convertUserTags(ctx context.Context, client *mongo.Client, src string) ([]tag.FileTag, error) {
	srcDB := client.Database(src)

	tagColl := srcDB.Collection("tag")
	cursor, err := tagColl.Find(ctx, bson.M{
		"type": bson.M{
			"$eq": userTagType,
		},
	})
	if err != nil {
		return nil, err
	}
	var oldTags []struct {
		File types.FileID                 `bson:"file"`
		Word string                       `bson:"word"`
		Data map[string]map[string]string `bson:"data"`
		Type tag.Type                     `bson:"type"`
	}
	if err := cursor.All(ctx, &oldTags); err != nil {
		if err == mongo.ErrNoDocuments {
			return []tag.FileTag{}, nil
		}
		return nil, err
	}
	var newTags []tag.FileTag
	for _, currTag := range oldTags {
		if idMap, ok := currTag.Data["user"]; ok {
			for id, val := range idMap {
				userID, err := types.DecodeObjectIDString(id)
				if err != nil {
					return nil, err
				}
				if val == "d" {
					newTag := tag.FileTag{
						File:  currTag.File,
						Owner: userID,
						Tag: tag.Tag{
							Word: currTag.Word,
							Type: userTagType,
						},
					}
					newTags = append(newTags, newTag)
				}
			}
		}
	}
	return newTags, nil
}

func createNLPTags(ctx context.Context, client *mongo.Client, src string) ([]tag.StoreTag, error) {
	// go through mongo client, do a find on everything in the store collection
	// may need instances of a knaxim db for both databases.
	// Need to see what is needed to process a file store. At the very least, will
	// need store ids. After that, may need to query the db object with those
	// store ids. Afterwards, will want to use the new db's knaxim db object to
	// upsert these tags.
	srcStoreColl := client.Database(src).Collection("store")
	cursor, err := srcStoreColl.Find(ctx, bson.D{})
	if err != nil {
		return []tag.StoreTag{}, err
	}
	var fileStores []types.FileStore
	err = cursor.All(ctx, &fileStores)
	if err != nil {
		return []tag.StoreTag{}, err
	}
	fmt.Println("File Store IDs:")
	for _, fs := range fileStores {
		fmt.Println(fs.ID)
	}
	return []tag.StoreTag{}, nil

}
