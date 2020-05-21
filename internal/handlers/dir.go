package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

// AttachDir is for connecting folder related api paths
func AttachDir(r *mux.Router) {
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.Use(groupMiddleware)
	//r.HandleFunc("/dynamic", createDynDir).Methods("PUT")
	r.HandleFunc("", createDir).Methods("PUT")
	r.HandleFunc("", getDirs).Methods("GET")
	r.HandleFunc("/{id}", dirInfo).Methods("GET")
	r.HandleFunc("/{id}/search", searchDir).Methods("GET")
	r.HandleFunc("/{id}/content", adjustDir(true)).Methods("POST")
	r.HandleFunc("/{id}/content", adjustDir(false)).Methods("DELETE")
	//r.HandleFunc("/{id}/refresh", refreshDynDir).Methods("POST")
	r.HandleFunc("/{id}", deleteDir).Methods("DELETE")
}

var dirflag = "d"

func getDirs(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var owner types.Owner
	if gr := r.Context().Value(GROUP); gr != nil {
		owner = gr.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}

	var dirs []string
	if tags, err := r.Context().Value(types.TAG).(database.Tagbase).GetAll(
		tag.USER,
		owner.GetID(),
	); err == nil {
		folderset := make(map[string]bool)
		for _, t := range tags {
			if !folderset[t.Word] {
				folderset[t.Word] = true
				dirs = append(dirs, t.Word)
			}
		}
	} else {
		panic(err)
	}

	w.Set("folders", dirs)
}

func createDir(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	filebase := r.Context().Value(types.FILE).(database.Filebase)
	nname := r.FormValue("newname")
	if !validDirName(nname) {
		panic(srverror.Basic(400, "Invalid Directory Name"))
	}
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	var files []types.FileI
	if r.Form["content"] != nil {
		for _, fidstr := range r.Form["content"] {
			fid, err := types.DecodeFileID(fidstr)
			if err != nil {
				util.VerboseRequest(r, "unable to decode: %s, because: %v", fidstr, err)
				continue
			}
			file, err := filebase.Get(fid) //TODO make a GetAll to use here
			if err == nil {
				files = append(files, file)
			}
		}
	}
	if len(files) == 0 {
		w.Set("id", nname)
		w.Set("affectedFiles", 0)
		return
	}
	tagbase := r.Context().Value(types.TAG).(database.Tagbase)
	for _, file := range files {
		err := tagbase.Upsert(tag.FileTag{
			File:  file.GetID(),
			Owner: owner.GetID(),
			Tag: tag.Tag{
				Word: nname,
				Type: tag.USER,
			},
		})
		if err != nil {
			panic(srverror.New(err, 500, "Server Error", "Unable to add user tag"))
		}
	}

	w.Set("id", nname)
	w.Set("affectedFiles", len(files))
}

func dirInfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	vals := mux.Vars(r)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	tagbase := r.Context().Value(types.TAG).(database.Tagbase)
	tags, err := tagbase.GetAll(tag.USER, owner.GetID())
	if err != nil {
		if se := err.(srverror.Error); se.Status() == errors.ErrNoResults.Status() {
			w.WriteHeader(se.Status())
		} else {
			panic(srverror.New(err, 500, "Server Error", "unable to get file tags"))
		}
	}
	var filematches []types.FileID
	for _, t := range tags {
		if t.Word == vals["id"] {
			filematches = append(filematches, t.File)
		}
	}

	w.Set("name", vals["id"])
	w.Set("files", filematches)
}

func adjustDir(add bool) func(http.ResponseWriter, *http.Request) {
	return func(out http.ResponseWriter, r *http.Request) {
		w := out.(*srvjson.ResponseWriter)

		tagbase := r.Context().Value(types.TAG).(database.Tagbase)
		var owner types.Owner
		if group := r.Context().Value(GROUP); group != nil {
			owner = group.(types.Owner)
		} else {
			owner = r.Context().Value(USER).(types.Owner)
		}
		vals := mux.Vars(r)
		dirtagname := vals["id"]
		if len(dirtagname) == 0 {
			panic(srverror.Basic(400, "Please include a directory name"))
		}
		fidstrs := r.PostForm["id"]
		if len(fidstrs) == 0 {
			panic(srverror.Basic(400, "Please include file IDs for the directory"))
		}
		var fids []types.FileID
		for _, fidstr := range fidstrs {
			fid, err := types.DecodeFileID(fidstr)
			if err != nil {
				panic(srverror.New(err, 400, "Corrupt File ID"))
			}
			fids = append(fids, fid)
		}
		for _, fid := range fids {
			dirtag := tag.FileTag{
				File:  fid,
				Owner: owner.GetID(),
				Tag: tag.Tag{
					Word: dirtagname,
					Type: tag.USER,
				},
			}
			var err error
			if add {
				err = tagbase.Upsert(dirtag)
			} else {
				err = tagbase.Remove(dirtag)
			}
			if err != nil {
				panic(err)
			}
		}

		w.Set("message", "Complete")
	}
}

func searchDir(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var tagbase = r.Context().Value(types.TAG).(database.Tagbase)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	filters := make([]tag.FileTag, 0, 1+len(r.Form["find"]))
	filters = append(filters, tag.FileTag{
		Owner: owner.GetID(),
		Tag: tag.Tag{
			Word: vals["id"],
			Type: tag.USER,
		},
	})
	for _, find := range r.Form["find"] {
		// TODO: process find strings into regex
		filters = append(filters, tag.FileTag{
			Tag: tag.Tag{
				Word: find,
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
	ownedids, err := tagbase.SearchOwned(owner.GetID(), filters...)
	if err != nil {
		panic(err)
	}
	viewids, err := tagbase.SearchAccess(owner.GetID(), "view", filters...)
	if err != nil {
		panic(err)
	}
	matches := make([]types.FileID, 0, len(ownedids)+len(viewids))
	for _, o := range ownedids {
		matches = append(matches, o)
	}
	for _, v := range viewids {
		matches = append(matches, v)
	}
	w.Set("matches", matches)
}

func deleteDir(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	ftags, err := r.Context().Value(types.TAG).(database.Tagbase).GetAll(tag.USER, owner.GetID())
	if err != nil {
		panic(err)
	}
	var targettags []tag.FileTag
	for _, ft := range ftags {
		if ft.Word == vals["id"] {
			targettags = append(targettags, ft)
		}
	}
	err = r.Context().Value(types.TAG).(database.Tagbase).Remove(targettags...)
	if err != nil {
		panic(err)
	}
	w.Set("message", "Complete")
}
