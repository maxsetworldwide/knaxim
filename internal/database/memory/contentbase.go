package memory

type Contentbase struct {
	Database
}

// Insert(...ContentLine) error
// Len(id filehash.StoreID) (int64, error)
// Slice(id filehash.StoreID, start int, end int) ([]ContentLine, error)
// RegexSearchFile(regex string, file filehash.StoreID, start int, end int) ([]ContentLine, error)
