package mongo

import (
	"context"
	"sync"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Tagbase is a connection to the database with tag operations
type Tagbase struct {
	Database
}

func divideTags(tags []tag.FileTag) ([]tag.StoreTag, []tag.FileTag) {
	stags := make([]tag.StoreTag, 0, len(tags))
	ftags := make([]tag.FileTag, 0, len(tags))
	for _, t := range tags {
		if st := t.StoreTag(); st.Type != 0 {
			stags = append(stags, t.StoreTag())
		}
		if ft := t.Pure(); ft.Type != 0 {
			ftags = append(ftags, t.Pure())
		}
	}
	return stags, ftags
}

// Upsert adds tag to the database
func (tb *Tagbase) Upsert(tags ...tag.FileTag) error {
	stags, ftags := divideTags(tags)
	upsertctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	errch := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(stags) + len(ftags))
	if len(stags) > 0 { // upsert store tags
		storeColl := tb.client.Database(tb.DBName).Collection(tb.CollNames["storetags"])
		for _, st := range stags {
			go func(st tag.StoreTag) {
				defer wg.Done()
				updatefields := bson.M{
					"$setOnInsert": bson.M{
						"store": st.Store,
						"word":  st.Word,
					},
					"$bit": bson.M{"type": bson.M{"or": st.Type}},
				}
				if st.Data != nil {
					set := make(bson.M)
					for t, info := range st.Data {
						for k, v := range info {
							set["data."+t.String()+"."+k] = v
						}
					}
					updatefields["$set"] = set
				}
				_, err := storeColl.UpdateOne(
					upsertctx,
					bson.M{
						"store": st.Store,
						"word":  st.Word,
					},
					updatefields,
					options.Update().SetUpsert(true),
				)
				if err != nil {
					err = srverror.New(err, 500, "Database Error T1.1", "Upserting store tag failed")
					select {
					case errch <- err:
					case <-upsertctx.Done():
					}
				}
			}(st)
		}
	}
	if len(ftags) > 0 { // upsert file tags
		fileColl := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"])
		for _, ft := range ftags {
			go func(ft tag.FileTag) {
				defer wg.Done()
				updatefields := bson.M{
					"$setOnInsert": bson.M{
						"file":  ft.File,
						"owner": ft.Owner,
						"word":  ft.Word,
					},
					"$bit": bson.M{"type": bson.M{"or": ft.Type}},
				}
				if ft.Data != nil {
					set := make(bson.M)
					for t, info := range ft.Data {
						for k, v := range info {
							set["data."+t.String()+"."+k] = v
						}
					}
					updatefields["$set"] = set
				}
				_, err := fileColl.UpdateOne(
					upsertctx,
					bson.M{
						"file":  ft.File,
						"owner": ft.Owner,
						"word":  ft.Word,
					},
					updatefields,
					options.Update().SetUpsert(true),
				)
				if err != nil {
					err = srverror.New(err, 500, "Database Error T1.2", "Upserting file tag failed")
					select {
					case errch <- err:
					case <-upsertctx.Done():
					}
				}
			}(ft)
		}
	}
	go func() {
		wg.Wait()
		close(errch)
	}()
	return <-errch
}

// Remove tags from database
func (tb *Tagbase) Remove(tags ...tag.FileTag) error {
	stags, ftags := divideTags(tags)
	rmctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	errch := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(stags) + len(ftags))
	if len(stags) > 0 {
		storecoll := tb.client.Database(tb.DBName).Collection(tb.CollNames["storetags"])
		for _, st := range stags {
			go func(st tag.StoreTag) {
				defer wg.Done()
				var err error
				if st.Type == tag.ALLSTORE {
					_, err = storecoll.DeleteOne(rmctx, bson.M{
						"word":  st.Word,
						"store": st.Store,
					})
				} else {
					removeData := bson.M{
						"$bit": bson.M{"type": bson.M{"and": ^st.Type}},
					}
					if st.Data != nil {
						unset := make(bson.M)
						for t, mapping := range st.Data {
							if t&st.Type == 0 {
								for k := range mapping {
									unset["data."+t.String()+"."+k] = ""
								}
							} else {
								unset["data."+t.String()] = ""
							}
						}
						removeData["$unset"] = unset
					}
					_, err = storecoll.UpdateOne(
						rmctx,
						bson.M{
							"word":  st.Word,
							"store": st.Store,
						},
						removeData,
					)
				}
				if err != nil {
					select {
					case errch <- srverror.New(err, 500, "Database Error T2.1", "unable to remove storetag"):
					case <-rmctx.Done():
					}
				}
			}(st)
		}
	}
	if len(ftags) > 0 {
		fileColl := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"])
		for _, ft := range ftags {
			go func(ft tag.FileTag) {
				defer wg.Done()
				var err error
				if ft.Type == tag.ALLFILE {
					_, err = fileColl.DeleteOne(rmctx, bson.M{
						"word":  ft.Word,
						"file":  ft.File,
						"owner": ft.Owner,
					})
				} else {
					removeData := bson.M{
						"$bit": bson.M{"type": bson.M{"and": ^ft.Type}},
					}
					if ft.Data != nil {
						unset := make(bson.M)
						for t, mapping := range ft.Data {
							if t&ft.Type == 0 {
								for k := range mapping {
									unset["data."+t.String()+"."+k] = ""
								}
							} else {
								unset["data."+t.String()] = ""
							}
						}
						removeData["$unset"] = unset
					}
					_, err = fileColl.UpdateOne(
						rmctx,
						bson.M{
							"word":  ft.Word,
							"file":  ft.File,
							"owner": ft.Owner,
						},
						removeData,
					)
				}
				if err != nil {
					select {
					case errch <- srverror.New(err, 500, "Database Error T2.2", "unable to remove filetag"):
					case <-rmctx.Done():
					}
				}
			}(ft)
		}
	}
	go func() {
		wg.Wait()
		close(errch)
	}()
	return <-errch
}

func (tb *Tagbase) Get(fid types.FileID, oid types.OwnerID) ([]tag.FileTag, error) {
	type result struct {
		tags []tag.FileTag
		err  error
	}
	out := make(chan result)
	errch := make(chan error)
	storetags := make(chan []tag.StoreTag)
	filetags := make(chan []tag.FileTag)
	getctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	go func() {
		select {
		case err := <-errch:
			select {
			case out <- result{err: err}:
			case <-getctx.Done():
			}
		case <-getctx.Done():
		}
	}()
	go func() { // Get Store Tags
		cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["storetags"]).Find(
			getctx,
			bson.M{
				"store": fid.StoreID,
			},
		)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				select {
				case storetags <- nil:
				case <-getctx.Done():
				}
				return
			}
			select {
			case errch <- srverror.New(err, 500, "Database Error T3.1", "unable to find on storetags"):
			case <-getctx.Done():
			}
			return
		}
		var results []tag.StoreTag
		err = cursor.All(getctx, &results)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				select {
				case storetags <- nil:
				case <-getctx.Done():
				}
				return
			}
			select {
			case errch <- srverror.New(err, 500, "Database Error T3.2", "unable to decode storetags"):
			case <-getctx.Done():
			}
			return
		}
		select {
		case storetags <- results:
		case <-getctx.Done():
		}
	}()
	go func() { // Get File Tags
		cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"]).Find(
			getctx,
			bson.M{
				"file":  fid,
				"owner": oid,
			},
		)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				select {
				case filetags <- nil:
				case <-getctx.Done():
				}
				return
			}
			select {
			case errch <- srverror.New(err, 500, "Database Error T3.3", "unable to find on filetags"):
			case <-getctx.Done():
			}
			return
		}
		var results []tag.FileTag
		err = cursor.All(getctx, &results)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				select {
				case filetags <- nil:
				case <-getctx.Done():
				}
				return
			}
			select {
			case errch <- srverror.New(err, 500, "Database Error T3.4", "unable to decode filetags"):
			case <-getctx.Done():
			}
			return
		}
		select {
		case filetags <- results:
		case <-getctx.Done():
		}
	}()
	go func() {
		collecting := make(map[string]tag.FileTag)
		for i := 0; i < 2; i++ {
			select {
			case stags := <-storetags:
				for _, st := range stags {
					collecting[st.Word] = tag.FileTag{
						Tag:   collecting[st.Word].Update(st.Tag),
						File:  fid,
						Owner: oid,
					}
				}
			case ftags := <-filetags:
				for _, ft := range ftags {
					collecting[ft.Word] = tag.FileTag{
						Tag:   collecting[ft.Word].Update(ft.Tag),
						File:  fid,
						Owner: oid,
					}
				}
			case <-getctx.Done():
				return
			}
		}
		var res result
		res.tags = make([]tag.FileTag, 0, len(collecting))
		for _, v := range collecting {
			res.tags = append(res.tags, v)
		}
		select {
		case out <- res:
		case <-getctx.Done():
		}
	}()

	res := <-out
	return res.tags, res.err
}

func (tb *Tagbase) GetAll(typ tag.Type, oid types.OwnerID) ([]tag.FileTag, error) {}

func (tb *Tagbase) SearchOwned(oid types.OwnerID, tags ...tag.Tag) ([]types.FileID, error) {}

func (tb *Tagbase) SearchAccess(types.OwnerID, string, ...tag.Tag) ([]types.FileID, error) {}

func (tb *Tagbase) SearchFiles(fids []types.FileID, tags ...tag.FileTag) ([]types.FileID, error) {}

// // UpsertFile adds tags associated with file id
// func (tb *Tagbase) UpsertFile(id types.FileID, tags ...tag.Tag) error {
// 	var data []tagbson
// 	for _, t := range tags {
// 		data = append(data, filetag(id, t))
// 	}
// 	upsertctx, cancel := context.WithCancel(tb.ctx)
// 	defer cancel()
// 	errch := make(chan error)
// 	var wg sync.WaitGroup
// 	wg.Add(len(data))
// 	collection := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"])
// 	for _, d := range data {
// 		go func(tag tagbson) {
// 			defer wg.Done()
// 			updateops := bson.M{
// 				"$setOnInsert": bson.M{
// 					"file": tag.File,
// 					"word": strings.ToLower(tag.Word),
// 				},
// 				"$bit": bson.M{"type": bson.M{"or": tag.Type}},
// 			}
// 			if tag.Data != nil {
// 				set := make(bson.M)
// 				for typ, fields := range *tag.Data {
// 					key := "data." + typ.String() + "."
// 					for k, v := range fields {
// 						set[key+k] = v
// 					}
// 				}
// 				updateops["$set"] = set
// 			}
// 			_, err := collection.UpdateOne(
// 				upsertctx,
// 				bson.M{
// 					"file": tag.File,
// 					"word": strings.ToLower(tag.Word),
// 				},
// 				updateops,
// 				options.Update().SetUpsert(true),
// 			)
// 			if err != nil {
// 				errch <- srverror.New(err, 500, "Database Error T1", "unable to upsert file tag")
// 			}
// 		}(d)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(errch)
// 	}()
// 	return <-errch
// }
//
// // UpsertStore adds tags associated with store id
// func (tb *Tagbase) UpsertStore(id types.StoreID, tags ...tag.Tag) error {
// 	var data []tagbson
// 	for _, t := range tags {
// 		data = append(data, storetag(id, t))
// 	}
// 	upsertctx, cancel := context.WithCancel(tb.ctx)
// 	defer cancel()
// 	errch := make(chan error)
// 	var wg sync.WaitGroup
// 	wg.Add(len(data))
// 	collection := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"])
// 	for _, d := range data {
// 		go func(tag tagbson) {
// 			defer wg.Done()
// 			updateops := bson.M{
// 				"$setOnInsert": bson.M{
// 					"store": tag.Store,
// 					"word":  strings.ToLower(tag.Word),
// 				},
// 				"$bit": bson.M{"type": bson.M{"or": tag.Type}},
// 			}
// 			if tag.Data != nil {
// 				set := make(bson.M)
// 				for typ, fields := range *tag.Data {
// 					key := "data." + typ.String() + "."
// 					for k, v := range fields {
// 						set[key+k] = v
// 					}
// 				}
// 				updateops["$set"] = set
// 			}
// 			_, err := collection.UpdateOne(
// 				upsertctx,
// 				bson.M{
// 					"store": tag.Store,
// 					"word":  strings.ToLower(tag.Word),
// 				},
// 				updateops,
// 				options.Update().SetUpsert(true),
// 			)
// 			if err != nil {
// 				errch <- srverror.New(err, 500, "Database Error T2", "unable to upsert store tag")
// 			}
// 		}(d)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(errch)
// 	}()
// 	return <-errch
// }
//
// // FileTags returns all tags associated with file
// func (tb *Tagbase) FileTags(files ...types.FileID) (map[string][]tag.Tag, error) {
// 	stores := make([]types.StoreID, 0, len(files))
// 	for _, f := range files {
// 		stores = append(stores, f.StoreID)
// 	}
// 	var perr error
// 	{
// 		sb := tb.Store(nil)
// 		for _, sid := range stores {
// 			fs, err := sb.Get(sid)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if fs.Perr != nil {
// 				perr = fs.Perr
// 				break
// 			}
// 		}
// 	}
// 	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Find(
// 		tb.ctx,
// 		bson.M{
// 			"$or": bson.A{
// 				bson.M{"file": bson.M{"$in": files}},
// 				bson.M{"store": bson.M{"$in": stores}},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			if perr != nil {
// 				return nil, perr
// 			}
// 			return nil, errors.ErrNoResults.Extend("no tags")
// 		}
// 		return nil, srverror.New(err, 500, "Database Error T3", "unable to find tags")
// 	}
// 	var matches []*tagbson
// 	if err := cursor.All(tb.ctx, &matches); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			if perr != nil {
// 				return nil, perr
// 			}
// 			return nil, errors.ErrNoResults.Extend("no tags decoded")
// 		}
// 		return nil, srverror.New(err, 500, "Database Error T3.1", "unable to decode tags")
// 	}
// 	if len(matches) == 0 {
// 		if perr != nil {
// 			return nil, perr
// 		}
// 		return nil, errors.ErrNoResults.Extend("no tags found for files")
// 	}
// 	out := make(map[string][]tag.Tag)
// 	for _, match := range matches {
// 		if match.File != nil {
// 			out[match.File.String()] = append(out[match.File.String()], match.Tag())
// 		}
// 		if match.Store != nil {
// 			for _, fid := range files {
// 				if fid.StoreID.Equal(*match.Store) {
// 					out[fid.String()] = append(out[fid.String()], match.Tag())
// 				}
// 			}
// 		}
// 	}
// 	return out, perr
// }
//
// type tagAggReturn struct {
// 	Store types.StoreID `bson:"_id"`
// 	Tags  []tagbson     `bson:"tags"`
// }
//
// // GetFiles returns file ids and store ids that have matching tags
// // if any files are given in context, then only tags of those files are searched
// func (tb *Tagbase) GetFiles(filters []tag.Tag, context ...types.FileID) ([]types.FileID, []types.StoreID, error) {
// 	//Build Aggregation Pipeline
// 	aggmatch := make(bson.A, 0, len(filters))
// 	for _, filter := range filters {
// 		match := bson.M{
// 			"word": strings.ToLower(filter.Word),
// 			"type": bson.M{"$bitsAnySet": filter.Type},
// 		}
// 		for typ, fields := range filter.Data {
// 			prefix := "data." + typ.String() + "."
// 			for k, v := range fields {
// 				match[prefix+k] = v
// 			}
// 		}
// 		aggmatch = append(aggmatch, match)
// 	}
// 	initmatch := make(bson.M)
// 	if len(context) == 0 {
// 		initmatch["$or"] = aggmatch
// 	} else {
// 		storeids := make([]types.StoreID, 0, len(context))
// 		for _, fid := range context {
// 			storeids = append(storeids, fid.StoreID)
// 		}
// 		initmatch["$and"] = bson.A{
// 			bson.M{"$or": aggmatch},
// 			bson.M{"$or": bson.A{
// 				bson.M{"file": bson.M{"$in": context}},
// 				bson.M{"store": bson.M{"$in": storeids}},
// 			}}}
// 	}
// 	agg := bson.A{
// 		bson.M{"$match": initmatch},
// 		bson.M{"$group": bson.M{
// 			"_id": bson.M{"$ifNull": bson.A{"$store", "$file.storeid"}},
// 			"tags": bson.M{"$push": bson.M{
// 				"file": "$file",
// 				"word": "$word",
// 				"type": "$type",
// 				"data": "$data",
// 			}},
// 		}},
// 	}
// 	for _, filter := range aggmatch {
// 		agg = append(agg, bson.M{
// 			"$match": bson.M{
// 				"tags": bson.M{"$elemMatch": filter},
// 			},
// 		})
// 	}
// 	//Run Aggregation
// 	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Aggregate(
// 		tb.ctx,
// 		agg,
// 	)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil, errors.ErrNoResults.Extend("no files match tags")
// 		}
// 		return nil, nil, srverror.New(err, 500, "Database Error T4.1", "unable to get aggregate tags")
// 	}
// 	var results []tagAggReturn
// 	if err := cursor.All(tb.ctx, &results); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil, errors.ErrNoResults.Extend("no files decoded matching tags")
// 		}
// 		return nil, nil, srverror.New(err, 500, "Database Error T4", "unable to decode data")
// 	}
// 	//assemble fileids and storeids for return; check that fileids meet all filter conditions
// 	var fids []types.FileID
// 	var sids []types.StoreID
// 	for _, result := range results {
// 		sids = append(sids, result.Store)
// 	NEXTTAG:
// 		for _, tag := range result.Tags {
// 			if tag.File != nil && func() bool {
// 				for _, foundid := range fids {
// 					if tag.File.Equal(foundid) {
// 						return false
// 					}
// 				}
// 				return true
// 			}() {
// 			NEXTFILTER:
// 				for _, filter := range filters {
// 					for _, t := range result.Tags {
// 						if (t.File == nil || tag.File.Equal(*t.File)) && t.Word == filter.Word && t.Type&filter.Type != 0 && func() bool {
// 							for typ, fields := range filter.Data {
// 								for k, v := range fields {
// 									if t.Data == nil || *t.Data == nil || (*t.Data)[typ] == nil || (*t.Data)[typ][k] != v {
// 										return false
// 									}
// 								}
// 							}
// 							return true
// 						}() {
// 							continue NEXTFILTER
// 						}
// 					}
// 					continue NEXTTAG
// 				}
// 				fids = append(fids, *tag.File)
// 			}
// 		}
// 	}
// 	//If context provided add fileids that match to store ids
// 	if len(context) > 0 {
// 		for _, id := range context {
// 			if func() bool {
// 				for _, s := range sids {
// 					if s.Equal(id.StoreID) {
// 						return true
// 					}
// 				}
// 				return false
// 			}() && func() bool {
// 				for _, fid := range fids {
// 					if fid.Equal(id) {
// 						return false
// 					}
// 				}
// 				return true
// 			}() {
// 				fids = append(fids, id)
// 			}
// 		}
// 	}
// 	if len(fids) == 0 && len(sids) == 0 {
// 		return nil, nil, errors.ErrNoResults.Extend("no full matching files")
// 	}
// 	return fids, sids, nil
// }
//
// // SearchData returns tags that contain a particular type and matching data fields
// func (tb *Tagbase) SearchData(typ tag.Type, data tag.Data) ([]tag.Tag, error) {
// 	//TODO: bit-wise-or the keys of data rather then pass tag.Type
// 	filter := make(bson.M)
// 	filter["type"] = bson.M{"$bitsAnySet": typ}
// 	for t, m := range data {
// 		prefix := "data." + t.String() + "."
// 		for k, v := range m {
// 			filter[prefix+k] = v
// 		}
// 	}
// 	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["tag"]).Find(
// 		tb.ctx,
// 		filter,
// 	)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, errors.ErrNoResults.Extend("no matching tags")
// 		}
// 		return nil, srverror.New(err, 500, "Database Error T5", "tag.searchData mongo error")
// 	}
// 	var returned []tagbson
// 	if err := cursor.All(tb.ctx, &returned); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, errors.ErrNoResults.Extend("no decoded tags")
// 		}
// 		return nil, srverror.New(err, 500, "Database Error T5.1", "tag.searchData decode error")
// 	}
// 	result := make([]tag.Tag, 0, len(data))
// 	for _, d := range returned {
// 		result = append(result, d.Tag())
// 	}
// 	return result, nil
// }
