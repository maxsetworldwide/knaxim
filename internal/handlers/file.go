package handlers

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"

	"github.com/gorilla/mux"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

func AttachFile(r *mux.Router) {
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.Use(groupMiddleware)
	r.HandleFunc("/webpage", webPageUpload).Methods("PUT")
	r.HandleFunc("", createFile).Methods("PUT")
	//r.HandleFunc("/copy", copyFile).Methods("PUT")
	r.HandleFunc("/{id}", fileInfo).Methods("GET")
	r.HandleFunc("/{id}/slice/{start}/{end}", fileContent).Methods("GET")
	r.HandleFunc("/{id}/search/{start}/{end}", searchFile).Methods("GET")
	r.HandleFunc("/{id}/download", sendFile).Methods("GET")
	r.HandleFunc("/{id}/view", sendView).Methods("GET")
	//r.HandleFunc("/{id}/refresh", refreshWebPage).Methods("POST")
	r.HandleFunc("/{id}", deleteRecord).Methods("DELETE")
}

var csvextension = regexp.MustCompile("[.](([ct]sv)|(xlsx?))$")

func processContent(ctx context.Context, cancel context.CancelFunc, file database.FileI, fs *database.FileStore) error {
	if cancel != nil {
		defer cancel()
	}
	gotenbergErr := make(chan error, 1)
	go createView(ctx, config.)
	rcontent, err := fs.Reader()
	if err != nil {
		return err
	}
	tikapath := config.T.Path
	contentex := process.NewContentExtractor(nil, tikapath)
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
	err = config.DB.Content(ctx).Insert(contentlines...)
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
	return config.DB.Tag(ctx).UpsertStore(fs.ID, tags...)
}

func createView(ctx context.Context, db database.Database, file database.FileI, fs *database.FileStore, out chan error) {
	name := file.GetName()
	buf := &bytes.Buffer{}
	r, err := fs.Reader()
	if err != nil {
		out <- err
		return
	}
	if _, err = io.Copy(buf, r); err != nil {
		out <- err
		return
	}
	url := config.V.GotenPath
	converter := process.NewFileConverter(url)
	var result *bytes.Buffer
	gotenFinished := make(chan error)
	go func() {
		var err error
		result, err = converter.ConvertOffice(name, buf)
		gotenFinished <- err
	}()
	select{
	case err := <- gotenFinished:
		if err != nil {
			out <- err
			return
		}
	case <- ctx.Done():
		out <- ctx.Err()
		return
	}
	vb := db.View(nil)
	vs, err := database.NewViewStore(fs.ID, result)
	if err != nil {
		out <- err
		return
	}
	if err = vb.Insert(vs); err != nil {
		out <- err
		return
	}
	out <- nil
}

func createFile(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

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
	if timescale > config.V.MaxFileTimeout.Duration {
		timescale = config.V.MaxFileTimeout.Duration
	}
	if timescale < config.V.MinFileTimeout.Duration {
		timescale = config.V.MinFileTimeout.Duration
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
	fs, err := process.InjestFile(fctx, file, fheader.Header.Get("Content-Type"), freader, config.DB)
	if err != nil {
		panic(err)
	}
	pctx, cncl := context.WithTimeout(context.Background(), timescale*5)

	go func() {
		if err := processContent(pctx, cncl, file, fs); err != nil {
			util.VerboseRequest(r, "Processing Error: %s", err.Error())
			fs.Perr = &database.ProcessingError{
				Status:  242,
				Message: err.Error(),
			}
		} else {
			fs.Perr = nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		sb := config.DB.Store(ctx)
		err := sb.UpdateMeta(fs)
		if err != nil {
			util.VerboseRequest(r, "Unable to Update Processing Error: %s", err.Error())
		}
		sb.Close(ctx)
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

	w.Set("id", file.GetID())
	w.Set("name", file.GetName())
}

var getter http.Client

func webPageUpload(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
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
		if timescale > config.V.MaxFileTimeout.Duration {
			timescale = config.V.MaxFileTimeout.Duration
		}
		if timescale < config.V.MinFileTimeout.Duration {
			timescale = config.V.MinFileTimeout.Duration
		}

	} else {
		timescale = config.V.MaxFileTimeout.Duration
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
	fs, err := process.InjestFile(fctx, file, res.Header.Get("Content-Type"), res.Body, config.DB)
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

	w.Set("id", file.GetID())
	w.Set("name", file.GetName())
}

func fileInfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

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
	w.Set("file", finfo.File)
	w.Set("count", finfo.Count)
	w.Set("size", finfo.Size)
}

func fileContent(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
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
	if pe, ok := err.(*database.ProcessingError); ok {
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
	matched, err := r.Context().Value(database.CONTENT).(database.Contentbase).RegexSearchFile(regex, file.GetID().StoreID, start, end)
	if err != nil {
		if pe, ok := err.(*database.ProcessingError); ok {
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
	w.Set("message", "File Removed")
	// w.Write([]byte("File Removed"))
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	if jsw, ok := w.(*srvjson.ResponseWriter); ok {
		w = jsw.Internal
	}
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

func sendView(w http.ResponseWriter, r *http.Request) {
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
		panic(srverror.Basic(403, "Permission Denied", "sendView user not have view permission", owner.GetID().String(), rec.GetName(), rec.GetID().String()))
	}
	store, err := r.Context().Value(database.STORE).(database.Storebase).Get(fid.StoreID)
	if err != nil {
		panic(err)
	}
	rdr, err := store.Reader()
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(rdr)
	if err != nil {
		panic(err)
	}

	gotenClient := &gotenberg.Client{
		Hostname: "http://gotenberg:3000",
	}

	// capital extensions
	// landscape spreadsheets
	// deadline
	lower := strings.ToLower(rec.GetName())
	index, err := gotenberg.NewDocumentFromBytes(lower, bytes)

	req := gotenberg.NewOfficeRequest(index)
	res, err := gotenClient.Post(req)
	if err != nil {
		panic(err)
	}
	rdr = res.Body
	w.Header().Set("Content-Disposition", "attachment; filename=\""+rec.GetName()+"\"")
	w.Header().Set("Content-Type", store.ContentType)
	io.Copy(w, rdr)
}
