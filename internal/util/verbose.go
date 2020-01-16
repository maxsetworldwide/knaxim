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

func Verbose(a string, b ...interface{}) {
	if *verboseflag {
		initverbose()
		vlog.Printf(a, b...)
	}
}

func VerboseRequest(r *http.Request, a string, b ...interface{}) {
	if *verboseflag {
		Verbose(fmt.Sprintf("%s:\n\t%s(%s)\n\t\t%s", r.RemoteAddr, r.Method, r.URL.Path, a), b...)
	}
}
