package memory

import (
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

// Filebase is the memory database accessor for file operations
type Filebase struct {
	Database
}

// Reserve is the first step in inserting a new file and it reserves a FileID,
// mutating it if necessary. and returns the FileID that has been reserved
func (fb *Filebase) Reserve(id types.FileID) (types.FileID, error) {
	lock.Lock()
	defer lock.Unlock()
	for _, ok := fb.Files[id.String()]; ok; _, ok = fb.Files[id.String()] {
		id = id.Mutate()
	}
	fb.Files[id.String()] = nil
	return id, nil
}

// Insert addes file to datase, file's fileid must be already reserved
func (fb *Filebase) Insert(r types.FileI) error {
	lock.Lock()
	defer lock.Unlock()
	if expectnil, ok := fb.Files[r.GetID().String()]; !ok {
		return errors.ErrIDNotReserved
	} else if expectnil != nil {
		return errors.ErrNameTaken
	}
	fb.Files[r.GetID().String()] = r
	return nil
}

// Get returns file matching id
func (fb *Filebase) Get(fid types.FileID) (types.FileI, error) {
	lock.RLock()
	defer lock.RUnlock()
	return fb.get(fid)
}

func (fb *Filebase) get(fid types.FileID) (types.FileI, error) {
	if fb.Files[fid.String()] == nil {
		return nil, errors.ErrNotFound
	}
	return fb.Files[fid.String()].Copy(), nil
}

// GetAll returns all files matching file ids
func (fb *Filebase) GetAll(fids ...types.FileID) ([]types.FileI, error) {
	lock.RLock()
	defer lock.RUnlock()
	out := make([]types.FileI, 0, len(fids))
	for _, fid := range fids {
		temp, err := fb.get(fid)
		if err != nil {
			return nil, err
		}
		out = append(out, temp)
	}
	return out, nil
}

// Update replaces file matching fileid
func (fb *Filebase) Update(r types.FileI) error {
	lock.Lock()
	defer lock.Unlock()
	if fb.Files[r.GetID().String()] == nil {
		return errors.ErrNotFound
	}
	fb.Files[r.GetID().String()] = r.Copy()
	return nil
}

// Remove file from database
func (fb *Filebase) Remove(r types.FileID) error {
	lock.Lock()
	defer lock.Unlock()
	if fb.Files[r.String()] == nil {
		return errors.ErrNotFound
	}
	delete(fb.Files, r.String())
	return nil
}

// GetOwned returns all files owned by ownerid
func (fb *Filebase) GetOwned(uid types.OwnerID) ([]types.FileI, error) {
	lock.RLock()
	defer lock.RUnlock()
	return fb.getOwned(uid)
}

func (fb *Filebase) getOwned(uid types.OwnerID) ([]types.FileI, error) {
	var out []types.FileI
	for _, file := range fb.Files {
		if file.GetOwner().GetID().Equal(uid) {
			out = append(out, file.Copy())
		}
	}
	return out, nil
}

// GetPermKey returns all files that a given owner has a particular permission
func (fb *Filebase) GetPermKey(uid types.OwnerID, pkey string) ([]types.FileI, error) {
	lock.RLock()
	defer lock.RUnlock()
	return fb.getPermKey(uid, pkey)
}

func (fb *Filebase) getPermKey(uid types.OwnerID, pkey string) ([]types.FileI, error) {
	var out []types.FileI
LOOP:
	for _, file := range fb.Files {
		for _, o := range file.GetPerm(pkey) {
			if o.GetID().Equal(uid) {
				out = append(out, file.Copy())
				continue LOOP
			}
		}
	}
	return out, nil
}

// MatchStore returns all files that match one of the storeids,
// and is either owned by oid or oid has one of the form of permission
func (fb *Filebase) MatchStore(oid types.OwnerID, sids []types.StoreID, pkeys ...string) ([]types.FileI, error) {
	lock.RLock()
	defer lock.RUnlock()
	var out []types.FileI
	for _, file := range fb.Files {
		if func() bool {
			for _, sid := range sids {
				if sid.Equal(file.GetID().StoreID) {
					return true
				}
			}
			return false
		}() && (file.GetOwner().GetID().Equal(oid) ||
			func() bool {
				for _, pkey := range pkeys {
					for _, o := range file.GetPerm(pkey) {
						if o.GetID().Equal(oid) {
							return true
						}
					}
				}
				return false
			}()) {
			out = append(out, file.Copy())
		}
	}
	return out, nil
}
