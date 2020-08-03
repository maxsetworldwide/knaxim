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
	"context"
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

func groupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupidstr := r.FormValue("group")
		if len(groupidstr) > 0 {
			groupid, err := types.DecodeOwnerIDString(groupidstr)
			if err != nil {
				panic(srverror.New(err, 400, "Corrupt Group id"))
			}
			group, err := r.Context().Value(types.OWNER).(database.Ownerbase).Get(groupid)
			if err != nil {
				panic(err)
			}
			if !group.Match(r.Context().Value(USER).(types.Owner)) {
				panic(srverror.Basic(403, "Not Group Member"))
			}
			r = r.WithContext(context.WithValue(r.Context(), GROUP, group))
		}
		next.ServeHTTP(w, r)
	})
}

func groupidMiddleware(checkmembership bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			groupidstr := mux.Vars(r)["id"]
			groupid, err := types.DecodeOwnerIDString(groupidstr)
			if err != nil {
				panic(srverror.New(err, 400, "Corrupt Group id"))
			}
			group, err := r.Context().Value(types.OWNER).(database.Ownerbase).Get(groupid)
			if err != nil {
				panic(err)
			}
			if checkmembership {
				if !group.Match(r.Context().Value(USER).(types.Owner)) {
					panic(srverror.Basic(403, "Not Group Member"))
				}
			}
			r = r.WithContext(context.WithValue(r.Context(), GROUP, group))
			next.ServeHTTP(w, r)
		})
	}
}

// AttachGroup adds api paths related to group actions
func AttachGroup(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(UserCookie)
	r.Use(srvjson.JSONResponse)
	r.Use(ParseBody)
	r.HandleFunc("/options", getGroups).Methods("GET")
	r.HandleFunc("/name/{name}", lookupGroup).Methods("GET")
	{
		r = r.NewRoute().Subrouter()
		r.Use(groupMiddleware)
		r.HandleFunc("", createGroup).Methods("PUT")
	}
	{
		r = r.NewRoute().Subrouter()
		r.Use(groupidMiddleware(false))
		r.HandleFunc("/{id}", groupinfo).Methods("GET")
	}
	{
		r = r.NewRoute().Subrouter()
		r.Use(groupidMiddleware(true))
		r.HandleFunc("/options/{id}", getGroupsGroups).Methods("GET")
		r.HandleFunc("/{id}/search", searchGroupFiles).Methods("GET")
		r.HandleFunc("/{id}/member", updateGroupMember(true)).Methods("POST")
		r.HandleFunc("/{id}/member", updateGroupMember(false)).Methods("DELETE")
	}
}

func createGroup(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	newname := r.FormValue("newname")
	if !validGroupName(newname) {
		panic(srverror.Basic(400, "invalid group name", newname))
	}
	ng := types.NewGroup(newname, owner)
	var err error
	if ng.ID, err = ownerbase.Reserve(ng.GetID(), ng.GetName()); err != nil {
		panic(err)
	}
	if err := ownerbase.Insert(ng); err != nil {
		panic(err)
	}

	w.Set("message", "X31-CPO Group Created")
	// w.Write([]byte("Group Created"))
}

func getGroups(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	user := r.Context().Value(USER).(types.UserI)
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)

	ogroups, mgroups, err := ownerbase.GetGroups(user.GetID())
	if err != nil {
		panic(err)
	}

	own := []GroupInformation{}
	member := []GroupInformation{}
	for _, gr := range ogroups {
		own = append(own, BuildGroupInfo(gr))
	}
	for _, gr := range mgroups {
		member = append(member, BuildGroupInfo(gr))
	}
	w.Set("own", own)
	w.Set("member", member)
}

func getGroupsGroups(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	owner := r.Context().Value(GROUP).(types.GroupI)
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)

	ogroups, mgroups, err := ownerbase.GetGroups(owner.GetID())
	if err != nil {
		panic(err)
	}

	own := []GroupInformation{}
	member := []GroupInformation{}
	for _, gr := range ogroups {
		own = append(own, BuildGroupInfo(gr))
	}
	for _, gr := range mgroups {
		member = append(member, BuildGroupInfo(gr))
	}
	w.Set("own", own)
	w.Set("member", member)
}

func groupinfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	group, ok := r.Context().Value(GROUP).(types.GroupI)
	if !ok {
		panic(srverror.Basic(404, "Not Found", "id was an owner but not a group"))
	}
	user := r.Context().Value(USER).(types.UserI)
	var result GroupInformation
	if group.Match(user) {
		result = BuildGroupInfo(group)
	} else {
		result = GroupInformation{
			ID:   group.GetID().String(),
			Name: group.GetName(),
		}
	}

	w.Set("id", result.ID)
	w.Set("name", result.Name)
	if result.Owner != "" {
		w.Set("owner", result.Owner)
	}
	if len(result.Members) != 0 {
		w.Set("members", result.Members)
	}
}

func searchGroupFiles(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	if err := r.ParseForm(); err != nil {
		panic(srverror.New(err, 400, "Unable to parse form data"))
	}
	if len(r.Form["find"]) == 0 {
		panic(srverror.Basic(400, "No Search Term"))
	}
	group := r.Context().Value(GROUP).(types.Owner)
	filebase := r.Context().Value(types.FILE).(database.Filebase)
	owned, err := filebase.GetOwned(group.GetID())
	if err != nil {
		panic(err)
	}
	accessible, err := filebase.GetPermKey(group.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]types.FileID, 0, len(owned)+len(accessible))
	for _, o := range owned {
		fids = append(fids, o.GetID())
	}
	for _, a := range accessible {
		fids = append(fids, a.GetID())
	}

	filters := make([]tag.FileTag, 0, len(r.Form["find"]))
	for _, f := range util.SplitSearch(r.Form["find"]...) {
		if len(f) > 0 {
			filters = append(filters, tag.FileTag{
				Tag: tag.Tag{
					Word: f,
					Type: tag.CONTENT | tag.SEARCH,
					Data: tag.Data{
						tag.SEARCH: map[string]interface{}{
							"regex":        true,
							"regexoptions": "i",
						},
					},
				},
			})
		}
	}
	if len(filters) == 0 {
		panic(srverror.Basic(400, "No Search Condition"))
	}
	result, err := r.Context().Value(types.TAG).(database.Tagbase).SearchFiles(fids, filters...)
	if err != nil {
		panic(err)
	}
	w.Set("matched", BuildSearchResponse(r, result).Files)
}

func updateGroupMember(add bool) func(http.ResponseWriter, *http.Request) {
	return func(out http.ResponseWriter, r *http.Request) {
		w := out.(*srvjson.ResponseWriter)
		group := r.Context().Value(GROUP).(types.GroupI)
		actor := r.Context().Value(USER).(types.Owner)

		if !group.GetOwner().Match(actor) {
			panic(srverror.Basic(403, "Not Owner", actor.GetName(), actor.GetID().String(), group.GetName(), group.GetID().String()))
		}
		r.FormValue("id")
		if len(r.Form["id"]) == 0 {
			panic(srverror.Basic(400, "Missing Member ID"))
		}
		ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)
		var targets []types.Owner
		for _, idstr := range r.Form["id"] {
			id, err := types.DecodeOwnerIDString(idstr)
			if err != nil {
				panic(srverror.New(err, 400, "Bad Member ID"))
			}
			mem, err := ownerbase.Get(id)
			if err != nil {
				panic(err)
			}
			if group.Equal(mem) {
				panic(srverror.Basic(400, "Attempted to add group to itself"))
			}
			targets = append(targets, mem)
		}
		for _, t := range targets {
			if add {
				group.AddMember(t)
			} else {
				group.RemoveMember(t)
			}
		}
		err := ownerbase.Update(group)
		if err != nil {
			panic(err)
		}

		w.Set("message", "updated members")
	}
}

func lookupGroup(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	vals := mux.Vars(r)
	match, err := ownerbase.FindGroupName(vals["name"])
	if err != nil {
		panic(err)
	}
	response := BuildGroupInfo(match)
	if !match.Match(r.Context().Value(USER).(types.Owner)) {
		response.Members = nil
	}

	w.Set("id", response.ID)
	w.Set("name", response.Name)
	w.Set("members", response.Members)
}
