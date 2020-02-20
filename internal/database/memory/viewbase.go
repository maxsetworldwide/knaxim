package memory

import (
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

type Viewbase struct {
	Database
}

func (vb *Viewbase) Insert(vs *database.ViewStore) error {
	lock.Lock()
	defer lock.Unlock()
	vb.Views[vs.ID.String()] = vs
	return nil
}

func (vb *Viewbase) Get(id filehash.StoreID) (out *database.ViewStore, err error) {
	lock.RLock()
	defer lock.RUnlock()
	out, ok := vb.Views[id.String()]
	if !ok {
		return nil, database.ErrNotFound
	}
	return
}
