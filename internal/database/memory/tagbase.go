package memory

type Tagbase struct {
	Database
}

// UpsertFile(filehash.FileID, ...tag.Tag) error
// UpsertStore(filehash.StoreID, ...tag.Tag) error
// FileTags(...filehash.FileID) (map[string][]tag.Tag, error)
// GetFiles([]tag.Tag, ...filehash.FileID) ([]filehash.FileID, []filehash.StoreID, error)
// SearchData(tag.Type, tag.Data) ([]tag.Tag, error)
