package memory

type Filebase struct {
	Database
}

// Reserve(id filehash.FileID) (filehash.FileID, error)
// Insert(r FileI) error
// Get(fid filehash.FileID) (FileI, error)
// GetAll(fids ...filehash.FileID) ([]FileI, error)
// Update(r FileI) error
// Remove(r filehash.FileID) error
// GetOwned(uid OwnerID) ([]FileI, error)
// GetPermKey(uid OwnerID, pkey string) ([]FileI, error) // does not include owned records
// MatchStore(OwnerID, []filehash.StoreID, ...string) ([]FileI, error)
