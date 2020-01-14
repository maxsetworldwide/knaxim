package filehash

type FileMap map[StoreID][]FileID

func (fm FileMap) GetStoreID(fid FileID) StoreID {
	fm[fid.StoreID] = append(fm[fid.StoreID], fid)
	return fid.StoreID
}

func (fm FileMap) GetFileID(sid StoreID) []FileID {
	return fm[sid]
}

func (fm FileMap) StoreIDs(fids ...FileID) []StoreID {
	out := make([]StoreID, 0)
	for _, f := range fids {
		out = append(out, fm.GetStoreID(f))
	}
	return out
}

func (fm FileMap) FileIDs(sids ...StoreID) []FileID {
	out := make([]FileID, 0)
	for _, s := range sids {
		out = append(out, fm.GetFileID(s)...)
	}
	return out
}
