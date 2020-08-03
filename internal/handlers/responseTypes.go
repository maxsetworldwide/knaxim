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

package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/util"
)

// FileInfo json response type of file information
type FileInfo struct {
	File  types.FileI `json:"file"`
	Count int64       `json:"count,omitempty"` //sentence count
	Size  int64       `json:"size,omitempty"`  //size of original file in bytes
}

// SearchResponse json response of matched files to a search
type SearchResponse struct {
	Files []FileInfo `json:"matched"`
}

// BuildSearchResponse contructs SearchResponse from a list of matched fileids
func BuildSearchResponse(r *http.Request, fids []types.FileID) SearchResponse {
	filebase := r.Context().Value(types.FILE).(database.Filebase)
	var files []types.FileI
	var lengths map[string]int64
	errch := make(chan error, 2)
	go func() {
		var err error
		files, err = filebase.GetAll(fids...)
		errch <- err
	}()
	go func() {
		var err error
		lengths = make(map[string]int64)
		for _, fid := range fids {
			lengths[fid.String()], err = r.Context().Value(types.CONTENT).(database.Contentbase).Len(fid.StoreID)
			if err != nil {
				errch <- err
				return
			}
		}
		errch <- nil
	}()
	//get 2 errors from routines
	err := <-errch
	if err != nil {
		panic(err)
	}
	err = <-errch
	if err != nil {
		panic(err)
	}
	result := SearchResponse{
		Files: make([]FileInfo, 0, len(files)),
	}
	for _, file := range files {
		result.Files = append(result.Files, FileInfo{file, lengths[file.GetID().String()], 0})
	}
	return result
}

// DirInformation is the json encoding for folder information
type DirInformation struct {
	Name  string         `json:"name"`
	Files []types.FileID `json:"files"`
}

// FileContent is the json encoding for lines from a file
type FileContent struct {
	Length int                 `json:"size"`
	Vals   []types.ContentLine `json:"lines"`
}

// GroupInformation is the json encoding object for Groups
type GroupInformation struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Members []string `json:"members,omitempty"`
}

// BuildGroupInfo contructs Group Response Object from Group
func BuildGroupInfo(grp types.GroupI) GroupInformation {
	return GroupInformation{
		ID:    grp.GetID().String(),
		Name:  grp.GetName(),
		Owner: grp.GetOwner().GetID().String(),
		Members: func() []string {
			var out []string
			for _, mem := range grp.GetMembers() {
				out = append(out, mem.GetID().String())
			}
			return out
		}(),
	}
}

// UserInfo is the json contruction struct for user information
type UserInfo struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
	Data  struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"data,omitempty"`
}

// BuildUserInfo builds user info from user object
func BuildUserInfo(r *http.Request, u types.UserI) UserInfo {
	actor := r.Context().Value(USER).(types.UserI)
	var ui UserInfo
	ui.ID = u.GetID().String()
	ui.Name = u.GetName()
	if actor.Equal(u) {
		ui.Roles = u.GetRoles()
		userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
		var err error
		if ui.Data.Current, err = userbase.GetSpace(u.GetID()); err != nil {
			util.VerboseRequest(r, "unable to get current files")
			panic(err)
		}
		if ui.Data.Total, err = userbase.GetTotalSpace(u.GetID()); err != nil {
			util.VerboseRequest(r, "unable to get total files")
			panic(err)
		}
	}
	return ui
}
