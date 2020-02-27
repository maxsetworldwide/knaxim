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
