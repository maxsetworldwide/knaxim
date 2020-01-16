package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
)

type FileInfo struct {
	File  database.FileI `json:"file"`
	Count int64          `json:"count,omitempty"` //sentence count
	Size  int64          `json:"size,omitempty"`  //size of original file in bytes
}

type SearchResponse struct {
	Files []FileInfo `json:"matched"`
}

func BuildSearchResponse(r *http.Request, fids []filehash.FileID) SearchResponse {
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	var files []database.FileI
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
			lengths[fid.String()], err = r.Context().Value(database.CONTENT).(database.Contentbase).Len(fid.StoreID)
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

type DirInformation struct {
	Name  string            `json:"name"`
	Files []filehash.FileID `json:"files"`
}

type FileContent struct {
	Length int                    `json:"size"`
	Vals   []database.ContentLine `json:"lines"`
}

type GroupInformation struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members,omitempty"`
}

func BuildGroupInfo(grp database.GroupI) GroupInformation {
	return GroupInformation{
		ID:   grp.GetID().String(),
		Name: grp.GetName(),
		Members: func() []string {
			var out []string
			for _, mem := range grp.GetMembers() {
				out = append(out, mem.GetID().String())
			}
			return out
		}(),
	}
}

type UserInfo struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
	Data  struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"data,omitempty"`
}

func BuildUserInfo(r *http.Request, u database.UserI) UserInfo {
	actor := r.Context().Value(USER).(database.UserI)
	var ui UserInfo
	ui.ID = u.GetID().String()
	ui.Name = u.GetName()
	if actor.Equal(u) {
		ui.Roles = u.GetRoles()
		userbase := r.Context().Value(database.OWNER).(database.Ownerbase)
		var err error
		if ui.Data.Current, err = userbase.GetSpace(u.GetID()); err != nil {
			panic(err)
		}
		if ui.Data.Total, err = userbase.GetTotalSpace(u.GetID()); err != nil {
			panic(err)
		}
	}
	return ui
}
