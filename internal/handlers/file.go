package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"github.com/gorilla/mux"
)

func AttachFile(r *mux.Router) {
	r.Use(UserCookie)
	r.Use(groupMiddleware)
	r.HandleFunc("/webpage", webPageUpload).Methods("PUT")
	r.HandleFunc("", createFile).Methods("PUT")
	//r.HandleFunc("/copy", copyFile).Methods("PUT")
	r.HandleFunc("/{id}", fileInfo).Methods("GET")
	r.HandleFunc("/{id}/slice/{start}/{end}", fileContent).Methods("GET")
	r.HandleFunc("/{id}/search/{start}/{end}", searchFile).Methods("GET")
	r.HandleFunc("/{id}/download", sendFile).Methods("GET")
	//r.HandleFunc("/{id}/refresh", refreshWebPage).Methods("POST")
	r.HandleFunc("/{id}", deleteRecord).Methods("DELETE")
}

var csvextension = regexp.MustCompile("[.](([ct]sv)|(xlsx?))$")

func processContent(ctx context.Context, cancel context.CancelFunc, file database.FileI, fs *database.FileStore) error {
	if cancel != nil {
		defer cancel()
	}
	rcontent, err := fs.Reader()
	if err != nil {
		return err
	}
	tikapath := config.T.Path
	contentex := database.NewContentExtractor(nil, tikapath)
	var contentlines []database.ContentLine
	if csvextension.MatchString(file.GetName()) {
		contentlines, err = contentex.ExtractCSV(ctx, rcontent)
		if err != nil {
			return err
		}
	} else {
		contentlines, err = contentex.ExtractText(ctx, rcontent)
		if err != nil {
			return err
		}
	}
	for i := range contentlines {
		contentlines[i].ID = fs.ID
	}
	// util.Verbose("generated content: %v", contentlines)
	{
		cnt := config.DB.Content(ctx)
		defer cnt.Close(ctx)
		err = cnt.Insert(contentlines...)
	}
	if err != nil {
		return err
	}
	creader, err := database.NewContentReader(contentlines)
	if err != nil {
		return err
	}
	tags, err := tag.ExtractContentTags(creader)
	if err != nil {
		return err
	}
	tg := config.DB.Tag(ctx)
	defer tg.Close(ctx)
	return tg.UpsertStore(fs.ID, tags...)
}

func createFile(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	freader, fheader, err := r.FormFile("file")
	if err != nil {
		panic(srverror.New(err, 400, "Error Uploading File"))
	}
	if fheader.Size > config.V.FileLimit {
		panic(srverror.Basic(460, "File exceeds maximum file size"))
	}
	timescale := time.Duration((fheader.Size / 1024) * config.V.FileTimeoutRate)
	if timescale > config.V.MaxFileTimeout {
		timescale = config.V.MaxFileTimeout
	}
	if timescale < config.V.MinFileTimeout {
		timescale = config.V.MinFileTimeout
	}
	fctx, cancel := context.WithTimeout(context.Background(), timescale)
	defer cancel()
	file := &database.File{
		Permission: database.Permission{
			Own: owner,
		},
		Name: fheader.Filename,
		Date: database.FileTime{Upload: time.Now()},
	}
	fs, err := database.InjestFile(fctx, file, fheader.Header.Get("Content-Type"), freader, config.DB)
	if err != nil {
		panic(err)
	}
	pctx, cncl := context.WithTimeout(context.Background(), timescale*5)
	go func() {
		if err := processContent(pctx, cncl, file, fs); err != nil {
			util.VerboseRequest(r, "Processing Error: %s", err.Error())
		}
	}()
	if len(r.FormValue("dir")) > 0 {
		err = config.DB.Tag(fctx).UpsertFile(file.GetID(), tag.Tag{
			Word: r.FormValue("dir"),
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					owner.GetID().String(): dirflag,
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   file.GetID(),
		"name": file.GetName(),
	}); err != nil {
		panic(srverror.New(err, 500, "Server Error", "create file unable to encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

var getter http.Client

func webPageUpload(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	URL, err := url.Parse(r.FormValue("url"))
	if err != nil {
		panic(srverror.New(err, 400, "Bad URL", r.FormValue("url"), "Unable to Parse"))
	}
	res, err := getter.Get(URL.String())
	if err != nil {
		panic(srverror.New(err, 400, "Unable to Get Address", r.FormValue("url"), URL.String()))
	}
	if res.ContentLength > -1 && res.ContentLength > config.V.FileLimit {
		panic(srverror.Basic(460, "File at URL Exceeds File Limit", r.FormValue("url"), URL.String()))
	}

	var timescale time.Duration
	if res.ContentLength > -1 {
		timescale = time.Duration((res.ContentLength / 1024) * config.V.FileTimeoutRate)
		if timescale > config.V.MaxFileTimeout {
			timescale = config.V.MaxFileTimeout
		}
		if timescale < config.V.MinFileTimeout {
			timescale = config.V.MinFileTimeout
		}

	} else {
		timescale = config.V.MaxFileTimeout
	}
	fctx, cancel := context.WithTimeout(context.Background(), timescale)
	defer cancel()
	file := &database.WebFile{
		File: database.File{
			Permission: database.Permission{
				Own: owner,
			},
			Name: URL.String(),
			Date: database.FileTime{Upload: time.Now()},
		},
		URL: URL.String(),
	}
	fs, err := database.InjestFile(fctx, file, res.Header.Get("Content-Type"), res.Body, config.DB)
	if err != nil {
		panic(err)
	}
	pctx, cncl := context.WithTimeout(context.Background(), timescale*5)
	go processContent(pctx, cncl, file, fs)
	if len(r.FormValue("dir")) > 0 {
		err = config.DB.Tag(fctx).UpsertFile(file.GetID(), tag.Tag{
			Word: r.FormValue("dir"),
			Type: tag.USER,
			Data: tag.Data{
				tag.USER: map[string]string{
					owner.GetID().String(): dirflag,
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   file.GetID(),
		"name": file.GetName(),
	}); err != nil {
		panic(srverror.New(err, 500, "Server Error", "create file unable to encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func fileInfo(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(err)
	}
	frec, err := r.Context().Value(database.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !frec.GetOwner().Match(owner) && !frec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", owner.GetID().String(), frec.GetName(), frec.GetID().String()))
	}
	count, err := r.Context().Value(database.CONTENT).(database.Contentbase).Len(frec.GetID().StoreID)
	if err != nil {
		panic(err)
	}
	store, err := r.Context().Value(database.STORE).(database.Storebase).Get(frec.GetID().StoreID)
	if err != nil {
		panic(err)
	}
	finfo := FileInfo{frec, count, store.FileSize}
	if err = json.NewEncoder(w).Encode(finfo); err != nil {
		panic(srverror.New(err, 500, "Server Error", "fileInfo encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func fileContent(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	startindx, err := strconv.Atoi(vals["start"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "fileContent start value not a number"))
	}
	endindx, err := strconv.Atoi(vals["end"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "fileContent end value not a number"))
	}
	if startindx > endindx {
		panic(srverror.Basic(400, "Bad Request", "end must be greater then start"))
	}
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request"))
	}
	rec, err := r.Context().Value(database.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) && !rec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "fileContent user no view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}

	lines, err := r.Context().Value(database.CONTENT).(database.Contentbase).Slice(rec.GetID().StoreID, startindx, endindx)
	if err != nil {
		panic(err)
	}
	var result FileContent
	result.Length = len(lines)
	result.Vals = lines
	if err = json.NewEncoder(w).Encode(result); err != nil {
		panic(srverror.New(err, 500, "Server Error", "fileContent encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func searchFile(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	regex := util.BuildSearchRegex(r.FormValue("find"))
	start, err := strconv.Atoi(vals["start"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "searchFile start not a number"))
	}
	end, err := strconv.Atoi(vals["end"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "searchFile end not a number"))
	}
	if start > end {
		panic(srverror.Basic(400, "Bad Request", "end must be greater then start"))
	}
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "bad file id"))
	}
	file, err := r.Context().Value(database.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !file.GetOwner().Match(owner) && !file.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "user does not have view permission", owner.GetID().String(), file.GetName(), file.GetID().String()))
	}
	if matched, err := r.Context().Value(database.CONTENT).(database.Contentbase).RegexSearchFile(regex, file.GetID().StoreID, start, end); err != nil {
		panic(err)
	} else {
		var out FileContent
		out.Length = len(matched)
		out.Vals = matched
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(srverror.New(err, 500, "Server Error", "searchFile encode json"))
		}
		w.Header().Add("Content-Type", "application/json")
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "malformed file id"))
	}
	rec, err := r.Context().Value(database.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) {
		panic(srverror.Basic(403, "Permission Denied", "deleteRecord user not owner", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	if err = r.Context().Value(database.FILE).(database.Filebase).Remove(fid); err != nil {
		panic(err)
	}
	w.Write([]byte("File Removed"))
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	vals := mux.Vars(r)
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request", "bad file id"))
	}
	rec, err := r.Context().Value(database.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) && !rec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "sendFile user not have view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	store, err := r.Context().Value(database.STORE).(database.Storebase).Get(fid.StoreID)
	if err != nil {
		panic(err)
	}
	rdr, err := store.Reader()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+rec.GetName()+"\"")
	w.Header().Set("Content-Type", store.ContentType)
	io.Copy(w, rdr)
}
