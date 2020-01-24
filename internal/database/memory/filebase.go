package memory

import (
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

type Filebase struct {
	Database
}

func (fb *Filebase) Reserve(id filehash.FileID) (filehash.FileID, error) {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	for _, ok := fb.Files[id.String()]; ok; _, ok = fb.Files[id.String()] {
		id = id.Mutate()
	}
	fb.Files[id.String()] = nil
	return id, nil
}

func (fb *Filebase) Insert(r database.FileI) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	if expectnil, ok := fb.Files[r.GetID().String()]; !ok {
		return database.ErrIDNotReserved
	} else if expectnil != nil {
		return database.ErrNameTaken
	}
	fb.Files[r.GetID().String()] = r
	return nil
}

// Get(fid filehash.FileID) (FileI, error)
// GetAll(fids ...filehash.FileID) ([]FileI, error)
// Update(r FileI) error
// Remove(r filehash.FileID) error
// GetOwned(uid OwnerID) ([]FileI, error)
// GetPermKey(uid OwnerID, pkey string) ([]FileI, error) // does not include owned records
// MatchStore(OwnerID, []filehash.StoreID, ...string) ([]FileI, error)
