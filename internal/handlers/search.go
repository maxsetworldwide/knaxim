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
	"encoding/json"
	"net/http"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/query"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

// AttachSearch adds handlers for searching the tags of files
func AttachSearch(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(UserCookie)
	r.Use(srvjson.JSONResponse)
	r.HandleFunc("/tags", searchFileTags).Methods("POST")
}

func searchFileTags(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var query query.Q
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		panic(srverror.New(err, 400, "Malformed Query, type 1"))
	}
	user := r.Context().Value(USER).(types.Owner)
	for _, c := range query.Context {
		if access, err := c.CheckAccess(user, r.Context().Value(types.DATABASE).(database.Database), "view"); !access {
			if err != nil {
				serr, ok := err.(srverror.Error)
				if ok {
					panic(serr)
				} else {
					panic(srverror.New(err, 500, "Error H3", "Unable to check permisison"))
				}
			}
			panic(srverror.Basic(403, "Access Denied"))
		}
	}
	matches, err := query.FindMatching(r.Context(), config.DB)
	if err != nil {
		switch e := err.(type) {
		case srverror.Error:
			panic(e)
		case error:
			panic(srverror.New(e, 400, "Malformed Query, type 2"))
		}
	}
	w.Set("matched", BuildSearchResponse(r, matches).Files)
}
