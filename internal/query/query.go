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
	"encoding/json"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

//Q is the primary query type that represents the combination of Context and Matching condition
type Q struct {
	Context []C `json:"context"`
	Match   []M `json:"match"`
}

// UnmarshalJSON reads json into Query object
func (q *Q) UnmarshalJSON(b []byte) error {
	var target struct {
		C interface{} `json:"context"`
		M interface{} `json:"match"`
	}
	err := json.Unmarshal(b, &target)
	if err != nil {
		return err
	}
	if q.Context, err = decodeC(target.C); err != nil {
		return err
	}
	if q.Match, err = decodeM(target.M); err != nil {
		return err
	}
	return nil
}

// FindMatching finds all matching fileids based on query
func (q *Q) FindMatching(ctx context.Context, dbConfig database.Database) (files []types.FileID, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	db, err := dbConfig.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)

	filelistch := make(chan []types.FileID, len(q.Context))
	errch := make(chan error)
	fileWG := new(sync.WaitGroup)
	fileWG.Add(len(q.Context))
	for _, c := range q.Context {
		go func(c C) {
			defer fileWG.Done()
			subset, err := c.getFileSet(db)
			filelistch <- subset
			select {
			case errch <- err:
			case <-ctx.Done():
			}
		}(c)
	}
	go func() {
		fileWG.Wait()
		close(filelistch)
	}()

	fullListCh := make(chan []types.FileID)
	go func() {
		fileset := make(map[string]bool)
		var fullList []types.FileID
		for list := range filelistch {
			for _, fid := range list {
				if fstr := fid.String(); !fileset[fstr] {
					fullList = append(fullList, fid)
					fileset[fstr] = true
				}
			}
		}
		select {
		case fullListCh <- fullList:
		case <-ctx.Done():
		}
	}()

	for i := 0; i < len(q.Context); i++ {
		if e := <-errch; e != nil {
			return nil, e
		}
	}
	filelist := <-fullListCh
	var matchTags []tag.FileTag
	for _, m := range q.Match {
		matchTags = append(matchTags, m.SearchTag())
	}
	return db.Tag().SearchFiles(filelist, matchTags...)
}
