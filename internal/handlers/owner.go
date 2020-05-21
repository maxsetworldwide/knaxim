package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

// AttachOwner adds api paths related to owner actions
func AttachOwner(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(srvjson.JSONResponse)

	r.HandleFunc("/{id}", getOwner).Methods("GET")
}

func getOwner(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	vals := mux.Vars(r)

	id, err := types.DecodeOwnerIDString(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Invalid owner id"))
	}
	o, err := r.Context().Value(types.OWNER).(database.Ownerbase).Get(id)
	if err != nil {
		panic(err)
	}
	w.Set("id", o.GetID())
	w.Set("name", o.GetName())
	switch o.(type) {
	case *types.User:
		w.Set("type", "user")
	case *types.Group:
		w.Set("type", "group")
	default:
		if o.Equal(types.Public) {
			w.Set("type", "public")
		} else {
			w.Set("type", "unknown")
		}
	}
}
