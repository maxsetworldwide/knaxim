package srvjson

import (
	"net/http"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

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
