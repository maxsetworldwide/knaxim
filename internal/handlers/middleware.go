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
	"time"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/email"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// ResWrtrCapturer http.ResponseWriter that saves status code
type ResWrtrCapturer struct {
	internal   http.ResponseWriter
	StatusCode int
}

// Header implements http.ResponseWriter
func (rwc *ResWrtrCapturer) Header() http.Header { return rwc.internal.Header() }

// Write implements http.ResponseWriter
func (rwc *ResWrtrCapturer) Write(b []byte) (int, error) { return rwc.internal.Write(b) }

// WriteHeader implements http.ResponseWriter
func (rwc *ResWrtrCapturer) WriteHeader(sc int) {
	rwc.StatusCode = sc
	rwc.internal.WriteHeader(sc)
}

// Logging is a middleware to log requests and responses
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := new(ResWrtrCapturer)
		nw.StatusCode = 200
		nw.internal = w
		util.VerboseRequest(r, "Recieved: %+v", r.Header)
		next.ServeHTTP(nw, r)
		util.VerboseRequest(r, "Complete(%d): %+v", nw.StatusCode, nw.Header())
	})
}

var debugflag = flag.Bool("debug", false, "write debug messages to response Writer")

// Recovery is a middleware to recover from panics in handling requests
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var logMsg string
				if se, ok := err.(srverror.Error); ok {
					se.ServeHTTP(w, r)
					util.VerboseRequest(r, se.Error())
					if se.Status() == 500 {
						logMsg = srverror.LogString(se, r, w)
					}
				} else {
					util.Verbose("Non-standard panic")
					w.WriteHeader(500)
					w.Write([]byte("{\"message\": \"Server Error\"}"))
					w.Header().Set("Content-Type", "application/json")
					switch v := err.(type) {
					case string:
						util.VerboseRequest(r, "Panic: %s", v)
						logMsg = srverror.LogString(errors.New(v), r, w)
					case error:
						util.VerboseRequest(r, "Panic: %s", v.Error())
						logMsg = srverror.LogString(v, r, w)
					default:
						util.VerboseRequest(r, "Panic: %v", v)
						logMsg = srverror.LogString(fmt.Errorf("%v", v), r, w)
					}
					if *debugflag {
						dbug.PrintStack()
					}
				}
				if len(logMsg) > 0 {
					logErr := srverror.WriteToFile(logMsg)
					if logErr != nil {
						util.Verbose("Failed to write log to file: %s", logErr.Error())
					}
					logErr = email.SendErrorEmail(logMsg)
					if logErr != nil {
						util.Verbose("Failed to send log email: %s", logErr.Error())
					} else {
						util.Verbose("An email has been sent to %s describing this issue", config.V.ErrorEmail)
					}
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// UserCookie Token verification MiddleWare
func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userbase := r.Context().Value(types.OWNER).(database.Ownerbase)
		uid, err := types.GetCookieUID(r)
		if err != nil {
			panic(srverror.New(err, 401, "Please log in", "invalid cookie, error getting userid from cookie"))
		}
		user, err := userbase.Get(uid)
		if err != nil {
			panic(srverror.New(err, 401, "Please log in", "unable to get user record to validate token", uid.String()))
		}
		if u, ok := user.(types.UserI); !ok {
			panic(srverror.New(errors.New("id is not a user"), 401, "Please log in", uid.String()))
		} else if u.GetRole("Guest") && r.Method != "GET" {
			panic(srverror.New(errors.New("Guest User cannot perform action"), 401, "Please log in", "Invalid Guest Action", r.Method, r.URL.Path))
		} else if !u.GetRole("Guest") && !u.CheckCookie(r) {
			panic(srverror.New(errors.New("Cookie not valid"), 401, "Please log in", "Cookie Invalid", uid.String()))
		} else {
			u.RefreshCookie(time.Now().Add(config.V.UserTimeouts.Inactivity.Duration))
			if err := userbase.Update(u); err != nil {
				panic(err)
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), USER, user))
		next.ServeHTTP(w, r)
	})
}

// ConnectDatabase is a middleware the opens a connection to the database and populates the request context with connection objects
func ConnectDatabase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		r = r.WithContext(ctx)

		dbConnection, err := config.DB.Connect(r.Context())
		if err != nil {
			panic(srverror.New(err, 500, "Error M1", "Failed to Connect to Database"))
		}
		defer dbConnection.Close(r.Context())
		r = r.WithContext(context.WithValue(r.Context(), types.DATABASE, dbConnection))

		filebase := dbConnection.File()
		r = r.WithContext(context.WithValue(r.Context(), types.FILE, filebase))

		ownerbase := dbConnection.Owner()
		r = r.WithContext(context.WithValue(r.Context(), types.OWNER, ownerbase))

		storebase := dbConnection.Store()
		r = r.WithContext(context.WithValue(r.Context(), types.STORE, storebase))

		contentbase := dbConnection.Content()
		r = r.WithContext(context.WithValue(r.Context(), types.CONTENT, contentbase))

		tagbase := dbConnection.Tag()
		r = r.WithContext(context.WithValue(r.Context(), types.TAG, tagbase))

		viewbase := dbConnection.View()
		r = r.WithContext(context.WithValue(r.Context(), types.VIEW, viewbase))

		acronymbase := dbConnection.Acronym()
		r = r.WithContext(context.WithValue(r.Context(), types.ACRONYM, acronymbase))

		next.ServeHTTP(w, r)
	})
}

// Timeout puts a timeout on request length
func Timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.V.BasicTimeout.Duration > 0 {
			c, cancel := context.WithTimeout(r.Context(), config.V.BasicTimeout.Duration)
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

// ParseBody parses request's body allows for requests to be formed as json or normal request forms
func ParseBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var jf jsonform
			if err := json.NewDecoder(r.Body).Decode(&jf); err != nil {
				panic(srverror.New(err, 400, "Unable to decode json object"))
			}
			r.PostForm = jf.Values()
			if err := r.ParseForm(); err != nil {
				panic(srverror.New(err, 400, "Unable to parse form values"))
			}
		} else if r.Header.Get("Content-Type") == "multipart/form-data" {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				panic(srverror.New(err, 400, "Unable to parse multipart form values"))
			}
		}
		next.ServeHTTP(w, r)
	})
}
