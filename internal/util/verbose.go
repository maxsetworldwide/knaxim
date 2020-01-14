package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"net/http"
)

var verboseflag = flag.Bool("v", false, "log to stdout event messages")
var vlog *log.Logger

//var debugflag = flag.Bool("debug", false, "write debug messages to response Writer")
// var dlog *log.Logger

func initverbose() {
	if vlog == nil {
		vlog = log.New(os.Stdout, "", log.LstdFlags)
	}
}

func verbose(a string, b ...interface{}) {
	if *verboseflag {
		initverbose()
		vlog.Printf(a, b...)
	}
}

func verboseRequest(r *http.Request, a string, b ...interface{}) {
	if *verboseflag {
		verbose(fmt.Sprintf("%s:\n\t%s(%s)\n\t\t%s", r.RemoteAddr, r.Method, r.URL.Path, a), b...)
	}
}

// func initdebug() {
// 	if dlog == nil {
// 		dlog = log.New(os.Stdout, "", log.LstdFlags)
// 	}
// }
//
// func debug(a string, b ...interface{}) {
// 	if *debugflag {
// 		initdebug()
// 		dlog.Printf(a, b...)
// 	}
// }
//
// func debugRequest(r *http.Request, a string, b ...interface{}) {
// 	if *debugflag {
// 		debug(fmt.Sprintf("%s:\n\t%s(%s)\n\t\t%s", r.RemoteAddr, r.Method, r.URL.Path, a), b...)
// 	}
// }
