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

package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

func TestTag(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	tb := DB.Tag()
	defer DB.Close(nil)
	t.Parallel()
	fileids := []types.FileID{
		types.FileID{
			StoreID: types.StoreID{
				Hash:  123456,
				Stamp: 789,
			},
			Stamp: []byte("abctag"),
		},
		types.FileID{
			StoreID: types.StoreID{
				Hash:  4568196,
				Stamp: 4897,
			},
			Stamp: []byte("tagtest"),
		},
	}
	ownerids := []types.OwnerID{
		types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'e', 's', 't'},
			Stamp:       []byte{'1', 't', 'a', 'g'},
		},
		types.OwnerID{
			Type:        't',
			UserDefined: [3]byte{'e', 's', 't'},
			Stamp:       []byte{'2'},
		},
	}
	tags := []tag.FileTag{
		tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "test",
				Type: tag.USER | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[1],
			Owner: ownerids[0],
			Tag: tag.Tag{
				Word: "test2",
				Type: tag.TOPIC | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[0],
			Owner: ownerids[1],
			Tag: tag.Tag{
				Word: "test3",
				Type: tag.ACTION | tag.CONTENT,
			},
		},
		tag.FileTag{
			File:  fileids[1],
			Owner: ownerids[1],
			Tag: tag.Tag{
				Word: "test4",
				Type: tag.ALLFILE,
			},
		},
	}
	t.Logf("fileids: %v", fileids)
	t.Logf("ownerids: %v", ownerids)
	t.Log("Upsert")
	err := tb.Upsert(tags...)
	if err != nil {
		t.Fatalf("Unable to Upsert tags: %s", err.Error())
	}

	t.Log("Get")
	{
		gtags, err := tb.Get(fileids[0], ownerids[0])
		if err != nil {
			t.Fatalf("Unable to Get tags: %s", err.Error())
		}
		if len(gtags) != 2 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("GetType")
	{
		gtags, err := tb.GetType(fileids[0], ownerids[0], tag.USER)
		if err != nil {
			t.Fatalf("Unable to GetType: %s", err.Error())
		}
		if len(gtags) != 1 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("GetAll")
	{
		gtags, err := tb.GetAll(tag.USER, ownerids[0])
		if err != nil {
			t.Fatalf("Unable to GetAll: %s", err.Error())
		}
		if len(gtags) != 1 {
			t.Fatalf("Incorrect Return: %v", gtags)
		}
	}

	t.Log("SearchFiles")
	{
		stags, err := tb.SearchFiles(fileids, tag.FileTag{
			Tag: tag.Tag{
				Word: "test",
				Type: tag.ALLTYPES,
				Data: tag.Data{
					tag.SEARCH: map[string]interface{}{
						"regex": true,
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("Unable to Search Files: %s", err.Error())
		}
		if len(stags) != 2 {
			t.Fatalf("Incorrect Return: %v", stags)
		}
	}
	owner := &types.User{
		ID:   ownerids[0],
		Name: "tagtestowner",
	}
	files := []*types.File{
		&types.File{
			ID: fileids[0],
			Permission: types.Permission{
				Own: owner,
			},
		},
		&types.File{
			ID: fileids[1],
			Permission: types.Permission{
				Own: owner,
			},
		},
	}
	{
		ob := tb.Owner()
		if _, err := ob.Reserve(owner.GetID(), owner.GetName()); err != nil {
			t.Fatalf("unable to reserve owner: %s", err.Error())
		}
		if err := ob.Insert(owner); err != nil {
			t.Fatalf("unable to insert owner: %s", err.Error())
		}
		fb := tb.File()
		for _, fid := range fileids {
			if _, err := fb.Reserve(fid); err != nil {
				t.Fatalf("unable to reserve fid: %s", fid.String())
			}
		}
		for i, f := range files {
			if err := fb.Insert(f); err != nil {
				t.Fatalf("unable to insert file %d", i)
			}
		}
	}
	t.Log("SearchOwned")
	{
		stags, err := tb.SearchOwned(ownerids[0], tag.FileTag{
			Tag: tag.Tag{
				Word: "test",
				Type: tag.ALLTYPES,
				Data: tag.Data{
					tag.SEARCH: map[string]interface{}{
						"regex": true,
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to search Owned: %s", err.Error())
		}
		if len(stags) != 2 {
			t.Fatalf("incorrect return: %v", stags)
		}
	}
}
