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

	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

func initContentIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["lines"]).Indexes()
	_, err := I.CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "id", Value: 1}, bson.E{Key: "position", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// Contentbase database connection with content lines operations
type Contentbase struct {
	Database
}

// Insert adds lines to the database
func (cb *Contentbase) Insert(lines ...types.ContentLine) error {
	var docs []interface{}
	for _, line := range lines {
		docs = append(docs, line)
	}
	_, err := cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).InsertMany(cb.ctx, docs)
	if err != nil {
		return srverror.New(err, 500, "Error C1", "Unable to Insert")
	}
	return nil
}

// Len returns number of lines associated with StoreID
func (cb *Contentbase) Len(id types.StoreID) (count int64, err error) {
	count, err = cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).CountDocuments(cb.ctx, bson.M{
		"id": id,
	})
	if err != nil {
		err = srverror.New(err, 500, "Error C2")
	}
	return
}

// Slice returns slices associated with StoreID within bounds
func (cb *Contentbase) Slice(id types.StoreID, start int, end int) ([]types.ContentLine, error) {
	fs, err := cb.Store().Get(id)
	if err != nil {
		return nil, err
	}
	var perr error
	if fs.Perr != nil {
		perr = fs.Perr
	}
	cursor, err := cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).Find(cb.ctx,
		bson.M{
			"id": id,
			"position": bson.M{
				"$gte": start,
				"$lt":  end,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no slices in range")
		}
		return nil, srverror.New(err, 500, "Error C3", "Failed to find lines")
	}
	var out []types.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no slices when decoded")
		}
		return nil, srverror.New(err, 500, "Error C3.1", "failed to decode lines")
	}
	return out, perr
}

// RegexSearchFile returns lines associated with StoreID, within bounds, and matches regular expression
func (cb *Contentbase) RegexSearchFile(regex string, id types.StoreID, start int, end int) ([]types.ContentLine, error) {
	fs, err := cb.Store().Get(id)
	if err != nil {
		return nil, err
	}
	var perr error
	if fs.Perr != nil {
		perr = fs.Perr
	}
	cursor, err := cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).Find(cb.ctx, bson.M{
		"id": id,
		"position": bson.M{
			"$gte": start,
			"$lt":  end,
		},
		"content": bson.M{"$regex": regex, "$options": "i"},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no matches in range")
		}
		return nil, srverror.New(err, 500, "Error C4", "Failed to find lines")
	}
	var out []types.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no matches decoded")
		}
		return nil, srverror.New(err, 500, "Error C4.1", "failed to decode lines")
	}
	return out, perr
}
