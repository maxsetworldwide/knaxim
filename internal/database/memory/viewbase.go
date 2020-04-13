package memory

import (
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

// Viewbase wraps database and provides view operations
type Viewbase struct {
	Database
}

// Insert adds new viewstore to the database
func (vb *Viewbase) Insert(vs *types.ViewStore) error {
	lock.Lock()
	defer lock.Unlock()
	vb.Views[vs.ID.String()] = vs
	return nil
}

// Get viewstore of associated id
func (vb *Viewbase) Get(id types.StoreID) (out *types.ViewStore, err error) {
	lock.RLock()
	defer lock.RUnlock()
	out, ok := vb.Views[id.String()]
	if !ok {
		return nil, errors.ErrNotFound
	}
	return
}
