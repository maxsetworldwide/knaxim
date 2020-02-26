package util

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var verboseflag = flag.Bool("v", false, "log to stdout event messages")
var vlog *log.Logger

func initverbose() {
	if vlog == nil {
		vlog = log.New(os.Stdout, "", log.LstdFlags)
	}
}

// SetLogger changes the logger used for verbose messages
func SetLogger(l *log.Logger) {
	vlog = l
}

// Verbose write a verbose message if the verboseflag is true
func Verbose(a string, b ...interface{}) {
	if *verboseflag {
		initverbose()
		vlog.Printf(a, b...)
	}
}

// VerboseRequest writes a verbose message containing containing info
// from the http.Request, but only if the verboseflag is true
func VerboseRequest(r *http.Request, a string, b ...interface{}) {
	if *verboseflag {
		Verbose(fmt.Sprintf("%s:\n\t%s(%s)\n\t\t%s", r.RemoteAddr, r.Method, r.URL.Path, a), b...)
	}
}
