package handlers

import (
	"encoding/json"
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"github.com/gorilla/mux"
)

func AttachPublic(r *mux.Router) {
	r.Use(UserCookie)
	r.HandleFunc("/search", searchPublic).Methods("GET")
}

func searchPublic(w http.ResponseWriter, r *http.Request) {
	if len(r.Form["find"]) == 0 {
		panic(srverror.Basic(400, "No Search Term"))
	}
	filters := make([]tag.Tag, 0, len(r.Form["find"]))
	for _, f := range util.SplitSearch(r.Form["find"]...) {
		if len(f) > 0 {
			filters = append(filters, tag.Tag{
				Word: f,
				Type: tag.CONTENT,
			})
		}
	}
	if len(filters) == 0 {
		panic(srverror.Basic(400, "No Search Condition"))
	}
	publicfiles, err := r.Context().Value(database.FILE).(database.Filebase).GetPermKey(database.Public.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]filehash.FileID, 0, len(publicfiles))
	for _, pf := range publicfiles {
		fids = append(fids, pf.GetID())
	}
	fids, _, err = r.Context().Value(database.TAG).(database.Tagbase).GetFiles(filters, fids...)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(BuildSearchResponse(r, fids)); err != nil {
		panic(srverror.New(err, 500, "Failed to encode responce"))
	}
	w.Header().Add("Content-Type", "application/json")
}
