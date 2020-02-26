package srverror

// This package implements an error type that is meant to be used
// within a server. The Error also implements the http.Handler so that
// if an error occurs it can be used to respond to the request directly

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	//"fmt"
)

// Error type that can also be used to handle http requests
type Error interface {
	error
	// normally only sends the first string in the messsages. In debug
	// mode sends the full Error message
	http.Handler
	Unwrap() error
	// creates new error with addition messages
	Extend(msgs ...string) Error
	Status() int
}

type srverr struct {
	status int
	msgs   []string
	e      error
}

// New creates a new instance of an Error wrapping the given error.
// status and msgs define what to respond to the request with.
func New(e error, status int, msgs ...string) Error {
	if e == nil {
		return nil
	}
	if se, ok := e.(*srverr); ok {
		return &srverr{
			status: status,
			msgs:   append(append(msgs, "\n\tPrevious Server Error= Previous Status Code >", strconv.Itoa(se.status)), se.msgs...),
			e:      se.e,
		}
	}
	return &srverr{
		status: status,
		msgs:   msgs,
		e:      e,
	}
}

// Basic creates a new instance of the server error without wrapping
// another error. See New
func Basic(status int, msgs ...string) Error {
	return &srverr{
		status: status,
		msgs:   msgs,
		e:      errors.New("Server Error"),
	}
}

// Error implements error
func (se *srverr) Error() string {
	if se == nil {
		return ""
	}
	return strings.Join(se.msgs, "--") + "--" + se.e.Error()
}

// DEBUG is used to flag if the server errors are to respond with debug
// messages when handling http. Should be false when running during
// production
var DEBUG = false

// ServeHTTP to implement http.Handler, normally only sends the first
// string in the msgs. In debug mode sends the full Error message
func (se *srverr) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if se == nil {
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	} else {
		w.WriteHeader(se.status)
		var msg string
		if DEBUG {
			msg = se.Error()
		} else if len(se.msgs) > 0 {
			msg = se.msgs[0]
		}
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": msg,
		}); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Server Error Error"))
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
	}
}

func (se *srverr) Unwrap() error {
	if se == nil {
		return nil
	}
	return se.e
}

func (se *srverr) Extend(msgs ...string) Error {
	if se == nil {
		return &srverr{
			status: 500,
			msgs:   msgs,
			e:      errors.New("nil srverrror"),
		}
	}
	nmsgs := make([]string, 0, len(se.msgs)+len(msgs))
	nmsgs = append(nmsgs, se.msgs...)
	nmsgs = append(nmsgs, msgs...)
	return &srverr{
		status: se.status,
		e:      se.e,
		msgs:   nmsgs,
	}
}

func (se *srverr) Status() int {
	return se.status
}

// TODO: Add feature that adds Header Fields
