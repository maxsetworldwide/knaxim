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

func (fb *Filebase) Get(fid filehash.FileID) (database.FileI, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()
	if fb.Files[fid.String()] == nil {
		return nil, database.ErrNotFound
	}
	return fb.Files[fid.String()].Copy(), nil
}

func (fb *Filebase) GetAll(fids ...filehash.FileID) ([]database.FileI, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()
	out := make([]database.FileI, 0, len(fids))
	for _, fid := range fids {
		temp, err := fb.Get(fid)
		if err != nil {
			return nil, err
		}
		out = append(out, temp)
	}
	return out, nil
}

func (fb *Filebase) Update(r database.FileI) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	if fb.Files[r.GetID().String()] == nil {
		return database.ErrNotFound
	}
	fb.Files[r.GetID().String()] = r.Copy()
	return nil
}

func (fb *Filebase) Remove(r filehash.FileID) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	if fb.Files[r.String()] == nil {
		return database.ErrNotFound
	}
	delete(fb.Files, r.String())
	return nil
}

func (fb *Filebase) GetOwned(uid database.OwnerID) ([]database.FileI, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()
	var out []database.FileI
	for _, file := range fb.Files {
		if file.GetOwner().GetID().Equal(uid) {
			out = append(out, file.Copy())
		}
	}
	return out, nil
}

func (fb *Filebase) GetPermKey(uid database.OwnerID, pkey string) ([]database.FileI, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()
	var out []database.FileI
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

func (fb *Filebase) MatchStore(oid database.OwnerID, sids []filehash.StoreID, pkeys ...string) ([]database.FileI, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()
	var out []database.FileI
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
