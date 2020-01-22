package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storebase struct {
	Database
}

func (db *Storebase) Reserve(id filehash.StoreID) (filehash.StoreID, error) {
	var out *filehash.StoreID
	for out == nil {
		timeout := time.Now().Add(time.Hour * 24)
		//check if id reservation has timed out
		result, err := db.client.Database(db.DBName).Collection(db.CollNames["store"]).UpdateOne(
			db.ctx,
			bson.M{
				"id":      id,
				"reserve": bson.M{"$lte": time.Now()},
			},
			bson.M{
				"$set": bson.M{"reserve": timeout},
			},
		)
		if err != nil {
			return id, srverror.New(err, 500, "Database Error S1", "unable to update reserve")
		}
		if result.ModifiedCount > 0 {
			out = &id
		} else {
			// add id if no in database
			result, err = db.client.Database(db.DBName).Collection(db.CollNames["store"]).UpdateOne(
				db.ctx,
				bson.M{"id": id},
				bson.M{"$setOnInsert": bson.M{
					"id":      id,
					"reserve": timeout,
				}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return id, srverror.New(err, 500, "Database Error S1.1", "unable to insert new id")
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

func (db *Storebase) Insert(fs *database.FileStore) error {
	{
		result, e := db.client.Database(db.DBName).Collection(db.CollNames["store"]).UpdateOne(
			db.ctx,
			bson.M{"id": fs.ID, "reserve": bson.M{"$gt": time.Now()}},
			bson.M{
				"$unset": bson.M{"reserve": ""},
				"$set":   fs,
			},
		)
		if e != nil {
			return srverror.New(e, 500, "Database Error S2", "unable to insert store")
		}
		if result.ModifiedCount == 0 {
			return database.ErrIDNotReserved
		}
	}
	{
		chunks := chunkify(fs.ID, fs.Content)
		_, e := db.client.Database(db.DBName).Collection(db.CollNames["chunk"]).InsertMany(
			db.ctx,
			chunks,
			options.InsertMany().SetOrdered(false),
		)
		if e != nil {
			return srverror.New(e, 500, "Database Error S3", "failed to insert data chunks")
		}
	}
	return nil
}

func (db *Storebase) Get(id filehash.StoreID) (out *database.FileStore, err error) {
	result := db.client.Database(db.DBName).Collection(db.CollNames["store"]).FindOne(
		db.ctx,
		bson.M{"id": id},
	)
	var store = new(database.FileStore)
	if err := result.Decode(store); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error S4", "failed to find file store")
	}
	var chunks []*contentchunk
	cursor, err := db.client.Database(db.DBName).Collection(db.CollNames["chunk"]).Find(
		db.ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error S5", "failed to get file data chunks")
	}
	if err = cursor.All(db.ctx, &chunks); err != nil {
		return nil, srverror.New(err, 500, "Database Error S6", "failed to decode file data chunks")
	}
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = srverror.New(v, 500, "Database Error S7", "unable to build chunks")
				out = nil
			default:
				err = srverror.New(fmt.Errorf("GetStore: %+#v", v), 500, "Database Error S8")
				out = nil
			}
		}
	}()
	store.Content = appendchunks(chunksort(chunks))
	return store, nil
}

func (db *Storebase) MatchHash(h uint32) (out []*database.FileStore, err error) {
	ctx, cancel := context.WithCancel(db.ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	cherr := make(chan error, 2)
	go func() {
		defer wg.Done()
		cursor, err := db.client.Database(db.DBName).Collection(db.CollNames["store"]).Find(ctx, bson.M{
			"id.hash": h,
		})
		if err != nil {
			cherr <- srverror.New(err, 500, "Database Error S8", "Match Hash Find error")
			return
		}
		if err = cursor.All(ctx, &out); err != nil {
			cherr <- srverror.New(err, 500, "Database Error S9", "MatchHash unable to decode FileStore's")
			return
		}
	}()
	go func() {
		cursor, err := db.client.Database(db.DBName).Collection(db.CollNames["chunk"]).Find(ctx, bson.M{
			"id.hash": h,
		})
		var chunks []*contentchunk
		if err != nil {
			cherr <- srverror.New(err, 500, "Database Error S10", "Unable to get filestore chunks")
			return
		}
		if err = cursor.All(ctx, &chunks); err != nil {
			cherr <- srverror.New(err, 500, "Database Error S11", "Unable to decode chunks")
			return
		}
		data := make(map[int64][]byte)
		for _, fchunks := range filterchunks(chunks) {
			data[fchunks[0].ID.ToNum()] = appendchunks(chunksort(fchunks))
		}
		wg.Wait()
		for i, fs := range out {
			out[i].Content = data[fs.ID.ToNum()]
		}
		cherr <- nil
	}()
	err = <-cherr
	return
}
