package handlers

import (
	"net/http"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/tag"
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

var missingContextErr = srverror.Basic(500, "Unable to access search context")

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
	if tags, err := r.Context().Value(types.TAG).(types.Tagbase).SearchData(
		tag.USER,
		tag.Data{
			tag.USER: map[string]string{
				owner.GetID().String(): dirflag},
		},
	); err == nil {
		for _, t := range tags {
			dirs = append(dirs, t.Word)
		}
	} else {
		panic(err)
	}

	w.Set("folders", dirs)
}

func createDir(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	filebase := r.Context().Value(types.FILE).(types.Filebase)
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
	dirtag := tag.Tag{
		Word: nname,
		Type: tag.USER,
		Data: tag.Data{
			tag.USER: map[string]string{
				owner.GetID().String(): dirflag,
			},
		},
	}
	tagbase := r.Context().Value(types.TAG).(types.Tagbase)
	for _, file := range files {
		err := tagbase.UpsertFile(file.GetID(), dirtag)
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
	tagbase := r.Context().Value(types.TAG).(types.Tagbase)
	tagfilter := tag.Tag{
		Word: vals["id"],
		Type: tag.USER,
		Data: tag.Data{
			tag.USER: map[string]string{
				owner.GetID().String(): dirflag,
			},
		},
	}
	filematches, _, err := tagbase.GetFiles([]tag.Tag{tagfilter})
	if err != nil {
		if se := err.(srverror.Error); se.Status() == types.ErrNoResults.Status() {
			w.WriteHeader(se.Status())
		} else {
			panic(srverror.New(err, 500, "Server Error", "unable to get file tags"))
		}
	}

	w.Set("name", strings.ToLower(tagfilter.Word))
	w.Set("files", filematches)
}

func adjustDir(add bool) func(http.ResponseWriter, *http.Request) {
	return func(out http.ResponseWriter, r *http.Request) {
		w := out.(*srvjson.ResponseWriter)

		tagbase := r.Context().Value(types.TAG).(types.Tagbase)
		var owner types.Owner
		if group := r.Context().Value(GROUP); group != nil {
			owner = group.(types.Owner)
		} else {
			owner = r.Context().Value(USER).(types.Owner)
		}
		vals := mux.Vars(r)
		dirtagname := vals["id"]
		if len(dirtagname) == 0 {
			panic(srverror.Basic(400, "dir name missing"))
		}
		fidstrs := r.PostForm["id"]
		if len(fidstrs) == 0 {
			panic(srverror.Basic(400, "missing file ids"))
		}
		var fids []types.FileID
		for _, fidstr := range fidstrs {
			fid, err := types.DecodeFileID(fidstr)
			if err != nil {
				panic(srverror.New(err, 400, "Corrupt File ID"))
			}
			fids = append(fids, fid)
		}
		dirtag := tag.Tag{
			Word: dirtagname,
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					owner.GetID().String(): func() string {
						if add {
							return dirflag
						}
						return ""
					}(),
				},
			},
		}
		for _, fid := range fids {
			err := tagbase.UpsertFile(fid, dirtag)
			if err != nil {
				panic(err)
			}
		}

		w.Set("message", "Complete")
	}
}

func searchDir(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var tagbase = r.Context().Value(types.TAG).(types.Tagbase)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	filters := make([]tag.Tag, 0, 1+len(r.Form["find"]))
	filters = append(filters, tag.Tag{
		Word: vals["id"],
		Type: tag.USER,
		Data: tag.Data{
			tag.USER: map[string]string{
				owner.GetID().String(): dirflag,
			},
		},
	})
	for _, find := range r.Form["find"] {
		filters = append(filters, tag.Tag{
			Word: find,
			Type: tag.CONTENT,
		})
	}
	fids, _, err := tagbase.GetFiles(filters)
	if err != nil {
		panic(err)
	}
	filebase := r.Context().Value(types.FILE).(types.Filebase)
	matches := make([]types.FileID, 0, len(fids))
	for _, fid := range fids {
		file, err := filebase.Get(fid)
		if err != nil && err != types.ErrNotFound {
			panic(err)
		}
		if err == types.ErrNotFound {
			continue
		}
		if file.GetOwner().Match(owner) || file.CheckPerm(owner, "view") {
			matches = append(matches, fid)
		}
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
	fids, _, err := r.Context().Value(types.TAG).(types.Tagbase).GetFiles([]tag.Tag{tag.Tag{
		Word: vals["id"],
		Type: tag.USER,
		Data: tag.Data{
			tag.USER: map[string]string{
				owner.GetID().String(): dirflag,
			},
		},
	}})
	if err != nil {
		panic(err)
	}
	for _, fid := range fids {
		err := r.Context().Value(types.TAG).(types.Tagbase).UpsertFile(fid, tag.Tag{
			Word: vals["id"],
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					owner.GetID().String(): "",
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}
	w.Set("message", "Complete")
}
