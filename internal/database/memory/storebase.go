package memory

type Storebase struct {
	Database
}

// Reserve(id filehash.StoreID) (filehash.StoreID, error)
// Insert(fs *FileStore) error
// Get(id filehash.StoreID) (*FileStore, error)
// MatchHash(h uint32) ([]*FileStore, error)
