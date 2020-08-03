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
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

// AttachPublic is for searching public available files
func AttachPublic(r *mux.Router) {
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.HandleFunc("/search", searchPublic).Methods("GET")
}

func searchPublic(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	if len(r.Form["find"]) == 0 {
		panic(srverror.Basic(400, "No Search Term"))
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
	publicfiles, err := r.Context().Value(types.FILE).(database.Filebase).GetPermKey(types.Public.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]types.FileID, 0, len(publicfiles))
	for _, pf := range publicfiles {
		fids = append(fids, pf.GetID())
	}
	fids, err = r.Context().Value(types.TAG).(database.Tagbase).SearchFiles(fids, filters...)
	if err != nil {
		panic(err)
	}
	w.Set("matched", BuildSearchResponse(r, fids).Files)
}
