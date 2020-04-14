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
	fmt.Printf("DEBUG: useTagType: %d\n", userTagType)
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
		fmt.Println("Cursor error")
		if err == mongo.ErrNoDocuments {
			return []tag.FileTag{}, nil
		}
		return nil, err
	}
	var newTags []tag.FileTag
	fmt.Printf("oldTags length:%d\n", len(oldTags))
	for _, currTag := range oldTags {
		if idMap, ok := currTag.Data["user"]; ok {
			for id, val := range idMap {
				userID, err := types.DecodeObjectIDString(id)
				if err != nil {
					fmt.Println("Error decoding user id")
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
