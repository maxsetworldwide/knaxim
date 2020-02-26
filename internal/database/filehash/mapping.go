package filehash

// FileMap is a mapping of StoreIDs to associated FileIDs
type FileMap map[StoreID][]FileID

// GetStoreID returns associated StoreID of a FileID and saves the associaion in the map
func (fm FileMap) GetStoreID(fid FileID) StoreID {
	fm[fid.StoreID] = append(fm[fid.StoreID], fid)
	return fid.StoreID
}

// GetFileID returns all saved associated FileIDs to that StoreID
func (fm FileMap) GetFileID(sid StoreID) []FileID {
	return fm[sid]
}

// StoreIDs converts FileIDs to StoreIDs while saving the association
func (fm FileMap) StoreIDs(fids ...FileID) []StoreID {
	out := make([]StoreID, 0)
	for _, f := range fids {
		out = append(out, fm.GetStoreID(f))
	}
	return out
}

// FileIDs returns the list of FileIDs associated with any of the storeids in the map
func (fm FileMap) FileIDs(sids ...StoreID) []FileID {
	out := make([]FileID, 0)
	for _, s := range sids {
		out = append(out, fm.GetFileID(s)...)
	}
	return out
}
