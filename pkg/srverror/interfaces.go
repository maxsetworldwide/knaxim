package srverror

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	//"fmt"
)

type Error interface {
	error
	http.Handler
	Unwrap() error
}

type srverr struct {
	status int
	msgs   []string
	e      error
}

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

func Basic(status int, msgs ...string) Error {
	return &srverr{
		status: status,
		msgs:   msgs,
		e:      errors.New("Server Error"),
	}
}

func (se *srverr) Error() string {
	if se == nil {
		return ""
	}
	return strings.Join(se.msgs, "--") + "--" + se.e.Error()
}

func (se *srverr) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if se == nil {
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	} else {
		w.WriteHeader(se.status)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": se.msgs[0],
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

//Add feature that adds Header Fields
