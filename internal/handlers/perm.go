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
	"errors"
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"git.maxset.io/web/knaxim/internal/database/types"
	"github.com/gorilla/mux"
)

// AttachPerm is for paths related to permission actions
func AttachPerm(r *mux.Router) {
	// TODO: Move neccessary middleware use commands here and each attach??

	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.HandleFunc("/{type}/{id}/public", setPermissionPublic(true)).Methods("POST")
	r.HandleFunc("/{type}/{id}/public", setPermissionPublic(false)).Methods("DELETE")
	r.HandleFunc("/{type}/{id}", getPermissions).Methods("GET")
	r.HandleFunc("/{type}/{id}", setPermission(true)).Methods("POST")
	r.HandleFunc("/{type}/{id}", setPermission(false)).Methods("DELETE")
}

type permissionReport struct {
	Owner       string              `json:"owner"`
	IsOwner     bool                `json:"isOwned"`
	Permissions map[string][]string `json:"vals,omitempty"`
}

func buildPermissionReport(perm types.PermissionI, actor types.Owner) permissionReport {
	var out permissionReport
	out.Owner = perm.GetOwner().GetID().String()
	if perm.GetOwner().Match(actor) {
		out.IsOwner = true
		out.Permissions = make(map[string][]string)
		keys := perm.PermTypes()
		for _, k := range keys {
			actors := perm.GetPerm(k)
			temp := make([]string, 0, len(actors))
			for _, a := range actors {
				temp = append(temp, a.GetID().String())
			}
			out.Permissions[k] = temp
		}
	}
	return out
}

var errInvalidPerm = srverror.New(errors.New("perm request is invalid"), 404, "Not Found")

func pullPerm(w http.ResponseWriter, r *http.Request) types.PermissionI {

	vals := mux.Vars(r)
	var permobj types.PermissionI
	if vals["type"] == "group" {
		userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
		id, err := types.DecodeOwnerIDString(vals["id"])
		if err != nil {
			panic(srverror.New(err, 400, "Bad Group ID"))
		}
		g, err := userbase.Get(id)
		if err != nil {
			util.VerboseRequest(r, "error getting group for pull permission")
			panic(err)
		}
		var ok bool
		permobj, ok = g.(types.PermissionI)
		if !ok {
			panic(srverror.Basic(404, "Not Found", g.GetID().String()))
		}
	} else if vals["type"] == "file" {
		filebase := r.Context().Value(types.FILE).(database.Filebase)
		id, err := types.DecodeFileID(vals["id"])
		if err != nil {
			panic(srverror.New(err, 400, "Bad File ID"))
		}
		rec, err := filebase.Get(id)
		if err != nil {
			util.VerboseRequest(r, "error getting file for pull permission")
			panic(err)
		}
		return rec.(types.PermissionI)
	} else {
		panic(srverror.Basic(404, "Not Found", "Unrecognized Permission Type"))
	}
	return permobj
}

func getPermissions(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	// user := r.Context().Value(USER).(types.UserI)
	// permobj := pullPerm(w, r)
	bpr := buildPermissionReport(pullPerm(w, r),
		r.Context().Value(USER).(types.UserI))
	w.Set("owner", bpr.Owner)
	w.Set("isOwned", bpr.IsOwner)
	w.Set("permission", bpr.Permissions)
}

func setPermission(permval bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(USER).(types.UserI)
		permobj := pullPerm(w, r)
		if !permobj.GetOwner().Match(user) {
			panic(srverror.Basic(403, "Permission Denied", user.GetID().String()))
		}
		r.PostFormValue("id")
		for _, idstr := range r.PostForm["id"] {
			id, err := types.DecodeOwnerIDString(idstr)
			if err != nil {
				panic(srverror.New(err, 400, "Bad Target ID"))
			}
			if id.Type == 'p' {
				panic(srverror.Basic(400, "Cannot Manipulate Public Permission"))
			}
			target, err := r.Context().Value(types.OWNER).(database.Ownerbase).Get(id)
			if err != nil {
				util.VerboseRequest(r, "Failed to get owner from id to change permission for")
				panic(err)
			}
			permobj.SetPerm(target, "view", permval)
		}
		var err error
		switch v := permobj.(type) {
		case types.Owner:
			err = r.Context().Value(types.OWNER).(database.Ownerbase).Update(v)
		case types.FileI:
			err = r.Context().Value(types.FILE).(database.Filebase).Update(v)
		default:
			err = errInvalidPerm
		}
		if err != nil {
			util.VerboseRequest(r, "unable to update permissions")
			panic(err)
		}
		w.Write([]byte("Permission Updated"))
	}
}

func setPermissionPublic(permval bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(USER).(types.UserI)
		userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
		filebase := r.Context().Value(types.FILE).(database.Filebase)
		permobj := pullPerm(w, r)
		if !permobj.GetOwner().Match(user) || !user.GetRole("admin") {
			panic(srverror.Basic(403, "Permission Denied", user.GetID().String()))
		}
		permobj.SetPerm(types.Public, "view", permval)
		var err error
		switch v := permobj.(type) {
		case types.Owner:
			err = userbase.Update(v)
		case types.FileI:
			err = filebase.Update(v)
		default:
			err = errInvalidPerm
		}
		if err != nil {
			panic(err)
		}
		w.Write([]byte("Permission Updated"))
	}
}
