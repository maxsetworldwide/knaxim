package handlers

import "github.com/gorilla/mux"

func AttachSearch(r *mux.Router) {
	r.Use(UserCookie)

	r.HandleFunc("/tags", searchFileTags).Methods("POST")
}
