package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	dbug "runtime/debug"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

type ResWrtrCapturer struct {
	internal   http.ResponseWriter
	StatusCode int
}

func (rwc *ResWrtrCapturer) Header() http.Header         { return rwc.internal.Header() }
func (rwc *ResWrtrCapturer) Write(b []byte) (int, error) { return rwc.internal.Write(b) }
func (rwc *ResWrtrCapturer) WriteHeader(sc int) {
	rwc.StatusCode = sc
	rwc.internal.WriteHeader(sc)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := new(ResWrtrCapturer)
		nw.StatusCode = 200
		nw.internal = w
		verboseRequest(r, "Recieved: %+v", r.Header)
		next.ServeHTTP(nw, r)
		verboseRequest(r, "Complete(%d): %+v", nw.StatusCode, nw.Header())
	})
}

var debugflag = flag.Bool("debug", false, "write debug messages to response Writer")

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if se, ok := err.(srverror.Error); ok {
					se.ServeHTTP(w, r)
					verboseRequest(r, se.Error())
				} else {
					verbose("Non-standard panic")
					w.WriteHeader(500)
					w.Write([]byte("{\"message\": \"Server Error\"}"))
					w.Header().Set("Content-Type", "application/json")
					switch v := err.(type) {
					case string:
						verboseRequest(r, "Panic: %s", v)
					case error:
						verboseRequest(r, "Panic: %s", v.Error())
					default:
						verboseRequest(r, "Panic: %v", v)
					}
					if *debugflag {
						dbug.PrintStack()
					}
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//Token verification MiddleWare
func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userbase := r.Context().Value(database.OWNER).(database.Ownerbase)
		uid, err := database.GetCookieUid(r)
		if err != nil {
			panic(srverror.New(err, 401, "login", "invalid cookie, error getting userid from cookie"))
		}
		user, err := userbase.Get(uid)
		if err != nil {
			panic(srverror.New(err, 401, "login", "unable to get user record to validate token", uid.String()))
		}
		if u, ok := user.(database.UserI); !ok {
			panic(srverror.New(errors.New("id is not a user"), 401, "login", uid.String()))
		} else if u.GetRole("Guest") && r.Method != "GET" {
			panic(srverror.New(errors.New("Guest User cannot perform action"), 401, "login", "Invalid Guest Action", r.Method, r.URL.Path))
		} else if !u.GetRole("Guest") && !u.CheckCookie(r) {
			panic(srverror.New(errors.New("Cookie not valid"), 401, "login", "Cookie Invalid", uid.String()))
		}
		r = r.WithContext(context.WithValue(r.Context(), USER, user))
		next.ServeHTTP(w, r)
	})
}

func ConnectDatabase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		r = r.WithContext(ctx)

		filebase := db.File(r.Context())
		defer filebase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.FILE, filebase))

		ownerbase := filebase.Owner(nil)
		defer ownerbase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.OWNER, ownerbase))

		storebase := filebase.Store(nil)
		defer storebase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.STORE, storebase))

		contentbase := filebase.Content(nil)
		defer contentbase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.CONTENT, contentbase))

		tagbase := filebase.Tag(nil)
		defer tagbase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.TAG, tagbase))

		acronymbase := filebase.Acronym(nil)
		defer acronymbase.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), database.ACRONYM, acronymbase))

		next.ServeHTTP(w, r)
	})
}

func Timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if conf.BasicTimeout > 0 {
			c, cancel := context.WithTimeout(r.Context(), conf.BasicTimeout)
			defer cancel()
			r = r.WithContext(c)
		}
		next.ServeHTTP(w, r)
	})
}

type jsonform map[string]interface{}

func (jf jsonform) Values() url.Values {
	out := make(url.Values)
	var empty = true
	for k, v := range jf {
		switch tostr := v.(type) {
		case string:
			out.Add(k, tostr)
		case bool:
			if tostr {
				out.Add(k, "true")
			} else {
				out.Add(k, "false")
			}
		case float64:
			out.Add(k, fmt.Sprintf("%f", tostr))
		case []interface{}:
			for _, ele := range tostr {
				if elestr, ok := ele.(string); ok {
					out.Add(k, elestr)
				} else {
					jsonbytes, _ := json.Marshal(ele)
					out.Add(k, string(jsonbytes))
				}
			}
		case map[string]interface{}:
			for subkey, val := range tostr {
				out.Add(k, subkey)
				jsonbytes, _ := json.Marshal(val)
				out.Add(k, string(jsonbytes))
			}
		}
		empty = false
	}
	if empty {
		return nil
	}
	return out
}

func ParseBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var jf jsonform
			if err := json.NewDecoder(r.Body).Decode(&jf); err != nil {
				panic(srverror.New(err, 400, "Unable to decode json object"))
			}
			r.PostForm = jf.Values()
		}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			panic(srverror.New(err, 400, "Unable to parse url form values"))
		}
		next.ServeHTTP(w, r)
	})
}
