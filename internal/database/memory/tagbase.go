package memory

import (
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
)

type Tagbase struct {
	Database
}

func (tb *Tagbase) UpsertFile(fid filehash.FileID, tags ...tag.Tag) error {
	tb.lock.Lock()
	defer tb.lock.Unlock()
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

func (tb *Tagbase) UpsertStore(sid filehash.StoreID, tags ...tag.Tag) error {
	tb.lock.Lock()
	defer tb.lock.Unlock()
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

func (tb *Tagbase) FileTags(fids ...filehash.FileID) (map[string][]tag.Tag, error) {
	tb.lock.RLock()
	defer tb.lock.RUnlock()
	storeids := make([]filehash.StoreID, 0, len(fids))
	for _, fid := range fids {
		storeids = append(storeids, fid.StoreID)
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
	return out, nil
}

// GetFiles([]tag.Tag, ...filehash.FileID) ([]filehash.FileID, []filehash.StoreID, error)

func (tb *Tagbase) SearchData(tag.Type, tag.Data) ([]tag.Tag, error) {

}
