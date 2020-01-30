package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

func AttachAcronym(r *mux.Router) {
	r.Use(UserCookie)

	r.HandleFunc("/{acronym}", getAcronym).Methods("GET")
}

func getAcronym(out http.ResponseWriter, r *http.Request) {
	w, ok := out.(*srvjson.ResponseWriter)
	if !ok {
		panic(srverror.Basic(500, "Server Error", "expecting *srvjson.ResponseWriter"))
	}
	vals := mux.Vars(r)
	matches, err := r.Context().Value(database.ACRONYM).(database.Acronymbase).Get(vals["acronym"])
	if err != nil {
		panic(err)
	}
	w.Set("matched", matches)
}
