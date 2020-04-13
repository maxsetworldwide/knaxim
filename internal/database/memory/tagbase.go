package memory

import (
	"regexp"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

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

// Tagbase wraps database and provides tag operations
type Tagbase struct {
	Database
}

// Upsert adds tags to the database
func (tb *Tagbase) Upsert(tags ...tag.FileTag) error {
	lock.Lock()
	defer lock.Unlock()
	stags, ftags := divideTags(tags)
	for _, st := range stags {
		if tb.TagStores[st.Store.String()] == nil {
			tb.TagStores[st.Store.String()] = map[string]tag.StoreTag{
				st.Word: st,
			}
		} else {
			tb.TagStores[st.Store.String()][st.Word] = tb.TagStores[st.Store.String()][st.Word].Update(st)
		}
	}
	for _, ft := range ftags {
		if tb.TagFiles[ft.File.String()] == nil {
			tb.TagFiles[ft.File.String()] = map[string]map[string]tag.FileTag{
				ft.Owner.String(): map[string]tag.FileTag{
					ft.Word: ft,
				},
			}
		} else if tb.TagFiles[ft.File.String()][ft.Owner.String()] == nil {
			tb.TagFiles[ft.File.String()][ft.Owner.String()] = map[string]tag.FileTag{
				ft.Word: ft,
			}
		} else {
			tb.TagFiles[ft.File.String()][ft.Owner.String()][ft.Word] = tb.TagFiles[ft.File.String()][ft.Owner.String()][ft.Word].Update(ft)
		}
	}
	return nil
}

// Remove removes tags from the database
func (tb *Tagbase) Remove(tags ...tag.FileTag) error {
	lock.Lock()
	defer lock.Unlock()
	stags, ftags := divideTags(tags)
	for _, st := range stags {
		if tb.TagStores[st.Store.String()] != nil {
			if old, set := tb.TagStores[st.Store.String()][st.Word]; set {
				if old.Tag.Type&^st.Tag.Type == 0 {
					delete(tb.TagStores[st.Store.String()], st.Word)
				} else {
					updated := tag.StoreTag{
						Store: old.Store,
						Tag: tag.Tag{
							Word: old.Word,
							Type: old.Tag.Type &^ st.Tag.Type,
							Data: old.Tag.Data.FilterType(old.Tag.Type &^ st.Tag.Type),
						},
					}
					tb.TagStores[st.Store.String()][st.Word] = updated
				}
			}
		}
	}
	for _, ft := range ftags {
		if tb.TagFiles[ft.File.String()] != nil && tb.TagFiles[ft.File.String()][ft.Owner.String()] != nil {
			if old, set := tb.TagFiles[ft.File.String()][ft.Owner.String()][ft.Word]; set {
				if old.Tag.Type&^ft.Tag.Type == 0 {
					delete(tb.TagFiles[ft.File.String()][ft.Owner.String()], ft.Word)
				} else {
					tb.TagFiles[ft.File.String()][ft.Owner.String()][ft.Word] = tag.FileTag{
						File:  old.File,
						Owner: old.Owner,
						Tag: tag.Tag{
							Word: old.Word,
							Type: old.Tag.Type &^ ft.Tag.Type,
							Data: old.Tag.Data.FilterType(old.Tag.Type &^ ft.Tag.Type),
						},
					}
				}
			}
		}
	}
	return nil
}

// Get returns all tags associated with a particular file and owner
func (tb *Tagbase) Get(fid types.FileID, oid types.OwnerID) ([]tag.FileTag, error) {
	lock.RLock()
	defer lock.RUnlock()
	ftags := tb.TagFiles[fid.String()][oid.String()]
	stags := tb.TagStores[fid.StoreID.String()]
	words := map[string]bool{}
	var out []tag.FileTag
	for word, ft := range ftags {
		if st, set := stags[word]; set {
			ft.Tag = ft.Tag.Update(st.Tag)
		}
		words[word] = true
		out = append(out, ft)
	}
	for word, st := range stags {
		if !words[word] {
			out = append(out, tag.FileTag{
				File:  fid,
				Owner: oid,
				Tag:   st.Tag,
			})
		}
	}
	return out, nil
}

// GetType returns all tags of a particular type, associated with a particular file and owner
func (tb *Tagbase) GetType(fid types.FileID, oid types.OwnerID, typ tag.Type) (tags []tag.FileTag, err error) {
	lock.RLock()
	defer lock.RUnlock()
	if typ&tag.ALLFILE != 0 {
		if tb.TagFiles[fid.String()] != nil && tb.TagFiles[fid.String()][oid.String()] != nil {
			for _, t := range tb.TagFiles[fid.String()][oid.String()] {
				if t.Type&typ != 0 {
					tags = append(tags, t)
				}
			}
		}
	}
	if typ&tag.ALLSTORE != 0 {
		if tb.TagStores[fid.StoreID.String()] != nil {
			for _, t := range tb.TagStores[fid.StoreID.String()] {
				if t.Type&typ != 0 {
					tags = append(tags, tag.FileTag{
						File:  fid,
						Owner: oid,
						Tag:   t.Tag,
					})
				}
			}
		}
	}
	return
}

// GetAll returns all tags of a particular type for a particular owner
func (tb *Tagbase) GetAll(typ tag.Type, oid types.OwnerID) (tags []tag.FileTag, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, maps := range tb.TagFiles {
		if maps[oid.String()] != nil {
			for _, ft := range maps[oid.String()] {
				if ft.Type&typ != 0 {
					tags = append(tags, ft)
				}
			}
		}
	}
	return
}

// SearchOwned returns all fileids that is owned by the owner and matches the tag fileter conditions
func (tb *Tagbase) SearchOwned(oid types.OwnerID, tags ...tag.FileTag) ([]types.FileID, error) {
	lock.Lock()
	defer func() {
		if r := recover(); r != nil {
			lock.Unlock()
			panic(r)
		}
	}()
	fb := tb.file(nil).(*Filebase)
	fs, err := fb.getOwned(oid)
	defer fb.Close(nil)
	if err != nil {
		return nil, err
	}
	var fids []types.FileID
	for _, f := range fs {
		fids = append(fids, f.GetID())
	}
	lock.Unlock()
	return tb.SearchFiles(fids, tags...)
}

// SearchAccess returns all fileids that are accessable by owner with particular permission that match the tag filter conditions
func (tb *Tagbase) SearchAccess(oid types.OwnerID, key string, tags ...tag.FileTag) ([]types.FileID, error) {
	lock.Lock()
	defer func() {
		if r := recover(); r != nil {
			lock.Unlock()
			panic(r)
		}
	}()
	fb := tb.file(nil).(*Filebase)
	fs, err := fb.getPermKey(oid, key)
	defer fb.Close(nil)
	if err != nil {
		return nil, err
	}
	var fids []types.FileID
	for _, f := range fs {
		fids = append(fids, f.GetID())
	}
	lock.Unlock()
	return tb.SearchFiles(fids, tags...)
}

// SearchFiles returns all fileids that match the tag fileters
func (tb *Tagbase) SearchFiles(in []types.FileID, tags ...tag.FileTag) (out []types.FileID, err error) {
	lock.RLock()
	defer lock.RUnlock()
	var expectFileTag, expectStoreTag bool
	matched := make(map[string]bool)
	for _, t := range tags {
		if t.Type&tag.ALLFILE != 0 {
			expectFileTag = true
		}
		if t.Type&tag.ALLSTORE != 0 {
			expectStoreTag = true
		}
	}
	for _, fid := range in {
		valid := make([]bool, len(tags))
		if expectFileTag {

			for _, maps := range tb.TagFiles[fid.String()] {
			TAG:
				for i, t := range tags {
					if t.Type&tag.ALLFILE != 0 {
						for _, ft := range maps {
							if ft.Type&t.Type != 0 {
								if t.Type&tag.SEARCH != 0 &&
									t.Data[tag.SEARCH] != nil &&
									t.Data[tag.SEARCH]["regex"] == true {
									if matched, _ := regexp.MatchString(t.Word, ft.Word); matched {
										valid[i] = true
										continue TAG
									}
								} else {
									if t.Word == ft.Word {
										valid[i] = true
										continue TAG
									}
								}
							}
						}
					}
				}
			}
		}
		if expectStoreTag {
		STAG:
			for i, t := range tags {
				for _, st := range tb.TagStores[fid.StoreID.String()] {
					if st.Type&t.Type != 0 {
						if t.Type&tag.SEARCH != 0 &&
							t.Data[tag.SEARCH] != nil &&
							t.Data[tag.SEARCH]["regex"] == true {
							if matched, _ := regexp.MatchString(t.Word, st.Word); matched {
								valid[i] = true
								continue STAG
							}
						} else {
							if t.Word == st.Word {
								valid[i] = true
								continue STAG
							}
						}
					}
				}
			}
		}
		matched[fid.String()] = func() bool {
			for _, v := range valid {
				if !v {
					return false
				}
			}
			return true
		}()
	}
	for id, match := range matched {
		if match {
			fid, _ := types.DecodeFileID(id)
			out = append(out, fid)
		}
	}
	return
}
