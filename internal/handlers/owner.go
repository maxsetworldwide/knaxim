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
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

// AttachOwner adds api paths related to owner actions
func AttachOwner(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(srvjson.JSONResponse)

	r.HandleFunc("/id/{id}", getOwner).Methods("GET")
	r.HandleFunc("/name/{name}", lookupName).Methods("GET")
}

func getOwner(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	vals := mux.Vars(r)

	id, err := types.DecodeOwnerIDString(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Invalid owner id"))
	}
	o, err := r.Context().Value(types.OWNER).(database.Ownerbase).Get(id)
	if err != nil {
		panic(err)
	}
	w.Set("id", o.GetID())
	w.Set("name", o.GetName())
	switch o.(type) {
	case *types.User:
		w.Set("type", "user")
	case *types.Group:
		w.Set("type", "group")
	default:
		if o.Equal(types.Public) {
			w.Set("type", "public")
		} else {
			w.Set("type", "unknown")
		}
	}
}

func lookupName(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	vals := mux.Vars(r)

	var o types.Owner
	var err error
	if o, err = r.Context().Value(types.OWNER).(database.Ownerbase).FindUserName(vals["name"]); err != nil {
		if serr, ok := err.(srverror.Error); ok && serr.Status() == 404 {
			if o, err = r.Context().Value(types.OWNER).(database.Ownerbase).FindGroupName(vals["name"]); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	w.Set("id", o.GetID())
	w.Set("name", o.GetName())
	switch o.(type) {
	case types.UserI:
		w.Set("type", "user")
	case types.GroupI:
		w.Set("type", "group")
	}
}
