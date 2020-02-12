package mongo

import (
	"fmt"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Viewbase struct {
	Database
}

func (vb *Viewbase) Insert(vs *database.ViewStore) error {
	_, err := vb.client.Database(vb.DBName).Collection(vb.CollNames["view"]).UpdateOne(
		vb.ctx,
		bson.M{"id": vs.ID},
		bson.M{"$set": vs},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return srverror.New(err, 500, "Database Error V3", "unable to insert viewstore")
	}
	chunks := chunkify(vs.ID, vs.Content)
	_, err = vb.client.Database(vb.DBName).Collection(vb.CollNames["chunk"]).InsertMany(
		vb.ctx,
		chunks,
		options.InsertMany().SetOrdered(false),
	)
	if err != nil {
		return srverror.New(err, 500, "Database Error V4", "failed to insert data chunks")
	}
	return nil
}

func (vb *Viewbase) Get(id filehash.StoreID) (out *database.ViewStore, err error) {
	encodedVS := vb.client.Database(vb.DBName).Collection(vb.CollNames["view"]).FindOne(
		vb.ctx,
		bson.M{"id": id},
	)
	var decodedVS = new(database.ViewStore)
	if err := encodedVS.Decode(decodedVS); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Err V5", "failed to find viewstore")
	}
	var chunks []*contentchunk
	cursor, err := vb.client.Database(vb.DBName).Collection(vb.CollNames["chunk"]).Find(
		vb.ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error V6", "failed to get view data chunks")
	}
	if err = cursor.All(vb.ctx, &chunks); err != nil {
		return nil, srverror.New(err, 500, "Database Error V7", "failed to decode view chunks")
	}
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = srverror.New(v, 500, "Database Error V7", "unable to build chunks")
				out = nil
			default:
				err = srverror.New(fmt.Errorf("GetStore: %+#v", v), 500, "Database Error V8")
				out = nil
			}
		}
	}()
	decodedVS.Content = appendchunks(chunksort(chunks))
	return decodedVS, nil
}
