// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"context"
	"sync"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initFileTagsIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["filetags"]).Indexes()
	_, err := I.CreateMany(
		ctx,
		[]mongo.IndexModel{
			mongo.IndexModel{
				Keys: bson.D{
					bson.E{Key: "owner", Value: 1},
					bson.E{Key: "word", Value: 1},
					bson.E{Key: "file", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
			mongo.IndexModel{
				Keys: bson.M{"word": 1},
			},
		})
	return err
}

func initStoreTagIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["storetags"]).Indexes()
	_, err := I.CreateMany(
		ctx,
		[]mongo.IndexModel{
			mongo.IndexModel{
				Keys:    bson.D{bson.E{Key: "store", Value: 1}, bson.E{Key: "word", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			mongo.IndexModel{
				Keys: bson.M{"word": 1},
			},
		})
	return err
}

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
					err = srverror.New(err, 500, "Error T1.1", "Upserting store tag failed")
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
					err = srverror.New(err, 500, "Error T1.2", "Upserting file tag failed")
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
					case errch <- srverror.New(err, 500, "Error T2.1", "unable to remove storetag"):
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
					case errch <- srverror.New(err, 500, "Error T2.2", "unable to remove filetag"):
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

// Get returns all filetags for a particular file id and owner
func (tb *Tagbase) Get(fid types.FileID, oid types.OwnerID) ([]tag.FileTag, error) {
	return tb.GetType(fid, oid, tag.ALLTYPES)
}

// GetType returns FileTags of a particular file, for a particular owner, with a match to a particular type
func (tb *Tagbase) GetType(fid types.FileID, oid types.OwnerID, typ tag.Type) ([]tag.FileTag, error) {
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
		if typ&tag.ALLSTORE == 0 {
			select {
			case storetags <- nil:
			case <-getctx.Done():
			}
			return
		}
		cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["storetags"]).Find(
			getctx,
			bson.M{
				"store": fid.StoreID,
				"type": bson.M{
					"$bitsAnySet": typ,
				},
			},
		)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				select {
				case errch <- nil:
				case <-getctx.Done():
				}
				return
			}
			select {
			case errch <- srverror.New(err, 500, "Error T3.1", "unable to find on storetags"):
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
			case errch <- srverror.New(err, 500, "Error T3.2", "unable to decode storetags"):
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
		if typ&tag.ALLFILE == 0 {
			select {
			case filetags <- nil:
			case <-getctx.Done():
			}
			return
		}
		cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"]).Find(
			getctx,
			bson.M{
				"file":  fid,
				"owner": oid,
				"type": bson.M{
					"$bitsAnySet": typ,
				},
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
			case errch <- srverror.New(err, 500, "Error T3.3", "unable to find on filetags"):
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
			case errch <- srverror.New(err, 500, "Error T3.4", "unable to decode filetags"):
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
						Tag:   collecting[st.Word].Tag.Update(st.Tag),
						File:  fid,
						Owner: oid,
					}
				}
			case ftags := <-filetags:
				for _, ft := range ftags {
					collecting[ft.Word] = tag.FileTag{
						Tag:   collecting[ft.Word].Tag.Update(ft.Tag),
						File:  fid,
						Owner: oid,
					}
				}
			case <-getctx.Done():
				return
			}
		}
		if len(collecting) == 0 {
			select {
			case errch <- errors.ErrNoResults.Extend("no tags"):
			case <-getctx.Done():
			}
			return
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

// GetAll returns all FileTags that have a particular type and and owner.
// Since it looks for owner association, store tags are not searched
func (tb *Tagbase) GetAll(typ tag.Type, oid types.OwnerID) ([]tag.FileTag, error) {
	cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"]).Find(
		tb.ctx,
		bson.M{
			"type":  bson.M{"$bitsAnySet": typ},
			"owner": oid,
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNoResults.Extend("no file tags")
		}
		return nil, srverror.New(err, 500, "Error T4.1", "unable to get file tags")
	}
	var results []tag.FileTag
	err = cursor.All(tb.ctx, &results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNoResults.Extend("no file tags")
		}
		return nil, srverror.New(err, 500, "Error T4.2", "unable to decode file tags")
	}
	return results, nil
}

// SearchOwned returns all fileids that are owned by the owner and match the filtering tags
func (tb *Tagbase) SearchOwned(oid types.OwnerID, tags ...tag.FileTag) ([]types.FileID, error) {
	fb := tb.File()
	files, err := fb.GetOwned(oid)
	if err != nil {
		return nil, err
	}
	fids := make([]types.FileID, 0, len(files))
	for _, f := range files {
		fids = append(fids, f.GetID())
	}
	return tb.SearchFiles(fids, tags...)
}

// SearchAccess returns all fileids that the owner has access to in a particular permission and match the filtering tags
func (tb *Tagbase) SearchAccess(oid types.OwnerID, pkey string, tags ...tag.FileTag) ([]types.FileID, error) {
	fb := tb.File()
	files, err := fb.GetPermKey(oid, pkey)
	if err != nil {
		return nil, err
	}
	fids := make([]types.FileID, 0, len(files))
	for _, f := range files {
		fids = append(fids, f.GetID())
	}
	return tb.SearchFiles(fids, tags...)
}

// SearchFiles returns all the fileids that match the tags which define the filter conditions
func (tb *Tagbase) SearchFiles(fids []types.FileID, tags ...tag.FileTag) ([]types.FileID, error) {
	if len(tags) == 0 {
		return fids, nil
	}
	var searchStoreTags bool
	var searchFileTags bool
	for _, t := range tags {
		if t.Type&tag.ALLSTORE > 0 {
			searchStoreTags = true
		}
		if t.Type&tag.ALLFILE > 0 {
			searchFileTags = true
		}
	}
	if !searchStoreTags && !searchFileTags {
		return fids, nil
	}
	type result struct {
		ids []types.FileID
		err error
	}
	type storeagg struct {
		Store types.StoreID `bson:"_id"`
		Tags  []tag.Tag     `bson:"tags"`
	}
	type ftag struct {
		Owner types.OwnerID `bson:"owner"`
		Tag   tag.Tag       `bson:",inline"`
	}
	type fileagg struct {
		File types.FileID `bson:"_id"`
		Tags []ftag       `bson:"tags"`
	}
	out := make(chan result)
	errch := make(chan error)
	storeresults := make(chan []storeagg)
	fileresults := make(chan []fileagg)
	searchctx, cancel := context.WithCancel(tb.ctx)
	defer cancel()
	go func() {
		select {
		case err := <-errch:
			select {
			case out <- result{err: err}:
			case <-searchctx.Done():
			}
		case <-searchctx.Done():
		}
	}()

	if searchStoreTags {
		go func() { // Search StoreTags
			words := make([]string, 0, len(tags))
			regexs := make([]bson.M, 0, len(tags))
			for _, t := range tags {
				if t.Type&tag.ALLSTORE > 0 {
					if t.Type&tag.SEARCH > 0 {
						if data := t.Data[tag.SEARCH]; data != nil {
							if searchwithregex, ok := data["regex"].(bool); ok && searchwithregex {
								regexoptions, _ := data["regexoptions"].(string)
								regexs = append(regexs, bson.M{
									"word": bson.M{
										"$regex":   t.Word,
										"$options": regexoptions,
									},
								})
								continue
							}
						}
					}
					words = append(words, t.Word)
				}
			}
			sids := make([]types.StoreID, 0, len(fids))
			for _, id := range fids {
				sids = append(sids, id.StoreID)
			}
			wordConditions := []bson.M{
				bson.M{"word": bson.M{"$in": words}},
			}
			if len(regexs) > 0 {
				wordConditions = append(wordConditions, bson.M{"$or": regexs})
			}
			pipeline := []bson.M{
				bson.M{
					"$match": bson.M{
						"$or":   wordConditions,
						"store": bson.M{"$in": sids},
					},
				},
				bson.M{
					"$group": bson.M{
						"_id": "$store",
						"tags": bson.M{
							"$push": bson.M{
								"word": "$word",
								"type": "$type",
							},
						},
					},
				},
			}
			for _, t := range tags {
				if t.Type&tag.ALLSTORE > 0 {
					match := make(bson.M)
					isRegex := false
					if t.Type&tag.SEARCH > 0 {
						if data := t.Data[tag.SEARCH]; data != nil {
							if searchwithregex, ok := data["regex"].(bool); ok && searchwithregex {
								regexoptions, _ := data["regexoptions"].(string)
								match["word"] = bson.M{
									"$regex":   t.Word,
									"$options": regexoptions,
								}
								isRegex = true
							}
						}
					}
					if !isRegex {
						match["word"] = t.Word
					}
					match["type"] = bson.M{
						"$bitsAnySet": t.Type,
					}
					pipeline = append(pipeline, bson.M{
						"$match": bson.M{
							"tags": bson.M{
								"$elemMatch": match,
							},
						},
					})
				}
			}
			cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["storetags"]).Aggregate(searchctx, pipeline)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					select {
					case storeresults <- nil:
					case <-searchctx.Done():
					}
					return
				}
				select {
				case errch <- srverror.New(err, 500, "Error T5.1", "unable to aggregate store tags"):
				case <-searchctx.Done():
				}
				return
			}
			var results []storeagg
			err = cursor.All(searchctx, &results)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					select {
					case storeresults <- nil:
					case <-searchctx.Done():
					}
					return
				}
				select {
				case errch <- srverror.New(err, 500, "Error T5.2", "unable to decode store tags"):
				case <-searchctx.Done():
				}
				return
			}
			select {
			case storeresults <- results:
			case <-searchctx.Done():
			}
		}()
	} else {
		go func() {
			select {
			case storeresults <- nil:
			case <-searchctx.Done():
			}
		}()
	}
	if searchFileTags {
		go func() { // Search File Tags
			words := make([]string, 0, len(tags))
			regexs := make([]bson.M, 0, len(tags))
			for _, t := range tags {
				if t.Type&tag.ALLFILE > 0 {
					if t.Type&tag.SEARCH > 0 {
						if data := t.Data[tag.SEARCH]; data != nil {
							if searchwithregex, ok := data["regex"].(bool); ok && searchwithregex {
								regexoptions, _ := data["regexoptions"].(string)
								regexs = append(regexs, bson.M{
									"word": bson.M{
										"$regex":   t.Word,
										"$options": regexoptions,
									},
								})
								continue
							}
						}
					}
					words = append(words, t.Word)
				}
			}
			wordConditions := []bson.M{
				bson.M{"word": bson.M{"$in": words}},
			}
			if len(regexs) > 0 {
				wordConditions = append(wordConditions, bson.M{
					"$or": regexs,
				})
			}
			pipeline := []bson.M{
				bson.M{
					"$match": bson.M{
						"$or":  wordConditions,
						"file": bson.M{"$in": fids},
					},
				},
				bson.M{
					"$group": bson.M{
						"_id": "$file",
						"tags": bson.M{
							"$push": bson.M{
								"owner": "$owner",
								"word":  "$word",
								"type":  "$type",
							},
						},
					},
				},
			}
			for _, t := range tags {
				if t.Type&tag.ALLFILE > 0 {
					match := make(bson.M)
					isRegex := false
					if t.Type&tag.SEARCH > 0 {
						if data := t.Data[tag.SEARCH]; data != nil {
							if searchwithregex, ok := data["regex"].(bool); ok && searchwithregex {
								regexoptions, _ := data["regexoptions"].(string)
								match["word"] = bson.M{
									"$regex":   t.Word,
									"$options": regexoptions,
								}
								isRegex = true
							}
						}
					}
					if !isRegex {
						match["word"] = t.Word
					}
					match["type"] = bson.M{
						"$bitsAnySet": t.Type,
					}
					if !t.Owner.Equal(types.OwnerID{}) {
						match["owner"] = t.Owner
					}
					pipeline = append(pipeline, bson.M{
						"$match": bson.M{
							"tags": bson.M{
								"$elemMatch": match,
							},
						},
					})
				}
			}
			cursor, err := tb.client.Database(tb.DBName).Collection(tb.CollNames["filetags"]).Aggregate(searchctx, pipeline)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					select {
					case fileresults <- nil:
					case <-searchctx.Done():
					}
					return
				}
				select {
				case errch <- srverror.New(err, 500, "Error T5.3", "unable to aggregate file tags"):
				case <-searchctx.Done():
				}
				return
			}
			var results []fileagg
			err = cursor.All(searchctx, &results)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					select {
					case fileresults <- nil:
					case <-searchctx.Done():
					}
					return
				}
				select {
				case errch <- srverror.New(err, 500, "Error T5.4", "unable to decode file tags"):
				case <-searchctx.Done():
				}
				return
			}
			select {
			case fileresults <- results:
			case <-searchctx.Done():
			}
		}()
	} else {
		go func() {
			select {
			case fileresults <- nil:
			case <-searchctx.Done():
			}
		}()
	}
	go func() {
		var stores []storeagg
		var files []fileagg
		for i := 0; i < 2; i++ {
			select {
			case stores = <-storeresults:
			case files = <-fileresults:
			case <-searchctx.Done():
				return
			}
		}
		var res result
		if searchStoreTags && searchFileTags {
			for _, file := range files {
				for _, store := range stores {
					if store.Store.Equal(file.File.StoreID) {
						res.ids = append(res.ids, file.File)
						break
					}
				}
			}
		} else if searchStoreTags {
			for _, fid := range fids {
				for _, store := range stores {
					if store.Store.Equal(fid.StoreID) {
						res.ids = append(res.ids, fid)
						break
					}
				}
			}
		} else if searchFileTags {
			for _, file := range files {
				res.ids = append(res.ids, file.File)
			}
		}
		if len(res.ids) == 0 {
			res.err = errors.ErrNoResults.Extend("no matching file ids")
		}
		select {
		case out <- res:
		case <-searchctx.Done():
		}
	}()
	res := <-out
	return res.ids, res.err
}
