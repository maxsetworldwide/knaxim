package handlers

import (
	"context"
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

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
func AttachGroup(r *mux.Router) {
	r.Use(UserCookie)
	r.Handle("", groupMiddleware(http.HandlerFunc(createGroup))).Methods("PUT")
	r.HandleFunc("/options", getGroups).Methods("GET")
	r.Handle("/options/{id}", groupidMiddleware(http.HandlerFunc(getGroupsGroups), true)).Methods("GET")
	r.Handle("/{id}", groupidMiddleware(http.HandlerFunc(groupinfo), false)).Methods("GET")
	r.Handle("/{id}/search", groupidMiddleware(http.HandlerFunc(searchGroupFiles), true)).Methods("GET")
	r.Handle("/{id}/member", groupidMiddleware(http.HandlerFunc(updateGroupMember(true)), true)).Methods("POST")
	r.Handle("/{id}/member", groupidMiddleware(http.HandlerFunc(updateGroupMember(false)), true)).Methods("DELETE")
	r.HandleFunc("/name/{name}", lookupGroup).Methods("GET")
}

func createGroup(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
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

	w.Set("message", "X31-CPO Group Created")
	// w.Write([]byte("Group Created"))
}

func getGroups(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	user := r.Context().Value(USER).(database.UserI)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)

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

	owner := r.Context().Value(GROUP).(database.GroupI)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)

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
	group, ok := r.Context().Value(GROUP).(database.GroupI)
	if !ok {
		panic(srverror.Basic(404, "Not Found", "id was an owner but not a group"))
	}
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
	for _, f := range util.SplitSearch(r.Form["find"]...) {
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
	w.Set("matched", BuildSearchResponse(r, result).Files)
}

func updateGroupMember(add bool) func(http.ResponseWriter, *http.Request) {
	return func(out http.ResponseWriter, r *http.Request) {
		w := out.(*srvjson.ResponseWriter)

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

		w.Set("message", "updated members")
	}
}

func lookupGroup(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)
	vals := mux.Vars(r)
	match, err := ownerbase.FindGroupName(vals["name"])
	if err != nil {
		panic(err)
	}
	response := BuildGroupInfo(match)
	if !match.Match(r.Context().Value(USER).(database.Owner)) {
		response.Members = nil
	}

	w.Set("id", response.ID)
	w.Set("name", response.Name)
	w.Set("members", response.Members)
}
