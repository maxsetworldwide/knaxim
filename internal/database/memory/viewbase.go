/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
