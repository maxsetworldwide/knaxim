// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		Verbose(fmt.Sprintf("%s:\n\t(%s)%s\n\t\t%s", r.RemoteAddr, r.Method, r.RequestURI, a), b...)
	}
}
