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
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/email"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/passentropy"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

// AttachUser hooks up paths for handling user related requests
func AttachUser(r *mux.Router) {
	r = r.NewRoute().Subrouter()
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	if !config.V.PrivateMode {
		r.HandleFunc("", createUser).Methods("PUT")
	}
	// r.HandleFunc("/admin", createAdmin).Methods("PUT")
	r.HandleFunc("/login", loginUser).Methods("POST")
	r.HandleFunc("/reset", sendReset).Methods("PUT")
	r.HandleFunc("/reset", updateCredentialsReset).Methods("POST")
	{
		r = r.NewRoute().Subrouter()
		r.Use(UserCookie)
		if config.V.PrivateMode {
			r.HandleFunc("", createUser).Methods("PUT")
		}
		r.HandleFunc("", userInfo).Methods("GET")
		r.HandleFunc("/name/{name}", lookupUser).Methods("GET")
		r.HandleFunc("", signoutUser).Methods("DELETE")
		r.HandleFunc("/complete", completeUserInfo).Methods("GET")
		r.HandleFunc("/search", searchAllUserFiles).Methods("GET")
		r.HandleFunc("/pass", updateCredentials).Methods("POST")
		r.HandleFunc("/data", getUserData).Methods("GET")
	}
}

func sendReset(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	if len(username) == 0 {
		panic(srverror.Basic(400, "Missing Username"))
	}
	ob := r.Context().Value(types.OWNER).(database.Ownerbase)
	user, err := ob.FindUserName(username)
	if err != nil {
		panic(err)
	}
	key, err := ob.GetResetKey(user.GetID())
	if err != nil {
		panic(err)
	}
	err = email.SendResetEmail([]string{user.GetEmail()}, user.GetName(), config.V.Address, key)
	if err != nil {
		panic(srverror.New(err, 500, "Error H4", "unable to send reset email"))
	}
	w.Write([]byte("success"))
}

var updateCredentialsResetLock = &sync.Mutex{}

func updateCredentialsReset(w http.ResponseWriter, r *http.Request) {
	updateCredentialsResetLock.Lock()
	defer updateCredentialsResetLock.Unlock()
	key := r.FormValue("key")
	ob := r.Context().Value(types.OWNER).(database.Ownerbase)
	id, err := ob.CheckResetKey(key)
	if err != nil {
		panic(err)
	}
	owner, err := ob.Get(id)
	if err != nil {
		panic(err)
	}
	user, ok := owner.(types.UserI)
	if !ok {
		panic(srverror.Basic(404, "Not Found", "non-user owner assigned a key", key))
	}
	newpass := r.FormValue("newpass")
	if passentropy.Score(newpass) < passentropy.Char6Cap1num1 {
		panic(srverror.Basic(400, "Password not secure"))
	}
	user.SetLock(types.NewUserCredential(newpass))
	err = ob.Update(user)
	if err != nil {
		panic(err)
	}
	err = ob.DeleteResetKey(user.GetID())
	if err != nil {
		panic(err)
	}
	w.Write([]byte("password updated"))
}

func lookupUser(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	vals := mux.Vars(r)
	userName := vals["name"]
	if len(userName) == 0 {
		panic(srverror.Basic(400, "No user name"))
	}
	user, err := r.Context().Value(types.OWNER).(database.Ownerbase).FindUserName(userName)
	if err != nil {
		panic(err)
	}

	resp := BuildUserInfo(r, user)
	w.Set("id", resp.ID)
	w.Set("name", resp.Name)
	w.Set("data", resp.Data)
}

type dataUsage struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

func getUserData(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	user := r.Context().Value(USER).(types.UserI)
	userbase := r.Context().Value(types.OWNER).(database.Ownerbase)

	var du dataUsage
	var err error
	if du.Current, err = userbase.GetSpace(user.GetID()); err != nil {
		panic(err)
	}
	if du.Total, err = userbase.GetTotalSpace(user.GetID()); err != nil {
		panic(err)
	}

	w.Set("current", du.Current)
	w.Set("total", du.Total)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	if config.V.PrivateMode {
		user := r.Context().Value(USER).(types.UserI)
		if !user.GetRole("admin") {
			panic(srverror.Basic(400, "Private Mode, contact site administrator to create account."))
		}
	}
	errorLog := fmt.Sprintf("invalid values: user: %s, email: %s, pass: %s", r.FormValue("name"), r.FormValue("email"), r.FormValue("pass"))
	if !validUserName(r.FormValue("name")) {
		panic(srverror.Basic(400, "Invalid Username", errorLog))
	}
	if passentropy.Score(r.FormValue("pass")) < passentropy.Char6Cap1num1 {
		panic(srverror.Basic(400, "Password not strong enough", errorLog))
	}
	if !validEmail(r.FormValue("email")) {
		panic(srverror.Basic(400, "Invalid Email", errorLog))
	}
	nuser := types.NewUser(r.FormValue("name"), r.FormValue("pass"), r.FormValue("email"))
	var err error
	if nuser.ID, err = ownerbase.Reserve(nuser.ID, nuser.Name); err != nil {
		panic(err)
	}
	if err := ownerbase.Insert(nuser); err != nil {
		panic(err)
	}
	w.Write([]byte(nuser.ID.String()))
}

func createAdmin(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("adminkey") != config.V.AdminKey {
		panic(srverror.Basic(400, "Bad Request", "incorrect admin key", r.FormValue("adminkey"), config.V.AdminKey))
	}
	ownerbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	errorLog := fmt.Sprintf("invalid values: user: %s, email: %s, pass: %s", r.FormValue("name"), r.FormValue("email"), r.FormValue("pass"))
	if !validUserName(r.FormValue("name")) {
		panic(srverror.Basic(400, "Invalid Username", errorLog))
	}
	if passentropy.Score(r.FormValue("pass")) < passentropy.Char6Cap1num1 {
		panic(srverror.Basic(400, "Password not strong enough", errorLog))
	}
	if !validEmail(r.FormValue("email")) {
		panic(srverror.Basic(400, "Invalid Email", errorLog))
	}
	nuser := types.NewUser(r.FormValue("name"), r.FormValue("pass"), r.FormValue("email"))
	nuser.SetRole("admin", true)
	var err error
	if nuser.ID, err = ownerbase.Reserve(nuser.ID, nuser.Name); err != nil {
		panic(err)
	}
	if err := ownerbase.Insert(nuser); err != nil {
		panic(err)
	}
	w.Write([]byte(nuser.ID.String()))
}

func userInfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	user := r.Context().Value(USER).(types.UserI)
	userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	if r.FormValue("id") != "" {
		oid, err := types.DecodeOwnerIDString(r.FormValue("id"))
		if err != nil {
			panic(srverror.New(err, 400, "Unable to Decode UserID"))
		}
		util.VerboseRequest(r, "getting owner id: %+v", oid)
		target, err := userbase.Get(oid)
		if err != nil {
			panic(err)
		}
		var ok bool
		if user, ok = target.(types.UserI); !ok {
			panic(srverror.Basic(404, "ID is not a user"))
		}
	}

	resp := BuildUserInfo(r, user)
	w.Set("id", resp.ID)
	w.Set("name", resp.Name)
	w.Set("data", resp.Data)
}

func searchAllUserFiles(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	user := r.Context().Value(USER).(types.Owner)
	filebase := r.Context().Value(types.FILE).(database.Filebase)
	if err := r.ParseForm(); err != nil {
		panic(srverror.New(err, 400, "Unable to parse form data"))
	}
	if len(r.Form["find"]) == 0 {
		panic(srverror.Basic(404, "Not Found", "no search term"))
	}
	filters := make([]tag.FileTag, 0, len(r.Form["find"]))
	for _, f := range util.SplitSearch(r.Form["find"]...) {
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
	owned, err := filebase.GetOwned(user.GetID())
	if err != nil {
		panic(err)
	}
	viewable, err := filebase.GetPermKey(user.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]types.FileID, 0, len(owned)+len(viewable))
	for _, o := range owned {
		fids = append(fids, o.GetID())
	}
	for _, v := range viewable {
		fids = append(fids, v.GetID())
	}
	fids, err = r.Context().Value(types.TAG).(database.Tagbase).SearchFiles(fids, filters...)
	if err != nil {
		panic(err)
	}

	w.Set("matched", BuildSearchResponse(r, fids).Files)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	user, err := userbase.FindUserName(r.PostFormValue("name"))
	if err != nil {
		util.VerboseRequest(r, "request username: %s", r.PostFormValue("name"))
		panic(err)
	}
	if user.GetLock().Valid(map[string]interface{}{"pass": r.PostFormValue("pass")}) {
		cs := user.NewCookies(
			time.Now().Add(config.V.UserTimeouts.Inactivity.Duration),
			time.Now().Add(config.V.UserTimeouts.Total.Duration),
		)
		if err = userbase.Update(user); err != nil {
			panic(err)
		}
		for _, c := range cs {
			c.Path = "/api"
			http.SetCookie(w, c)
		}
		w.Write([]byte("logged in"))
	} else {
		panic(srverror.Basic(404, "Not Found", "Password Wrong"))
	}
}

func signoutUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(types.UserI)
	userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	user.NewCookies(time.Time{}, time.Time{})
	if err := userbase.Update(user); err != nil {
		panic(err)
	}
	w.Write([]byte("Signed Out"))
}

func updateCredentials(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(types.UserI)
	userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
	if !user.GetLock().Valid(map[string]interface{}{"pass": r.FormValue("oldpass")}) {
		panic(srverror.Basic(404, "Not Found"))
	}
	if passentropy.Score(r.FormValue("newpass")) < passentropy.Char6Cap1num1 {
		panic(srverror.Basic(400, "Password not strong enough"))
	}
	user.SetLock(types.NewUserCredential(r.FormValue("newpass")))
	if err := userbase.Update(user); err != nil {
		panic(err)
	}
	w.Write([]byte("password updated"))
}
