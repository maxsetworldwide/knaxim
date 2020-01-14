package main

import (
	"encoding/json"
	"net/http"

	"git.maxset.io/server/knaxim/database"
	"git.maxset.io/server/knaxim/srverror"
	"github.com/gorilla/mux"
)

func setupAcronym(r *mux.Router) {
	r.Use(cookieMiddleware)

	r.HandleFunc("/{acronym}", getAcronym).Methods("GET")
}

func getAcronym(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)
	matches, err := r.Context().Value(database.ACRONYM).(database.Acronymbase).Get(vals["acronym"])
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"matched": matches,
	}); err != nil {
		panic(srverror.New(err, 500, "Server Error", "Unable to encode acronym matches"))
	}
}
