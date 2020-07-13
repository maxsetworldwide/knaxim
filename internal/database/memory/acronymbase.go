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

// Acronymbase is a memory Database accessor of acronym operators
type Acronymbase struct {
	Database
}

// Put adds association of acronym and phrase into database
func (ab *Acronymbase) Put(acronym string, phrase string) error {
	lock.Lock()
	defer lock.Unlock()
	ab.Acronyms[acronym] = append(ab.Acronyms[acronym], phrase)
	return nil
}

// Get get a list of all associated phrases of a particular acronym8
func (ab *Acronymbase) Get(acronym string) ([]string, error) {
	lock.RLock()
	defer lock.RUnlock()
	defCopy := make([]string, len(ab.Acronyms[acronym]))
	copy(defCopy, ab.Acronyms[acronym])
	return defCopy, nil
}
