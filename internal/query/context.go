// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"context"
	"errors"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
)

// CType determines how the context is generated. More specifically, it determines how the id value is to be interpreted,
type CType uint8

const (
	// OWNER means that the id is an OwnerID
	OWNER CType = iota
	// FILE means that the id is a FileID
	FILE
)

func decodeCType(s string) (CType, error) {
	switch s {
	case "o":
		fallthrough
	case "owner":
		return OWNER, nil
	case "f":
		fallthrough
	case "file":
		return FILE, nil
	default:
		return 0, errors.New("unrecognized Context Type")
	}
}

// CRestriction gives greater specificity to context types
type CRestriction uint8

const (
	// ALL means both owned and viewable files appear in the context
	ALL CRestriction = 3
	// OWNED means only owned files appear in the context
	OWNED CRestriction = 1
	// VIEW means only viewable file appear in the context
	VIEW CRestriction = 2
)

func decodeCRestriction(s string) (CRestriction, error) {
	switch s {
	case "":
		fallthrough
	case "all":
		return ALL, nil
	case "o":
		fallthrough
	case "owned":
		return OWNED, nil
	case "viewable":
		fallthrough
	case "v":
		return VIEW, nil
	default:
		return 0, errors.New("Unrecognized Context Restriction")
	}
}

// C is the context that the search is performed over. Data is used to identify a list of files
type C struct {
	Type  CType        `json:"type"`
	ID    string       `json:"id"`
	Limit CRestriction `json:"only,omitempty"`
}

func decodeC(i interface{}) (contexts []C, err error) {
	switch v := i.(type) {
	case []interface{}:
		for _, ele := range v {
			var temp []C
			temp, err = decodeC(ele)
			if err != nil {
				return
			}
			contexts = append(contexts, temp...)
		}
	case map[string]interface{}:
		tstr, ok := v["type"].(string)
		if !ok {
			return nil, errors.New("Missing Context Type")
		}
		var t CType
		t, err = decodeCType(tstr)
		if err != nil {
			return
		}
		id, ok := v["id"].(string)
		if !ok {
			return nil, errors.New("Missing ID of context")
		}
		restriction := ALL
		if i, assigned := v["only"]; assigned {
			if r, ok := i.(string); ok {
				restriction, err = decodeCRestriction(r)
			}
		}
		contexts = append(contexts, C{
			Type:  t,
			ID:    id,
			Limit: restriction,
		})
	case string:
		contexts = append(contexts, C{
			Type:  OWNER,
			ID:    v,
			Limit: ALL,
		})
	default:
		return nil, errors.New("unrecognized Context Value")
	}
	return
}

// GetFileSet returns the list of fileids that the context maps to
func (c C) GetFileSet(ctx context.Context, dbConfig database.Database) ([]types.FileID, error) {
	db, err := dbConfig.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)
	return c.getFileSet(db)
}

func (c C) getFileSet(db database.Database) ([]types.FileID, error) {
	switch c.Type {
	case OWNER:
		id, err := types.DecodeOwnerIDString(c.ID)
		if err != nil {
			return nil, err
		}
		var list []types.FileID
		if c.Limit&OWNED != 0 {
			owned, err := db.File().GetOwned(id)
			if err != nil {
				return nil, err
			}
			for _, file := range owned {
				list = append(list, file.GetID())
			}
		}
		if c.Limit&VIEW != 0 {
			viewable, err := db.File().GetPermKey(id, "view")
			if err != nil {
				return nil, err
			}
			for _, file := range viewable {
				list = append(list, file.GetID())
			}
		}
		return list, nil
	case FILE:
		id, err := types.DecodeFileID(c.ID)
		if err != nil {
			return nil, err
		}
		return []types.FileID{id}, nil
	default:
		return nil, errors.New("Unrecognized Context Type")
	}
}

// CheckAccess returns true if the provided owner has permission to access the files contexts, extra provides additional permissions to check on file type contexts
func (c C) CheckAccess(o types.Owner, dbConnection database.Database, extra ...string) (bool, error) {
	switch c.Type {
	case OWNER:
		oid, err := types.DecodeOwnerIDString(c.ID)
		if err != nil {
			return false, err
		}
		O, err := dbConnection.Owner().Get(oid)
		if err != nil {
			return false, err
		}
		return O.Match(o), nil
	case FILE:
		fid, err := types.DecodeFileID(c.ID)
		if err != nil {
			return false, err
		}
		file, err := dbConnection.File().Get(fid)
		if err != nil {
			return false, err
		}
		access := file.GetOwner().Match(o)
		for _, ex := range extra {
			access = access || file.CheckPerm(o, ex)
		}
		return access, nil
	default:
		return false, errors.New("unrecognized context type")
	}
}
