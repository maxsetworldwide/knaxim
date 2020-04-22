package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

func AttachNLP(r *mux.Router) {
	r.Use(ConnectDatabase)
	r.Use(ParseBody)
	r.Use(UserCookie)
	r.Use(srvjson.JSONResponse)
	r.HandleFunc("/file/{fid}/{synth}/{start}/{end}", sendNLP).Methods("GET")
}

func sendNLP(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)

	vals := mux.Vars(r)
	start, err := strconv.Atoi(vals["start"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, expected to provide range as numbers"))
	}
	end, err := strconv.Atoi(vals["end"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, expected to provide range as numbers"))
	}
	if start >= end {
		panic(srverror.Basic(400, "Bad Request, initial index must be less then final index"))
	}
	if end > 50 {
		panic(srverror.Basic(400, "Bad Request, final index too large"))
	}
	var tagtype tag.Type
	switch strings.ToLower(vals["synth"]) {
	case "t":
		fallthrough
	case "topic":
		tagtype = tag.TOPIC
	case "a":
		fallthrough
	case "action":
		tagtype = tag.ACTION
	case "r":
		fallthrough
	case "resource":
		tagtype = tag.RESOURCE
	case "p":
		fallthrough
	case "process":
		tagtype = tag.PROCESS
	default:
		panic(srverror.Basic(400, "Bad Request, unrecognized nlp category"))
	}
	fid, err := types.DecodeFileID(vals["fid"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad Request, Bad file id"))
	}
	fb := r.Context().Value(types.FILE).(database.Filebase)
	file, err := fb.Get(fid)
	if err != nil {
		panic(err)
	}
	user := r.Context().Value(USER).(types.Owner)
	if !file.CheckPerm(user, "view") && !file.GetOwner().Match(user) {
		panic(srverror.Basic(403, "Access Denied"))
	}
	sb := r.Context().Value(types.STORE).(database.Storebase)
	fs, err := sb.Get(fid.StoreID)
	if err != nil {
		panic(err)
	}
	if fs.Perr != nil {
		panic(srverror.Basic(fs.Perr.Status, fs.Perr.Message))
	}
	tb := r.Context().Value(types.TAG).(database.Tagbase)
	tags, err := tb.GetType(fid, user.GetID(), tagtype)
	if err != nil {
		panic(err)
	}
	if len(tags) < end {
		end = len(tags)
		if start >= end {
			return
		}
	}
	result := make([]struct {
		Word  string
		Count int
	}, end-start)
	for _, t := range tags {
		position := t.Data[tagtype]["significance"].(int)
		if position >= start && position < end {
			result[position-start].Word = t.Word
			result[position-start].Count = t.Data[tagtype]["count"].(int)
		}
	}
	w.Set("fid", fid)
	w.Set("info", result)
}
