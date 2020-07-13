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
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initFileIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["file"]).Indexes()
	_, err := I.CreateMany(
		ctx,
		[]mongo.IndexModel{
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetUnique(true),
			},
			mongo.IndexModel{
				Keys: bson.M{"name": 1},
			},
			mongo.IndexModel{
				Keys: bson.M{"own": 1},
			},
		})
	return err
}

// Filebase is a database connection with file operations
type Filebase struct {
	Database
}

// Reserve a fileid, will mutate if fileid not available, returns reserved file id
func (fb *Filebase) Reserve(id types.FileID) (types.FileID, error) {
	var out *types.FileID
	for out == nil {
		timeout := time.Now().Add(time.Hour * 24)
		result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).UpdateOne(fb.ctx, bson.M{
			"id":      id,
			"reserve": bson.M{"$lte": time.Now()},
		}, bson.M{
			"$set": bson.M{"reserve": timeout},
		})
		if err != nil {
			return id, srverror.New(err, 500, "Error F4", "Unable to update id reserve")
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

// Insert file into database. file id must all ready been reserved, see Reserve
func (fb *Filebase) Insert(r types.FileI) error {
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
		return srverror.New(err, 500, "Error F5", "Unable to insert")
	}
	if result.ModifiedCount == 0 {
		return errors.ErrIDNotReserved.Extend("missing fileid")
	}
	return nil
}

// Get file from id
func (fb *Filebase) Get(fid types.FileID) (types.FileI, error) {
	result := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).FindOne(fb.ctx, bson.M{
		"id": fid,
	})
	fd := new(types.FileDecoder)
	if err := result.Decode(fd); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound.Extend(fid.String())
		}
		return nil, srverror.New(err, 500, "Error F6", "Unable to get file")
	}
	f := fd.File()
	err := f.Populate(fb.Owner())
	if err != nil {
		return nil, err
	}
	return f, nil
}

// GetAll files from ids
func (fb *Filebase) GetAll(fids ...types.FileID) ([]types.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"id": bson.M{"$in": fids},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ids := make([]string, 0, len(fids))
			for _, fid := range fids {
				ids = append(ids, fid.String())
			}
			return nil, errors.ErrNoResults.Extend("GetAll files: ").Extend(ids...)
		}
		return nil, srverror.New(err, 500, "Error F3.1", "Unable to get files")
	}
	return fb.decodefiles(cursor)
}

// Update File
func (fb *Filebase) Update(r types.FileI) error {
	result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).ReplaceOne(fb.ctx, bson.M{
		"id": r.GetID(),
	}, r)
	if err != nil {
		return srverror.New(err, 500, "Error F7", "error updating file")
	}
	if result.ModifiedCount == 0 {
		return errors.ErrNotFound.Extend("unable to update:", r.GetID().String())
	}
	return nil
}

// Remove file with id
func (fb *Filebase) Remove(r types.FileID) error {
	result, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).DeleteOne(fb.ctx, bson.M{
		"id": r,
	})
	if err != nil {
		return srverror.New(err, 500, "Error F8", "unable to remove file", r.String())
	}
	if result.DeletedCount == 0 {
		return errors.ErrNotFound.Extend("File id: ", r.String())
	}
	return nil
}

func (fb *Filebase) decodefiles(cursor *mongo.Cursor) ([]types.FileI, error) {
	var reference []*types.FileDecoder
	if err := cursor.All(fb.ctx, &reference); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNoResults.Extend("unable to decode files")
		}
		return nil, srverror.New(err, 500, "Error F9", "unable to decode file list")
	}
	files := make([]types.FileI, 0, len(reference))
	for _, ref := range reference {
		files = append(files, ref.File())
	}
	for _, file := range files {
		err := file.Populate(fb.Owner())
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

// GetOwned returns all files owned by owner id
func (fb *Filebase) GetOwned(uid types.OwnerID) ([]types.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"own": uid,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNoResults.Extend("no owned files")
		}
		return nil, srverror.New(err, 500, "Error F10", "unable to send request")
	}
	return fb.decodefiles(cursor)
}

// GetPermKey returns all files that have owner id with provided permission
func (fb *Filebase) GetPermKey(uid types.OwnerID, pkey string) ([]types.FileI, error) {
	cursor, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).Find(fb.ctx, bson.M{
		"perm." + pkey: uid,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNoResults.Extend("no files with:", pkey)
		}
		return nil, srverror.New(err, 500, "Error F11", "unable to send request")
	}
	return fb.decodefiles(cursor)
}

// Count returns the total number of files accessible by the ownerid, the optional pkeys are permission values to check in addition to ownership
func (fb *Filebase) Count(uid types.OwnerID, pkeys ...string) (int64, error) {
	var or bson.A
	or = append(or, bson.M{
		"own": uid,
	})
	for _, k := range pkeys {
		or = append(or, bson.M{
			"perm." + k: uid,
		})
	}
	count, err := fb.client.Database(fb.DBName).Collection(fb.CollNames["file"]).CountDocuments(fb.ctx, bson.M{
		"$or": or,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return -1, srverror.New(err, 500, "Error F12", "unable to count documents")
	}
	return count, nil
}

// MatchStore returns all files where an owner either owns the file or has a particular permission,
// and the file matches one of the provided StoreIDs
func (fb *Filebase) MatchStore(oid types.OwnerID, sid []types.StoreID, pkeys ...string) ([]types.FileI, error) {
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
			return nil, errors.ErrNoResults.Extend("no filestores match file")
		}
		return nil, srverror.New(err, 500, "Error F13", "unable to send request")
	}
	return fb.decodefiles(cursor)
}
