package mongo

import (
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

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
		return srverror.New(err, 500, "Database Error C1", "Unable to Insert")
	}
	return nil
}

// Len returns number of lines associated with StoreID
func (cb *Contentbase) Len(id types.StoreID) (count int64, err error) {
	count, err = cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).CountDocuments(cb.ctx, bson.M{
		"id": id,
	})
	if err != nil {
		err = srverror.New(err, 500, "Database Error C2")
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
		return nil, srverror.New(err, 500, "Database Error C3", "Failed to find lines")
	}
	var out []types.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no slices when decoded")
		}
		return nil, srverror.New(err, 500, "Database Error C3.1", "failed to decode lines")
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
		return nil, srverror.New(err, 500, "Database Error C4", "Failed to find lines")
	}
	var out []types.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			if perr != nil {
				return nil, perr
			}
			return nil, errors.ErrNoResults.Extend("no matches decoded")
		}
		return nil, srverror.New(err, 500, "Database Error C4.1", "failed to decode lines")
	}
	return out, perr
}
