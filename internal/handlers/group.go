package main

import (
	"context"
	"encoding/json"
	"net/http"

	"git.maxset.io/server/knaxim/database"
	"git.maxset.io/server/knaxim/database/filehash"
	"git.maxset.io/server/knaxim/database/tag"

	"git.maxset.io/server/knaxim/srverror"

	"github.com/gorilla/mux"
)

func groupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupidstr := r.FormValue("group")
		if len(groupidstr) > 0 {
			groupid, err := database.DecodeObjectIDString(groupidstr)
			if err != nil {
				panic(srverror.New(err, 400, "Corrupt Group id"))
			}
			group, err := r.Context().Value(database.OWNER).(database.Ownerbase).Get(groupid)
			if err != nil {
				panic(err)
			}
			if !group.Match(r.Context().Value(USER).(database.Owner)) {
				panic(srverror.Basic(403, "Not Group Member"))
			}
			r = r.WithContext(context.WithValue(r.Context(), GROUP, group))
		}
		next.ServeHTTP(w, r)
	})
}

func groupidMiddleware(next http.Handler, checkmembership bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupidstr := mux.Vars(r)["id"]
		groupid, err := database.DecodeObjectIDString(groupidstr)
		if err != nil {
			panic(srverror.New(err, 400, "Corrupt Group id"))
		}
		group, err := r.Context().Value(database.OWNER).(database.Ownerbase).Get(groupid)
		if err != nil {
			panic(err)
		}
		if checkmembership {
			if !group.Match(r.Context().Value(USER).(database.Owner)) {
				panic(srverror.Basic(403, "Not Group Member"))
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), GROUP, group))
		next.ServeHTTP(w, r)
	})
}

// server sends: /api/group
func setupGroup(r *mux.Router) {
	r.Use(cookieMiddleware)
	r.Handle("", groupMiddleware(http.HandlerFunc(createGroup))).Methods("PUT")
	r.HandleFunc("/options", getGroups).Methods("GET")
	r.Handle("/options/{id}", groupidMiddleware(http.HandlerFunc(getGroupsGroups), true)).Methods("GET")
	r.Handle("/{id}", groupidMiddleware(http.HandlerFunc(groupinfo), false)).Methods("GET")
	r.Handle("/{id}/search", groupidMiddleware(http.HandlerFunc(searchGroupFiles), true)).Methods("GET")
	r.Handle("/{id}/member", groupidMiddleware(http.HandlerFunc(updateGroupMember(true)), true)).Methods("POST")
	r.Handle("/{id}/member", groupidMiddleware(http.HandlerFunc(updateGroupMember(false)), true)).Methods("DELETE")
	r.HandleFunc("/name/{name}", lookupGroup).Methods("GET")
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)
	newname := r.FormValue("newname")
	if !validGroupName(newname) {
		panic(srverror.Basic(400, "Bad Request", "invalid group name", newname))
	}
	ng := database.NewGroup(newname, owner)
	var err error
	if ng.ID, err = ownerbase.Reserve(ng.GetID(), ng.GetName()); err != nil {
		panic(err)
	}
	if err := ownerbase.Insert(ng); err != nil {
		panic(err)
	}
	w.Write([]byte("Group Created"))
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(database.UserI)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)

	ogroups, mgroups, err := ownerbase.GetGroups(user.GetID())
	if err != nil {
		panic(err)
	}
	result := make(map[string][]GroupInformation)
	for _, gr := range ogroups {
		result["own"] = append(result["own"], BuildGroupInfo(gr))
	}
	for _, gr := range mgroups {
		result["member"] = append(result["member"], BuildGroupInfo(gr))
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(srverror.New(err, 500, "Server Error", "getGroups encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func getGroupsGroups(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(GROUP).(database.GroupI)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)

	ogroups, mgroups, err := ownerbase.GetGroups(owner.GetID())
	if err != nil {
		panic(err)
	}
	result := make(map[string][]GroupInformation)
	for _, gr := range ogroups {
		result["own"] = append(result["own"], BuildGroupInfo(gr))
	}
	for _, gr := range mgroups {
		result["member"] = append(result["member"], BuildGroupInfo(gr))
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(srverror.New(err, 500, "Server Error", "getGroups encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func groupinfo(w http.ResponseWriter, r *http.Request) {
	group := r.Context().Value(GROUP).(database.GroupI)
	user := r.Context().Value(USER).(database.UserI)
	var result GroupInformation
	if group.Match(user) {
		result = BuildGroupInfo(group)
	} else {
		result = GroupInformation{
			ID:   group.GetID().String(),
			Name: group.GetName(),
		}
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(srverror.New(err, 500, "Server Error", "groupinfo encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func searchGroupFiles(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(srverror.New(err, 400, "Bad Request", "Unable to parse form data"))
	}
	if len(r.Form["find"]) == 0 {
		panic(srverror.Basic(400, "No Search Term"))
	}
	group := r.Context().Value(GROUP).(database.Owner)
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	owned, err := filebase.GetOwned(group.GetID())
	if err != nil {
		panic(err)
	}
	accessible, err := filebase.GetPermKey(group.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]filehash.FileID, 0, len(owned)+len(accessible))
	for _, o := range owned {
		fids = append(fids, o.GetID())
	}
	for _, a := range accessible {
		fids = append(fids, a.GetID())
	}

	filters := make([]tag.Tag, 0, len(r.Form["find"]))
	for _, f := range splitSearch(r.Form["find"]...) {
		if len(f) > 0 {
			filters = append(filters, tag.Tag{
				Word: f,
				Type: tag.CONTENT,
			})
		}
	}
	if len(filters) == 0 {
		panic(srverror.Basic(400, "No Search Condition"))
	}
	result, _, err := r.Context().Value(database.TAG).(database.Tagbase).GetFiles(filters, fids...)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(BuildSearchResponse(r, result)); err != nil {
		panic(srverror.New(err, 500, "Failed to encode responce"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func updateGroupMember(add bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.FormValue("id")
		if len(r.Form["id"]) == 0 {
			panic(srverror.Basic(400, "Missing Member ID"))
		}
		group := r.Context().Value(GROUP).(database.GroupI)
		ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)
		for _, idstr := range r.Form["id"] {
			id, err := database.DecodeObjectIDString(idstr)
			if err != nil {
				panic(srverror.New(err, 400, "Bad Member ID"))
			}
			mem, err := ownerbase.Get(id)
			if err != nil {
				panic(err)
			}
			if add {
				group.AddMember(mem)
			} else {
				group.RemoveMember(mem)
			}
		}
		err := ownerbase.Update(group)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("updated members"))
	}
}

func lookupGroup(w http.ResponseWriter, r *http.Request) {
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)
	vals := mux.Vars(r)
	match, err := ownerbase.FindGroupName(vals["name"])
	if err != nil {
		panic(err)
	}
	responce := BuildGroupInfo(match)
	if !match.Match(r.Context().Value(USER).(database.Owner)) {
		responce.Members = nil
	}
	if err = json.NewEncoder(w).Encode(responce); err != nil {
		panic(srverror.New(err, 500, "Server Error", "unable to encode json"))
	}
	w.Header().Set("Content-Type", "application/json")
}
