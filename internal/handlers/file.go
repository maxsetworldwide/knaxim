package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/internal/decode"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
)

// AttachFile is to add file api paths
func AttachFile(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.Use(groupMiddleware)
	r.HandleFunc("/{id}/download", sendFile).Methods("GET")
	r.HandleFunc("/{id}/view", sendView).Methods("GET")
	{
		r = r.NewRoute().Subrouter()
		r.Use(srvjson.JSONResponse)
		r.HandleFunc("/webpage", webPageUpload).Methods("PUT")
		r.HandleFunc("", createFile).Methods("PUT")
		r.HandleFunc("/{id}", fileInfo).Methods("GET")
		r.HandleFunc("/{id}/slice/{start}/{end}", fileContent).Methods("GET")
		r.HandleFunc("/{id}/search/{start}/{end}", searchFile).Methods("GET")
		r.HandleFunc("/{id}", deleteRecord).Methods("DELETE")
	}

}

var csvextension = regexp.MustCompile("[.](([ct]sv)|(xlsx?))$")

var creationlocks = sync.Map{}

type threadTracker struct {
	L     chan bool
	Count uint
	CL    *sync.Mutex
}

func createFile(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var owner types.Owner
	var maxfiles int64
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
		maxfiles = owner.MaxFiles()
	} else {
		owner = r.Context().Value(USER).(types.Owner)
		if maxfiles = owner.MaxFiles(); maxfiles == 0 {
			maxfiles = config.V.MaxFileCount
		}
	}
	freader, fheader, err := r.FormFile("file")
	if err != nil {
		panic(srverror.New(err, 400, "Error Uploading File"))
	}
	if fheader.Size > config.V.FileLimit {
		panic(srverror.Basic(460, "File exceeds maximum file size"))
	}
	timescale := time.Duration((fheader.Size / 1024) * config.V.FileTimeoutRate)
	if timescale > config.V.MaxFileTimeout.Duration {
		timescale = config.V.MaxFileTimeout.Duration
	}
	if timescale < config.V.MinFileTimeout.Duration {
		timescale = config.V.MinFileTimeout.Duration
	}
	fctx, cancel := context.WithTimeout(context.Background(), timescale)
	defer cancel()
	file := &types.File{
		Permission: types.Permission{
			Own: owner,
		},
		Name: fheader.Filename,
		Date: types.FileTime{Upload: time.Now()},
	}
	nametags, err := tag.BuildNameTags(file.GetName())

	if err != nil {
		panic(srverror.New(err, 400, "Unable to parse filename"))
	}
	var fs *types.FileStore
	// Closure used to ensure lock is garunteed to unlock
	for restart := true; restart; {
		restart = func() bool {
			iTracker, oldTracker := creationlocks.LoadOrStore(owner.GetID().String(), &threadTracker{
				L:     make(chan bool, 1),
				Count: 0,
				CL:    new(sync.Mutex),
			})
			tracker := iTracker.(*threadTracker)
			if !oldTracker {
				tracker.L <- false
			}
			tracker.CL.Lock()
			tracker.Count++
			tracker.CL.Unlock()
			if <-tracker.L {
				tracker.L <- true
				return true
			}
			defer func() {
				tracker.CL.Lock()
				tracker.Count--
				if tracker.Count < 1 {
					creationlocks.Delete(owner.GetID().String())
					tracker.L <- true
				} else {
					tracker.L <- false
				}
				tracker.CL.Unlock()
			}()
			if count, err := r.Context().Value(types.DATABASE).(database.Database).File().Count(owner.GetID()); err != nil {
				panic(err)
			} else if maxfiles > -1 && count >= maxfiles {
				panic(srverror.Basic(461, fmt.Sprintf("Too many files, you can only have %d files. Delete files and empty the trash to make space", maxfiles), fmt.Sprintf("count: %d, maxfiles: %d", count, maxfiles)))
			}
			fs, err = process.InjestFile(fctx, file, fheader.Header.Get("Content-Type"), freader, config.DB)
			if err != nil {
				panic(err)
			}
			return false
		}()
	}
	nameErrCh := make(chan error, 1)
	go func() {
		var filetags []tag.FileTag
		for _, nt := range nametags {
			filetags = append(filetags, tag.FileTag{
				File:  file.GetID(),
				Owner: owner.GetID(),
				Tag:   nt,
			})
		}
		select {
		case nameErrCh <- r.Context().Value(types.TAG).(database.Tagbase).Upsert(filetags...):
		case <-r.Context().Done():
		}
	}()
	if fs.Perr != nil {
		pctx := context.WithValue(context.Background(), decode.TIMEOUT, timescale*5)
		pctx = context.WithValue(pctx, decode.PROCESSING, config.GetResourceTracker())
		go decode.Read(pctx, nil, fs, config.DB, config.T.Path, config.V.GotenPath)
	}
	if len(r.FormValue("dir")) > 0 {
		db, err := config.DB.Connect(fctx)
		if err != nil {
			panic(srverror.New(err, 500, "Failed to add to folder", "Failed to connect to database"))
		}
		defer db.Close(fctx)
		err = db.Tag().Upsert(tag.FileTag{
			File:  file.GetID(),
			Owner: owner.GetID(),
			Tag: tag.Tag{
				Word: r.FormValue("dir"),
				Type: tag.USER,
			},
		})
		if err != nil {
			panic(err)
		}
	}
	select {
	case err := <-nameErrCh:
		if err != nil {
			panic(err)
		}
	case <-r.Context().Done():
	}
	w.Set("id", file.GetID())
	w.Set("name", file.GetName())
}

var getter http.Client

func webPageUpload(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var owner types.Owner
	var maxFiles int64
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
		maxFiles = owner.MaxFiles()
	} else {
		owner = r.Context().Value(USER).(types.Owner)
		if maxFiles = owner.MaxFiles(); maxFiles == 0 {
			maxFiles = config.V.MaxFileCount
		}
	}
	URL, err := url.Parse(r.FormValue("url"))
	if err != nil {
		panic(srverror.New(err, 400, "Bad URL", r.FormValue("url"), "Unable to Parse"))
	}
	nametags, err := tag.BuildNameTags(URL.String())
	if err != nil {
		panic(srverror.New(err, 400, "Bad URL", "Unable to tokenize URL", r.FormValue("url")))
	}

	resp, err := getter.Get(URL.String())
	if err != nil {
		panic(srverror.New(err, 400, "Unable to Get Address", r.FormValue("url"), URL.String()))
	}

	var fs *types.FileStore
	var file types.FileI
	var timescale time.Duration
	var fctx context.Context
	for restart := true; restart; {
		restart = func() bool {
			iTracker, oldTracker := creationlocks.LoadOrStore(owner.GetID().String(), &threadTracker{
				L:     make(chan bool, 1),
				Count: 0,
				CL:    new(sync.Mutex),
			})
			tracker := iTracker.(*threadTracker)
			if !oldTracker {
				tracker.L <- false
			}
			tracker.CL.Lock()
			tracker.Count++
			tracker.CL.Unlock()
			if <-tracker.L {
				tracker.L <- true
				return true
			}
			defer func() {
				tracker.CL.Lock()
				tracker.Count--
				if tracker.Count < 1 {
					creationlocks.Delete(owner.GetID().String())
					tracker.L <- true
				} else {
					tracker.L <- false
				}
				tracker.CL.Unlock()
			}()
			if count, err := r.Context().Value(types.DATABASE).(database.Database).File().Count(owner.GetID()); err != nil {
				panic(err)
			} else if maxFiles > -1 && count >= maxFiles {
				panic(srverror.Basic(461, fmt.Sprintf("Too many files, you can only have %d files. Delete files and empty the trash to make space", maxFiles), fmt.Sprintf("count: %d, maxfiles: %d", count, maxFiles)))
			}
			if process.MapContentType(resp.Header.Get("Content-Type")) == process.URL {

				res, err := process.NewFileConverter(config.V.GotenPath).ConvertURL(URL.String())
				if err != nil {
					panic(srverror.New(err, 400, "Unable to Get Address", "gotenburg", r.FormValue("url"), URL.String()))
				}

				if int64(len(res)) > config.V.FileLimit {
					panic(srverror.Basic(460, "File at URL Exceeds File Limit", r.FormValue("url"), URL.String()))
				}

				timescale = time.Duration(int64(len(res)/1024) * config.V.FileTimeoutRate)
				if timescale > config.V.MaxFileTimeout.Duration {
					timescale = config.V.MaxFileTimeout.Duration
				}
				if timescale < config.V.MinFileTimeout.Duration {
					timescale = config.V.MinFileTimeout.Duration
				}

				var cancel context.CancelFunc
				fctx, cancel = context.WithTimeout(context.Background(), timescale)
				defer cancel()
				file = &types.WebFile{
					File: types.File{
						Permission: types.Permission{
							Own: owner,
						},
						Name: URL.String(),
						Date: types.FileTime{Upload: time.Now()},
					},
					URL: URL.String(),
				}
				fs, err = process.InjestFile(fctx, file, "application/pdf", bytes.NewReader(res), config.DB)
				if err != nil {
					panic(err)
				}
			} else {
				if resp.ContentLength > config.V.FileLimit {
					panic(srverror.Basic(460, "File at URL Exceeds File Limit", r.FormValue("url"), URL.String()))
				}

				timescale = time.Duration((resp.ContentLength / 1024) * config.V.FileTimeoutRate)
				if timescale > config.V.MaxFileTimeout.Duration {
					timescale = config.V.MaxFileTimeout.Duration
				}
				if timescale < config.V.MinFileTimeout.Duration {
					timescale = config.V.MinFileTimeout.Duration
				}

				var cancel context.CancelFunc
				fctx, cancel = context.WithTimeout(context.Background(), timescale)
				defer cancel()
				file = &types.WebFile{
					File: types.File{
						Permission: types.Permission{
							Own: owner,
						},
						Name: URL.String(),
						Date: types.FileTime{Upload: time.Now()},
					},
					URL: URL.String(),
				}
				fs, err = process.InjestFile(fctx, file, resp.Header.Get("Content-Type"), resp.Body, config.DB)
				if err != nil {
					panic(err)
				}
			}
			return false
		}()
	}
	nameErrCh := make(chan error, 1)
	go func() {
		var filetags []tag.FileTag
		for _, nt := range nametags {
			filetags = append(filetags, tag.FileTag{
				File:  file.GetID(),
				Owner: owner.GetID(),
				Tag:   nt,
			})
		}
		select {
		case nameErrCh <- r.Context().Value(types.TAG).(database.Tagbase).Upsert(filetags...):
		case <-r.Context().Done():
		}
	}()
	if fs.Perr != nil {
		pctx := context.WithValue(context.Background(), decode.TIMEOUT, timescale*5)
		pctx = context.WithValue(pctx, decode.PROCESSING, config.GetResourceTracker())
		go decode.Read(pctx, nil, fs, config.DB, config.T.Path, config.V.GotenPath)
	}
	if len(r.FormValue("dir")) > 0 {
		db, err := config.DB.Connect(fctx)
		if err != nil {
			panic(srverror.New(err, 500, "Failed to add to folder", "Failed to connect to database"))
		}
		defer db.Close(fctx)
		err = db.Tag().Upsert(tag.FileTag{
			File:  file.GetID(),
			Owner: owner.GetID(),
			Tag: tag.Tag{
				Word: r.FormValue("dir"),
				Type: tag.USER,
			},
		})
		if err != nil {
			panic(err)
		}
	}
	select {
	case err := <-nameErrCh:
		if err != nil {
			panic(err)
		}
	case <-r.Context().Done():
	}
	w.Set("id", file.GetID())
	w.Set("name", file.GetName())
}

func fileInfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(err)
	}
	frec, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !frec.GetOwner().Match(owner) && !frec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", owner.GetID().String(), frec.GetName(), frec.GetID().String()))
	}
	count, err := r.Context().Value(types.CONTENT).(database.Contentbase).Len(frec.GetID().StoreID)
	if err != nil {
		panic(err)
	}
	store, err := r.Context().Value(types.STORE).(database.Storebase).Get(frec.GetID().StoreID)
	if err != nil {
		panic(err)
	}

	finfo := FileInfo{frec, count, store.FileSize}
	w.Set("file", finfo.File)
	w.Set("count", finfo.Count)
	w.Set("size", finfo.Size)
}

func fileContent(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	startindx, err := strconv.Atoi(vals["start"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, fileContent start value not a number"))
	}
	endindx, err := strconv.Atoi(vals["end"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, fileContent end value not a number"))
	}
	if startindx > endindx {
		panic(srverror.Basic(400, "Bad Request, end must be greater then start"))
	}
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad File ID"))
	}
	rec, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) && !rec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "fileContent user no view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}

	lines, err := r.Context().Value(types.CONTENT).(database.Contentbase).Slice(rec.GetID().StoreID, startindx, endindx)
	if pe, ok := err.(*errors.Processing); ok {
		w.WriteHeader(pe.Status)
		w.Set("ProcessingError", pe.Message)
	} else if err != nil && lines == nil {
		panic(err)
	}

	w.Set("size", len(lines))
	w.Set("lines", lines)
}

func searchFile(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	regex := util.BuildSearchRegex(r.FormValue("find"))
	start, err := strconv.Atoi(vals["start"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, searchFile start not a number"))
	}
	end, err := strconv.Atoi(vals["end"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, searchFile end not a number"))
	}
	if start > end {
		panic(srverror.Basic(400, "Bad Request, end must be greater then start"))
	}
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, bad file id"))
	}
	file, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !file.GetOwner().Match(owner) && !file.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "user does not have view permission", owner.GetID().String(), file.GetName(), file.GetID().String()))
	}
	matched, err := r.Context().Value(types.CONTENT).(database.Contentbase).RegexSearchFile(regex, file.GetID().StoreID, start, end)
	if err != nil {
		if pe, ok := err.(*errors.Processing); ok {
			w.WriteHeader(pe.Status)
			w.Set("ProcessingError", pe.Message)
		} else if matched == nil {
			panic(err)
		}
	}
	w.Set("size", len(matched))
	w.Set("lines", matched)
}

func deleteRecord(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, malformed file id"))
	}
	rec, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) {
		panic(srverror.Basic(403, "Permission Denied", "deleteRecord user not owner", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	if err = r.Context().Value(types.FILE).(database.Filebase).Remove(fid); err != nil {
		panic(err)
	}
	w.Set("message", "File Removed")
	// w.Write([]byte("File Removed"))
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	if jsw, ok := w.(*srvjson.ResponseWriter); ok {
		w = jsw.Internal
	}
	var owner types.Owner
	// shouldn't need to check for group
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, bad file id"))
	}
	rec, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) && !rec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "sendFile user not have view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	store, err := r.Context().Value(types.STORE).(database.Storebase).Get(fid.StoreID)
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

func sendView(w http.ResponseWriter, r *http.Request) {
	var owner types.Owner
	// shouldn't need to check for group
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	vals := mux.Vars(r)
	fid, err := types.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, bad file id"))
	}
	rec, err := r.Context().Value(types.FILE).(database.Filebase).Get(fid)
	if err != nil {
		panic(err)
	}
	if !rec.GetOwner().Match(owner) && !rec.CheckPerm(owner, "view") {
		panic(srverror.Basic(403, "Permission Denied", "sendView user does not have view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	fs, err := r.Context().Value(types.STORE).(database.Storebase).Get(fid.StoreID)
	if err != nil {
		panic(err)
	}
	ext := fs.ContentType
	var rdr io.Reader
	if process.ExtMap[ext] == process.PDF {
		rdr, err = fs.Reader()
		if err != nil {
			panic(err)
		}
	} else {
		view, err := r.Context().Value(types.VIEW).(database.Viewbase).Get(fid.StoreID)
		if err != nil {
			if fs.Perr != nil {
				panic(srverror.Basic(fs.Perr.Status, fs.Perr.Message))
			} else {
				panic(srverror.Basic(303, "No View, use sentence view"))
			}
		}
		rdr, err = view.Reader()
		if err != nil {
			panic(err)
		}
	}
	pdfName := rec.GetName()
	dotIdx := strings.LastIndex(rec.GetName(), ".")
	if dotIdx > -1 {
		pdfName = rec.GetName()[:dotIdx+1] + "pdf"
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+pdfName+"\"")
	w.Header().Set("Content-Type", "application/pdf")
	io.Copy(w, rdr)
}
