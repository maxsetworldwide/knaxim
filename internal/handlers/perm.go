package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.maxset.io/server/knaxim/database"
	"git.maxset.io/server/knaxim/database/filehash"
	"git.maxset.io/server/knaxim/srverror"

	"github.com/gorilla/mux"
)

func setupPerm(r *mux.Router) {
	r.Use(cookieMiddleware)
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

func buildPermissionReport(perm database.PermissionI, actor database.Owner) permissionReport {
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

var errInvalidPerm = srverror.New(errors.New("prem request is invalid"), 404, "Not Found")

func pullPerm(w http.ResponseWriter, r *http.Request) database.PermissionI {

	vals := mux.Vars(r)
	var permobj database.PermissionI
	if vals["type"] == "group" {
		userbase := r.Context().Value(database.OWNER).(database.Ownerbase)
		id, err := database.DecodeObjectIDString(vals["id"])
		if err != nil {
			panic(srverror.New(err, 400, "Bad Group ID"))
		}
		g, err := userbase.Get(id)
		if err != nil {
			verboseRequest(r, "error getting group for pull permission")
			panic(err)
		}
		var ok bool
		permobj, ok = g.(database.PermissionI)
		if !ok {
			panic(srverror.Basic(404, "Not Found", g.GetID().String()))
		}
	} else if vals["type"] == "file" {
		filebase := r.Context().Value(database.FILE).(database.Filebase)
		id, err := filehash.DecodeFileID(vals["id"])
		if err != nil {
			panic(srverror.New(err, 400, "Bad File ID"))
		}
		rec, err := filebase.Get(id)
		if err != nil {
			verboseRequest(r, "error getting file for pull permission")
			panic(err)
		}
		return rec.(database.PermissionI)
	} else {
		panic(srverror.Basic(404, "Not Found", "Unrecognized Permission Type"))
	}
	return permobj
}

func getPermissions(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(database.UserI)
	permobj := pullPerm(w, r)
	if err := json.NewEncoder(w).Encode(buildPermissionReport(permobj, user)); err != nil {
		panic(srverror.New(err, 500, "Server Error", "getPermission encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func setPermission(permval bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(USER).(database.UserI)
		permobj := pullPerm(w, r)
		if !permobj.GetOwner().Match(user) {
			panic(srverror.Basic(403, "Permission Denied", user.GetID().String()))
		}
		r.PostFormValue("id")
		for _, idstr := range r.PostForm["id"] {
			id, err := database.DecodeObjectIDString(idstr)
			if err != nil {
				panic(srverror.New(err, 400, "Bad Target ID"))
			}
			if id.Type == 'p' {
				panic(srverror.Basic(400, "Cannot Manipulate Public Permission"))
			}
			target, err := r.Context().Value(database.OWNER).(database.Ownerbase).Get(id)
			if err != nil {
				verboseRequest(r, "Failed to get owner from id to change permission for")
				panic(err)
			}
			permobj.SetPerm(target, "view", permval)
		}
		var err error
		switch v := permobj.(type) {
		case database.Owner:
			err = r.Context().Value(database.OWNER).(database.Ownerbase).Update(v)
		case database.FileI:
			err = r.Context().Value(database.FILE).(database.Filebase).Update(v)
		default:
			err = errInvalidPerm
		}
		if err != nil {
			verboseRequest(r, "unable to update permissions")
			panic(err)
		}
		w.Write([]byte("Permission Updated"))
	}
}

func setPermissionPublic(permval bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(USER).(database.UserI)
		userbase := r.Context().Value(database.OWNER).(database.Ownerbase)
		filebase := r.Context().Value(database.FILE).(database.Filebase)
		permobj := pullPerm(w, r)
		if !permobj.GetOwner().Match(user) && user.GetRole("admin") {
			panic(srverror.Basic(403, "Permission Denied", user.GetID().String()))
		}
		permobj.SetPerm(database.Public, "view", permval)
		var err error
		switch v := permobj.(type) {
		case database.Owner:
			err = userbase.Update(v)
		case database.FileI:
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
