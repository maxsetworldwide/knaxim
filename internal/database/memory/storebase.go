// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memory

import (
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
)

// Storebase wraps database for file store actions
type Storebase struct {
	Database
}

// Reserve is the first step in adding a new file store. returns
// reserved StoreID, might have been mutated from input
func (sb *Storebase) Reserve(id types.StoreID) (types.StoreID, error) {
	lock.Lock()
	defer lock.Unlock()
	for _, assigned := sb.Stores[id.String()]; assigned; _, assigned = sb.Stores[id.String()] {
		id = id.Mutate()
	}
	sb.Stores[id.String()] = nil
	return id, nil
}

// Insert adds new filestore to database
func (sb *Storebase) Insert(fs *types.FileStore) error {
	lock.Lock()
	defer lock.Unlock()
	if expectnil, assigned := sb.Stores[fs.ID.String()]; !assigned {
		return errors.ErrIDNotReserved
	} else if expectnil != nil {
		return errors.ErrNameTaken
	}
	sb.Stores[fs.ID.String()] = fs
	return nil
}

// Get File Store
func (sb *Storebase) Get(id types.StoreID) (*types.FileStore, error) {
	lock.RLock()
	defer lock.RUnlock()
	return sb.get(id)
}

func (sb *Storebase) get(id types.StoreID) (*types.FileStore, error) {
	if sb.Stores[id.String()] == nil {
		return nil, errors.ErrNotFound
	}
	return sb.Stores[id.String()].Copy(), nil
}

// MatchHash returns all filestores that have a particular hash
func (sb *Storebase) MatchHash(h uint32) (out []*types.FileStore, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, store := range sb.Stores {
		if store.ID.Hash == h {
			out = append(out, store)
		}
	}
	return
}

// UpdateMeta update meta data values of a filestore
func (sb *Storebase) UpdateMeta(fs *types.FileStore) error {
	lock.Lock()
	defer lock.Unlock()
	if sb.Stores[fs.ID.String()] == nil {
		return errors.ErrNotFound
	}
	sb.Stores[fs.ID.String()].ContentType = fs.ContentType
	sb.Stores[fs.ID.String()].FileSize = fs.FileSize
	sb.Stores[fs.ID.String()].Perr = fs.Perr
	return nil
}
