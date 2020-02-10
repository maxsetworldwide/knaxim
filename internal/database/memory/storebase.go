package memory

import (
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

type Storebase struct {
	Database
}

func (sb *Storebase) Reserve(id filehash.StoreID) (filehash.StoreID, error) {
	lock.Lock()
	defer lock.Unlock()
	for _, assigned := sb.Stores[id.String()]; assigned; _, assigned = sb.Stores[id.String()] {
		id = id.Mutate()
	}
	sb.Stores[id.String()] = nil
	return id, nil
}

func (sb *Storebase) Insert(fs *database.FileStore) error {
	lock.Lock()
	defer lock.Unlock()
	if expectnil, assigned := sb.Stores[fs.ID.String()]; !assigned {
		return database.ErrIDNotReserved
	} else if expectnil != nil {
		return database.ErrNameTaken
	}
	sb.Stores[fs.ID.String()] = fs
	return nil
}

func (sb *Storebase) Get(id filehash.StoreID) (*database.FileStore, error) {
	lock.RLock()
	defer lock.RUnlock()
	return sb.get(id)
}

func (sb *Storebase) get(id filehash.StoreID) (*database.FileStore, error) {
	if sb.Stores[id.String()] == nil {
		return nil, database.ErrNotFound
	}
	return sb.Stores[id.String()].Copy(), nil
}

func (sb *Storebase) MatchHash(h uint32) (out []*database.FileStore, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, store := range sb.Stores {
		if store.ID.Hash == h {
			out = append(out, store)
		}
	}
	return
}

func (sb *Storebase) UpdateMeta(fs *database.FileStore) error {
	lock.Lock()
	defer lock.Unlock()
	if sb.Stores[fs.ID.String()] == nil {
		return database.ErrNotFound
	}
	sb.Stores[fs.ID.String()].ContentType = fs.ContentType
	sb.Stores[fs.ID.String()].FileSize = fs.FileSize
	sb.Stores[fs.ID.String()].Perr = fs.Perr
	return nil
}
