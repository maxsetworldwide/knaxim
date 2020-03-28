package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
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
	publicfiles, err := r.Context().Value(database.FILE).(types.Filebase).GetPermKey(types.Public.GetID(), "view")
	if err != nil {
		panic(err)
	}
	fids := make([]types.FileID, 0, len(publicfiles))
	for _, pf := range publicfiles {
		fids = append(fids, pf.GetID())
	}
	fids, _, err = r.Context().Value(types.TAG).(types.Tagbase).GetFiles(filters, fids...)
	if err != nil {
		panic(err)
	}
	w.Set("matched", BuildSearchResponse(r, fids).Files)
}
