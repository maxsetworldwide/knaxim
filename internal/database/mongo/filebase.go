package mongo

import (
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Filebase struct {
	Database
}

func (fb *Filebase) Reserve(id filehash.FileID) (filehash.FileID, error) {
	var out *filehash.FileID
	for out == nil {
		timeout := time.Now().Add(time.Hour * 24)
		result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).UpdateOne(fb.ctx, bson.M{
			"id":      id,
			"reserve": bson.M{"$lte": time.Now()},
		}, bson.M{
			"$set": bson.M{"reserve": timeout},
		})
		if err != nil {
			return id, srverror.New(err, 500, "Database Error F1", "Unable to update id reserve")
		}
		if result.ModifiedCount > 0 {
			out = &id
		} else {
			result, err = fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).UpdateOne(
				fb.ctx,
				bson.M{"id": id},
				bson.M{"$setOnInsert": bson.M{
					"id":      id,
					"reserve": timeout,
				}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return id, srverror.New(err, 500, "Database F1.1", "Unable to upsert id")
			}
			if result.UpsertedCount > 0 {
				out = &id
			} else {
				id = id.Mutate()
			}
		}
	}
	return *out, nil
}

func (fb *Filebase) Insert(r database.FileI) error {
	result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).UpdateOne(
		fb.ctx,
		bson.M{
			"id":      r.GetID(),
			"reserve": bson.M{"$gt": time.Now()},
		},
		bson.M{
			"$unset": bson.M{"reserve": ""},
			"$set":   r,
		},
	)
	if err != nil {
		return srverror.New(err, 500, "Database Error F2", "Unable to insert")
	}
	if result.ModifiedCount == 0 {
		return database.ErrIDNotReserved
	}
	return nil
}

func (fb *Filebase) Get(fid filehash.FileID) (database.FileI, error) {
	result := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).FindOne(fb.ctx, bson.M{
		"id": fid,
	})
	fd := new(database.FileDecoder)
	if err := result.Decode(fd); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F3", "Unable to get file")
	}
	f := fd.File()
	err := f.Populate(fb.Owner(nil))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fb *Filebase) GetAll(fids ...filehash.FileID) ([]database.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"id": bson.M{"$in": fids},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F3.1", "Unable to get files")
	}
	return fb.decodefiles(cursor)
}

func (fb *Filebase) Update(r database.FileI) error {
	result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).ReplaceOne(fb.ctx, bson.M{
		"id": r.GetID(),
	}, r)
	if err != nil {
		return srverror.New(err, 500, "Database Error F4", "error updating file")
	}
	if result.ModifiedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (fb *Filebase) Remove(r filehash.FileID) error {
	result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).DeleteOne(fb.ctx, bson.M{
		"id": r,
	})
	if err != nil {
		return srverror.New(err, 500, "Database Error F5", "unable to remove file", r.String())
	}
	if result.DeletedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (fb *Filebase) decodefiles(cursor *mongo.Cursor) ([]database.FileI, error) {
	var reference []*database.FileDecoder
	if err := cursor.All(fb.ctx, &reference); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F6", "unable to decode file list")
	}
	files := make([]database.FileI, 0, len(reference))
	for _, ref := range reference {
		files = append(files, ref.File())
	}
	for _, file := range files {
		err := file.Populate(fb.Owner(nil))
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func (fb *Filebase) GetOwned(uid database.OwnerID) ([]database.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"own": uid,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F7", "unable to send request")
	}
	return fb.decodefiles(cursor)
}

func (fb *Filebase) GetPermKey(uid database.OwnerID, pkey string) ([]database.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"perm." + pkey: uid,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F8", "unable to send request")
	}
	return fb.decodefiles(cursor)
}

func (fb *Filebase) MatchStore(oid database.OwnerID, sid []filehash.StoreID, pkeys ...string) ([]database.FileI, error) {
	query := bson.M{
		"id.storeid": bson.M{"$in": sid},
	}
	or := make(bson.A, 0, 1+len(pkeys))
	or = append(or, bson.M{"own": oid})
	for _, p := range pkeys {
		or = append(or, bson.M{"perm." + p: oid})
	}
	query["$or"] = or
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(
		fb.ctx,
		query,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error F9", "unable to send request")
	}
	return fb.decodefiles(cursor)
}
