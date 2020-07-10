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
