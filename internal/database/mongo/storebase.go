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
	"fmt"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initStoreIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["store"]).Indexes()
	_, err := I.CreateMany(
		ctx,
		[]mongo.IndexModel{
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetUnique(true),
			},
			mongo.IndexModel{
				Keys: bson.M{"id.hash": 1},
			},
		})
	return err
}

// Storebase is a connection to the database with file store operations
type Storebase struct {
	Database
}

// Reserve a store id, will mutate if store id not available, returns reserved store id
func (db *Storebase) Reserve(id types.StoreID) (types.StoreID, error) {
	var out *types.StoreID
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
			return id, srverror.New(err, 500, "Error S1", "unable to update reserve")
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
				return id, srverror.New(err, 500, "Error S1.1", "unable to insert new id")
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

// Insert file store, must reserve store id first, see Reserve
func (db *Storebase) Insert(fs *types.FileStore) error {
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
			return srverror.New(e, 500, "Error S2", "unable to insert store")
		}
		if result.ModifiedCount == 0 {
			return errors.ErrIDNotReserved
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
			return srverror.New(e, 500, "Error S3", "failed to insert data chunks")
		}
	}
	return nil
}

// Get file store
func (db *Storebase) Get(id types.StoreID) (out *types.FileStore, err error) {
	result := db.client.Database(db.DBName).Collection(db.CollNames["store"]).FindOne(
		db.ctx,
		bson.M{"id": id},
	)
	var store = new(types.FileStore)
	if err := result.Decode(store); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound.Extend("FileStore", id.String())
		}
		return nil, srverror.New(err, 500, "Error S4", "failed to find file store")
	}
	var chunks []*contentchunk
	cursor, err := db.client.Database(db.DBName).Collection(db.CollNames["chunk"]).Find(
		db.ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return nil, srverror.New(err, 500, "Error S5", "failed to get file data chunks")
	}
	if err = cursor.All(db.ctx, &chunks); err != nil {
		return nil, srverror.New(err, 500, "Error S6", "failed to decode file data chunks")
	}
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = srverror.New(v, 500, "Error S7", "unable to build chunks")
				out = nil
			default:
				err = srverror.New(fmt.Errorf("GetStore: %+#v", v), 500, "Error S8")
				out = nil
			}
		}
	}()
	store.Content = appendchunks(chunksort(chunks))
	return store, nil
}

// MatchHash returns all file stores with matching hashes
func (db *Storebase) MatchHash(h uint32) (out []*types.FileStore, err error) {
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
			cherr <- srverror.New(err, 500, "Error S9", "Match Hash Find error")
			return
		}
		if err = cursor.All(ctx, &out); err != nil {
			cherr <- srverror.New(err, 500, "Error S10", "MatchHash unable to decode FileStore's")
			return
		}
	}()
	go func() {
		cursor, err := db.client.Database(db.DBName).Collection(db.CollNames["chunk"]).Find(ctx, bson.M{
			"id.hash": h,
		})
		var chunks []*contentchunk
		if err != nil {
			cherr <- srverror.New(err, 500, "Error S11", "Unable to get filestore chunks")
			return
		}
		if err = cursor.All(ctx, &chunks); err != nil {
			cherr <- srverror.New(err, 500, "Error S12", "Unable to decode chunks")
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

// UpdateMeta of file store
func (db *Storebase) UpdateMeta(fs *types.FileStore) error {
	result, err := db.client.Database(db.DBName).Collection(db.CollNames["store"]).ReplaceOne(db.ctx, bson.M{
		"id": fs.ID,
	}, fs)
	if err != nil {
		return srverror.New(err, 500, "Error S13", "error updating file store metadata")
	}
	if result.ModifiedCount == 0 {
		return errors.ErrNotFound.Extend("no FileStore to update", fs.ID.String())
	}
	return nil
}
