/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
