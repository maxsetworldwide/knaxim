package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

func AttachDir(r *mux.Router) {
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
	var owner database.Owner
	if gr := r.Context().Value(GROUP); gr != nil {
		owner = gr.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}

	var dirs []string
	if tags, err := r.Context().Value(database.TAG).(database.Tagbase).SearchData(tag.USER, tag.Data{tag.USER: map[string]string{owner.GetID().String(): dirflag}}); err == nil {
		for _, t := range tags {
			dirs = append(dirs, t.Word)
		}
	} else {
		panic(err)
	}

	w.Set("folders", dirs)
	//if err := json.NewEncoder(w).Encode(map[string]interface{}{
	//	"folders": dirs,
	//}); err != nil {
	//	panic(srverror.New(err, 500, "Server Error", "Failed to encode json"))
	//}
	//w.Header().Set("Content-Type", "application/json")
}

func createDir(w http.ResponseWriter, r *http.Request) {
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	nname := r.FormValue("newname")
	if !validDirName(nname) {
		panic(srverror.Basic(400, "Invalid Directory Name"))
	}
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	var files []database.FileI
	if r.Form["content"] != nil {
		for _, fidstr := range r.Form["content"] {
			fid, err := filehash.DecodeFileID(fidstr)
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
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":            nname,
			"affectedFiles": 0,
		}); err != nil {
			panic(srverror.New(err, 500, "Server Error", "createDir failed to encode json 1"))
		}
		w.Header().Add("Content-Type", "application/json")
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
	tagbase := r.Context().Value(database.TAG).(database.Tagbase)
	for _, file := range files {
		err := tagbase.UpsertFile(file.GetID(), dirtag)
		if err != nil {
			panic(srverror.New(err, 500, "Server Error", "Unable to add user tag"))
		}
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            nname,
		"affectedFiles": len(files),
	}); err != nil {
		panic(srverror.New(err, 500, "Server Error", "createDir failed to encode json 2"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func dirInfo(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	tagbase := r.Context().Value(database.TAG).(database.Tagbase)
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
		panic(srverror.New(err, 500, "Server Error", "unable to get file tags"))
	}
	result := DirInformation{
		Name:  strings.ToLower(tagfilter.Word),
		Files: filematches,
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(srverror.New(err, 500, "Server Error", "dirInfo unable to encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func adjustDir(add bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tagbase := r.Context().Value(database.TAG).(database.Tagbase)
		var owner database.Owner
		if group := r.Context().Value(GROUP); group != nil {
			owner = group.(database.Owner)
		} else {
			owner = r.Context().Value(USER).(database.Owner)
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
		var fids []filehash.FileID
		for _, fidstr := range fidstrs {
			fid, err := filehash.DecodeFileID(fidstr)
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
		w.Write([]byte("Complete"))
	}
}

func searchDir(w http.ResponseWriter, r *http.Request) {
	var tagbase = r.Context().Value(database.TAG).(database.Tagbase)
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
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
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	matches := make([]filehash.FileID, 0, len(fids))
	for _, fid := range fids {
		file, err := filebase.Get(fid)
		if err != nil && err != database.ErrNotFound {
			panic(err)
		}
		if err == database.ErrNotFound {
			continue
		}
		if file.GetOwner().Match(owner) || file.CheckPerm(owner, "view") {
			matches = append(matches, fid)
		}
	}
	if err := json.NewEncoder(w).Encode(BuildSearchResponse(r, matches)); err != nil {
		panic(srverror.New(err, 500, "Server Error", "searchDir failed to encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func deleteDir(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	fids, _, err := r.Context().Value(database.TAG).(database.Tagbase).GetFiles([]tag.Tag{tag.Tag{
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
		err := r.Context().Value(database.TAG).(database.Tagbase).UpsertFile(fid, tag.Tag{
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
	w.Write([]byte("Complete"))
}
