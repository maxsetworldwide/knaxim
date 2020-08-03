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
