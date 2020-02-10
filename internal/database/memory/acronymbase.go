package memory

type Acronymbase struct {
	Database
}

func (ab *Acronymbase) Put(acronym string, phrase string) error {
	lock.Lock()
	defer lock.Unlock()
	ab.Acronyms[acronym] = append(ab.Acronyms[acronym], phrase)
	return nil
}

func (ab *Acronymbase) Get(acronym string) ([]string, error) {
	lock.RLock()
	defer lock.RUnlock()
	defCopy := make([]string, len(ab.Acronyms[acronym]))
	copy(defCopy, ab.Acronyms[acronym])
	return defCopy, nil
}
