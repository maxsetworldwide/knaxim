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

package srvjson

import (
	"net/http"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

// JSONResponse is a middleware function that replaces the incoming
// http.ResponseWriter with a srvjson.ResponseWriter
func JSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsw := NewRW(w)
		jsw.Set("Copywrite", "Maxset Worldwide Inc. 2020")
		jsw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(jsw, r)
		if err := jsw.Flush(); err != nil {
			panic(srverror.New(err, 500, "Server Error", "Unable to encode response to json"))
		}
	})
}
