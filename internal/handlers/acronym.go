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

// AttachAcronym is to add api paths related to acronyms
func AttachAcronym(r *mux.Router) {
	r = r.NewRoute().Subrouter()
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(UserCookie)

	r.HandleFunc("/{acronym}", getAcronym).Methods("GET")
}

func getAcronym(out http.ResponseWriter, r *http.Request) {
	w, ok := out.(*srvjson.ResponseWriter)
	if !ok {
		panic(srverror.Basic(500, "Error H5", "expecting *srvjson.ResponseWriter"))
	}
	vals := mux.Vars(r)
	matches, err := r.Context().Value(types.ACRONYM).(database.Acronymbase).Get(vals["acronym"])
	if err != nil {
		panic(err)
	}
	w.Set("matched", matches)
}
