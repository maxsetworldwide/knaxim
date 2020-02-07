package mongo

import (
	"context"
	"strings"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tagbson struct {
	File  *filehash.FileID  `bson:"file,omitempty" json:"file,omitempty"`
	Store *filehash.StoreID `bson:"store,omitempty" json:"store,omitempty"`
	Word  string            `bson:"word" json:"word"`
	Type  tag.Type          `bson:"type" json:"type"`
	Data  *tag.Data         `bson:"data,omitempty" json:"data,omitempty"`
}

func filetag(f filehash.FileID, t tag.Tag) tagbson {
	r := tagbson{
		File: &f,
		Word: t.Word,
		Type: t.Type,
	}
	if t.Data != nil {
		r.Data = &t.Data
	}
	return r
}

func (tb *tagbson) Tag() tag.Tag {
	t := tag.Tag{
		Type: tb.Type,
		Word: tb.Word,
	}
	if tb.Data != nil {
		t.Data = *tb.Data
	}
	return t
}

func storetag(s filehash.StoreID, t tag.Tag) tagbson {
	r := tagbson{
		Store: &s,
		Word:  t.Word,
		Type:  t.Type,
	}
	if t.Data != nil {
		r.Data = &t.Data
	}
	return r
}

type Tagbase struct {
	Database
}

func (tb *Tagbase) UpsertFile(id filehash.FileID, tags ...tag.Tag) error {
	var data []tagbson
	for _, t := range tags {
		data = append(data, filetag(id, t))
	}
	upsertctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	errch := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(data))
	collection := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"])
	for _, d := range data {
		go func(tag tagbson) {
			defer wg.Done()
			updateops := bson.M{
				"$setOnInsert": bson.M{
					"file": tag.File,
					"word": strings.ToLower(tag.Word),
				},
				"$bit": bson.M{"type": bson.M{"or": tag.Type}},
			}
			if tag.Data != nil {
				set := make(bson.M)
				for typ, fields := range *tag.Data {
					key := "data." + typ.String() + "."
					for k, v := range fields {
						set[key+k] = v
					}
				}
				updateops["$set"] = set
			}
			_, err := collection.UpdateOne(
				upsertctx,
				bson.M{
					"file": tag.File,
					"word": strings.ToLower(tag.Word),
				},
				updateops,
				options.Update().SetUpsert(true),
			)
			if err != nil {
				errch <- srverror.New(err, 500, "Database Error T1", "unable to upsert file tag")
			}
		}(d)
	}
	go func() {
		wg.Wait()
		close(errch)
	}()
	return <-errch
}

func (tb *Tagbase) UpsertStore(id filehash.StoreID, tags ...tag.Tag) error {
	var data []tagbson
	for _, t := range tags {
		data = append(data, storetag(id, t))
	}
	upsertctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	errch := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(data))
	collection := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"])
	for _, d := range data {
		go func(tag tagbson) {
			defer wg.Done()
			updateops := bson.M{
				"$setOnInsert": bson.M{
					"store": tag.Store,
					"word":  strings.ToLower(tag.Word),
				},
				"$bit": bson.M{"type": bson.M{"or": tag.Type}},
			}
			if tag.Data != nil {
				set := make(bson.M)
				for typ, fields := range *tag.Data {
					key := "data." + typ.String() + "."
					for k, v := range fields {
						set[key+k] = v
					}
				}
				updateops["$set"] = set
			}
			_, err := collection.UpdateOne(
				upsertctx,
				bson.M{
					"store": tag.Store,
					"word":  strings.ToLower(tag.Word),
				},
				updateops,
				options.Update().SetUpsert(true),
			)
			if err != nil {
				errch <- srverror.New(err, 500, "Database Error T2", "unable to upsert store tag")
			}
		}(d)
	}
	go func() {
		wg.Wait()
		close(errch)
	}()
	return <-errch
}

func (tb *Tagbase) FileTags(files ...filehash.FileID) (map[string][]tag.Tag, error) {
	stores := make([]filehash.StoreID, 0, len(files))
	for _, f := range files {
		stores = append(stores, f.StoreID)
	}
<<<<<<< Updated upstream
=======
	var perr error
	{
		sb := tb.Store(nil)
		for _, sid := range stores {
			fs, err := sb.Get(sid)
			if err != nil {
				return nil, err
			}
			if fs.Perr != nil {
				perr = fs.Perr
				break
			}
		}
	}
>>>>>>> Stashed changes
	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Find(
		tb.ctx,
		bson.M{
			"$or": bson.A{
				bson.M{"file": bson.M{"$in": files}},
				bson.M{"store": bson.M{"$in": stores}},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error T3", "unable to find tags")
	}
	var matches []*tagbson
	if err := cursor.All(tb.ctx, &matches); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error T3.1", "unable to decode tags")
	}
	if len(matches) == 0 {
		return nil, database.ErrNotFound
	}
	out := make(map[string][]tag.Tag)
	for _, match := range matches {
		if match.File != nil {
			out[match.File.String()] = append(out[match.File.String()], match.Tag())
		}
		if match.Store != nil {
			for _, fid := range files {
				if fid.StoreID.Equal(*match.Store) {
					out[fid.String()] = append(out[fid.String()], match.Tag())
				}
			}
		}
	}
	return out, nil
}

type tagAggReturn struct {
	Store filehash.StoreID `bson:"_id"`
	Tags  []tagbson        `bson:"tags"`
}

func (tb *Tagbase) GetFiles(filters []tag.Tag, context ...filehash.FileID) ([]filehash.FileID, []filehash.StoreID, error) {
	//Build Aggregation Pipeline
	aggmatch := make(bson.A, 0, len(filters))
	for _, filter := range filters {
		match := bson.M{
			"word": strings.ToLower(filter.Word),
			"type": bson.M{"$bitsAnySet": filter.Type},
		}
		for typ, fields := range filter.Data {
			prefix := "data." + typ.String() + "."
			for k, v := range fields {
				match[prefix+k] = v
			}
		}
		aggmatch = append(aggmatch, match)
	}
	initmatch := make(bson.M)
	if len(context) == 0 {
		initmatch["$or"] = aggmatch
	} else {
		storeids := make([]filehash.StoreID, 0, len(context))
		for _, fid := range context {
			storeids = append(storeids, fid.StoreID)
		}
		initmatch["$and"] = bson.A{
			bson.M{"$or": aggmatch},
			bson.M{"$or": bson.A{
				bson.M{"file": bson.M{"$in": context}},
				bson.M{"store": bson.M{"$in": storeids}},
			}}}
	}
	agg := bson.A{
		bson.M{"$match": initmatch},
		bson.M{"$group": bson.M{
			"_id": bson.M{"$ifNull": bson.A{"$store", "$file.storeid"}},
			"tags": bson.M{"$push": bson.M{
				"file": "$file",
				"word": "$word",
				"type": "$type",
				"data": "$data",
			}},
		}},
	}
	for _, filter := range aggmatch {
		agg = append(agg, bson.M{
			"$match": bson.M{
				"tags": bson.M{"$elemMatch": filter},
			},
		})
	}
	//Run Aggregation
	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Aggregate(
		tb.ctx,
		agg,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, database.ErrNotFound
		}
		return nil, nil, srverror.New(err, 500, "Database Error T4.1", "unable to get aggregate tags")
	}
	var results []tagAggReturn
	if err := cursor.All(tb.ctx, &results); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, database.ErrNotFound
		}
		return nil, nil, srverror.New(err, 500, "Database Error T4", "unable to decode data")
	}
	//assemble fileids and storeids for return; check that fileids meet all filter conditions
	var fids []filehash.FileID
	var sids []filehash.StoreID
	for _, result := range results {
		sids = append(sids, result.Store)
	NEXTTAG:
		for _, tag := range result.Tags {
			if tag.File != nil && func() bool {
				for _, foundid := range fids {
					if tag.File.Equal(foundid) {
						return false
					}
				}
				return true
			}() {
			NEXTFILTER:
				for _, filter := range filters {
					for _, t := range result.Tags {
						if (t.File == nil || tag.File.Equal(*t.File)) && t.Word == filter.Word && t.Type&filter.Type != 0 && func() bool {
							for typ, fields := range filter.Data {
								for k, v := range fields {
									if t.Data == nil || *t.Data == nil || (*t.Data)[typ] == nil || (*t.Data)[typ][k] != v {
										return false
									}
								}
							}
							return true
						}() {
							continue NEXTFILTER
						}
					}
					continue NEXTTAG
				}
				fids = append(fids, *tag.File)
			}
		}
	}
	//If context provided add fileids that match to store ids
	if len(context) > 0 {
		for _, id := range context {
			if func() bool {
				for _, s := range sids {
					if s.Equal(id.StoreID) {
						return true
					}
				}
				return false
			}() && func() bool {
				for _, fid := range fids {
					if fid.Equal(id) {
						return false
					}
				}
				return true
			}() {
				fids = append(fids, id)
			}
		}
	}
	if len(fids) == 0 && len(sids) == 0 {
		return nil, nil, database.ErrNoResults
	}
	return fids, sids, nil
}

//TODO: bit-wise-or the keys of data rather then pass tag.Type
func (tb *Tagbase) SearchData(typ tag.Type, data tag.Data) ([]tag.Tag, error) {
	filter := make(bson.M)
	filter["type"] = bson.M{"$bitsAnySet": typ}
	for t, m := range data {
		prefix := "data." + t.String() + "."
		for k, v := range m {
			filter[prefix+k] = v
		}
	}
	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Find(
		tb.ctx,
		filter,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error T5", "tag.searchData mongo error")
	}
	var returned []tagbson
	if err := cursor.All(tb.ctx, &returned); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound
		}
		return nil, srverror.New(err, 500, "Database Error T5.1", "tag.searchData decode error")
	}
	result := make([]tag.Tag, 0, len(data))
	for _, d := range returned {
		result = append(result, d.Tag())
	}
	return result, nil
}
