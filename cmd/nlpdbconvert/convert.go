// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

const userTagType = tag.USER

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
				userID, err := types.DecodeOwnerIDString(id)
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

// Assumes the store collection has already been copied over.
// This function goes through each file record and inserts name tags for each file.
// It will then clear previous processing errors that have occurred in the past and
// reprocesses each unique filestore via decode.Read() as if the file was just uploaded.
// This will add the views, lines, and nlp tags for each file to the database.
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

	err = clearPerrs(ctx, client, destDB.DBName)
	if err != nil {
		return err
	}

	var foundStoreIDs = make(map[string]bool)

	if !*quiet {
		fmt.Printf("Processing %d files:\n", len(files))
		defer fmt.Println()
	}
	for _, file := range files {
		if !*quiet {
			fmt.Print(".")
		}

		fs, err := destDB.File().Get(file.GetID())
		if err != nil {
			return err
		}

		nametags, err := tag.BuildNameTags(fs.GetName())
		if err != nil {
			return err
		}
		var fileNameTags []tag.FileTag
		for _, nt := range nametags {
			fileNameTags = append(fileNameTags, tag.FileTag{
				File:  fs.GetID(),
				Owner: fs.GetOwner().GetID(),
				Tag:   nt,
			})
		}
		err = destDB.Tag().Upsert(fileNameTags...)

		storeID := fs.GetID().StoreID
		if found := foundStoreIDs[storeID.String()]; !found {
			fs, err := destDB.Store().Get(storeID)
			if err != nil {
				return err
			}
			decode.Read(ctx, nil, file.GetName(), fs, destDB, *tikaPath, *gotenPath)
			foundStoreIDs[storeID.String()] = true
		}
	}

	return nil
}

func clearPerrs(ctx context.Context, client *mongo.Client, dbName string) error {
	_, err := client.Database(dbName).Collection("store").UpdateMany(ctx,
		bson.M{},
		bson.M{
			"$unset": bson.M{"perr": ""},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
