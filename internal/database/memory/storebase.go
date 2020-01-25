package memory

import (
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

type Storebase struct {
	Database
}

func (sb *Storebase) Reserve(id filehash.StoreID) (filehash.StoreID, error) {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	for _, assigned := sb.Stores[id.String()]; assigned; _, assigned = sb.Stores[id.String()] {
		id = id.Mutate()
	}
	sb.Stores[id.String()] = nil
	return id, nil
}

func (sb *Storebase) Insert(fs *database.FileStore) error {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	if expectnil, assigned := sb.Stores[fs.ID.String()]; !assigned {
		return database.ErrIDNotReserved
	} else if expectnil != nil {
		return database.ErrNameTaken
	}
	sb.Stores[fs.ID.String()] = fs
	return nil
}

func (sb *Storebase) Get(id filehash.StoreID) (*database.FileStore, error) {
	sb.lock.RLock()
	defer sb.lock.RUnlock()
	if sb.Stores[id.String()] == nil {
		return nil, database.ErrNotFound
	}
	return sb.Stores[id.String()].Copy(), nil
}

func (sb *Storebase) MatchHash(h uint32) (out []*database.FileStore, err error) {
	sb.lock.RLock()
	defer sb.lock.RUnlock()
	for _, store := range sb.Stores {
		if store.ID.Hash == h {
			out = append(out, store)
		}
	}
	return
}
