package main

import (
	"context"
	"fmt"

	CEMongo "git.maxset.io/web/knaxim/internal/database/mongo"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/decode"
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

// Assumes store collection has already been copied over.
// This function goes through each file record (not file store due to requiring a name),
// and calls decode.Read() on each one, as if the file was just uploaded. This will
// add the views, lines, and nlp tags for each file to the database.
func insertNLPTags(ctx context.Context, client *mongo.Client, destDB *CEMongo.Database) error {
	srcFileColl := client.Database(destDB.DBName).Collection("file")
	cursor, err := srcFileColl.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	var files []types.File
	err = cursor.All(ctx, &files)
	if err != nil {
		return err
	}

	sb := destDB.Store(ctx)
	defer sb.Close(ctx)
	if !*quiet {
		fmt.Printf("Processing %d files:\n", len(files))
		defer fmt.Println()
	}
	for _, file := range files {
		if !*quiet {
			fmt.Print(".")
		}
		storeID := file.GetID().StoreID
		fs, err := sb.Get(storeID)
		if err != nil {
			return err
		}
		decode.Read(ctx, nil, file.Name, fs, sb, *tikaPath, *gotenPath)
	}

	return nil
}
