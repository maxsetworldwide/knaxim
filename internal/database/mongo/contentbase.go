package mongo

import (
	"git.maxset.io/server/knaxim/database"
	"git.maxset.io/server/knaxim/database/filehash"
	"git.maxset.io/server/knaxim/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Contentbase struct {
	Database
}

func (cb *Contentbase) Insert(lines ...database.ContentLine) error {
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

func (cb *Contentbase) Len(id filehash.StoreID) (count int64, err error) {
	count, err = cb.client.Database(cb.DBName).Collection(cb.CollNames["lines"]).CountDocuments(cb.ctx, bson.M{
		"id": id,
	})
	if err != nil {
		err = srverror.New(err, 500, "Database Error C2")
	}
	return
}

func (cb *Contentbase) Slice(id filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
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
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error C3", "Failed to find lines")
	}
	var out []database.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error C3.1", "failed to decode lines")
	}
	return out, nil
}

func (cb *Contentbase) RegexSearchFile(regex string, id filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
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
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error C4", "Failed to find lines")
	}
	var out []database.ContentLine
	if err = cursor.All(cb.ctx, &out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error C4.1", "failed to decode lines")
	}
	return out, nil
}
