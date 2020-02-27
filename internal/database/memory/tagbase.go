package memory

import (
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
)

// Tagbase wraps database and provides tag operations
type Tagbase struct {
	Database
}

// UpsertFile adds tags attached to fileid
func (tb *Tagbase) UpsertFile(fid filehash.FileID, tags ...tag.Tag) error {
	lock.Lock()
	defer lock.Unlock()
	if tb.TagFiles[fid.String()] == nil {
		tb.TagFiles[fid.String()] = make(map[string]tag.Tag)
	}
	for _, t := range tags {
		if tb.TagFiles[fid.String()][t.Word].Word != t.Word {
			tb.TagFiles[fid.String()][t.Word] = t
		} else {
			tb.TagFiles[fid.String()][t.Word] = tb.TagFiles[fid.String()][t.Word].Update(t)
		}
	}
	return nil
}

// UpsertStore add tags attached to storeids
func (tb *Tagbase) UpsertStore(sid filehash.StoreID, tags ...tag.Tag) error {
	lock.Lock()
	defer lock.Unlock()
	if tb.TagStores[sid.String()] == nil {
		tb.TagStores[sid.String()] = make(map[string]tag.Tag)
	}
	for _, t := range tags {
		if tb.TagStores[sid.String()][t.Word].Word != t.Word {
			tb.TagStores[sid.String()][t.Word] = t
		} else {
			tb.TagStores[sid.String()][t.Word] = tb.TagStores[sid.String()][t.Word].Update(t)
		}
	}
	return nil
}

// FileTags returns all tags associated with a particular fileid
func (tb *Tagbase) FileTags(fids ...filehash.FileID) (map[string][]tag.Tag, error) {
	lock.RLock()
	defer lock.RUnlock()
	storeids := make([]filehash.StoreID, 0, len(fids))
	for _, fid := range fids {
		storeids = append(storeids, fid.StoreID)
	}
	var perr error
	{
		sb := tb.store(nil).(*Storebase)
		for _, sid := range storeids {
			fs, err := sb.get(sid)
			if err != nil {
				sb.close()
				return nil, err
			}
			if fs.Perr != nil {
				perr = fs.Perr
				break
			}
		}
		sb.close()
	}
	out := make(map[string][]tag.Tag)
	for _, fid := range fids {
		for w, tag := range tb.TagFiles[fid.String()] {
			out[w] = append(out[w], tag)
		}
	}
	for _, sid := range storeids {
		for w, tag := range tb.TagStores[sid.String()] {
			out[w] = append(out[w], tag)
		}
	}
	return out, perr
}

// GetFiles returns all fileids and storeids associated with particular
// tags, optionally allows only searching over certain FileIDs
func (tb *Tagbase) GetFiles(filters []tag.Tag, context ...filehash.FileID) (fileids []filehash.FileID, storeids []filehash.StoreID, err error) {
	lock.RLock()
	defer lock.RUnlock()
	if len(context) == 0 {
	STORES:
		for sidstr, tags := range tb.TagStores {
			for _, filter := range filters {
				tag, assigned := tags[filter.Word]
				if !assigned {
					continue STORES
				}
				if tag.Type&filter.Type == 0 {
					continue STORES
				}
				if !tag.Data.Contains(filter.Data) {
					continue STORES
				}
			}
			sid, _ := filehash.DecodeStoreID(sidstr)
			storeids = append(storeids, sid)
		}
		for fidstr := range tb.TagFiles {
			fid, _ := filehash.DecodeFileID(fidstr)
			context = append(context, fid)
		}
	}
FILES:
	for _, fid := range context {
		for _, filter := range filters {
			var filetag, storetag tag.Tag
			var fassigned, sassigned bool
			if tb.TagFiles[fid.String()] != nil {
				filetag, fassigned = tb.TagFiles[fid.String()][filter.Word]
			}
			if tb.TagStores[fid.StoreID.String()] != nil {
				storetag, sassigned = tb.TagStores[fid.StoreID.String()][filter.Word]
			}
			if !fassigned && !sassigned {
				continue FILES
			}
			if (filetag.Type|storetag.Type)&filter.Type == 0 {
				continue FILES
			}
			for typ, info := range filter.Data {
				for k, v := range info {
					if (filetag.Data[typ] == nil || filetag.Data[typ][k] != v) && (storetag.Data[typ] == nil || storetag.Data[typ][k] != v) {
						continue FILES
					}
				}
			}
		}
		fileids = append(fileids, fid)
		storeids = append(storeids, fid.StoreID)
	}
	return
}

// SearchData returns all tags that have matching data fields
func (tb *Tagbase) SearchData(typ tag.Type, d tag.Data) (out []tag.Tag, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, filetags := range tb.TagFiles {
		for _, tag := range filetags {
			if tag.Type == typ && tag.Data.Contains(d) {
				out = append(out, tag)
			}
		}
	}
	for _, storetags := range tb.TagStores {
		for _, tag := range storetags {
			if tag.Type == typ && tag.Data.Contains(d) {
				out = append(out, tag)
			}
		}
	}
	return
}
