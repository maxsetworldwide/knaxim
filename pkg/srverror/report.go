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

package srverror

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// report.go contains behavior for logging server errors by sending an email to
// the devs and by saving them to a file.

const (
	fileNameFormat = "2006-01-02"
	yearFormat     = "2006"
	monthFormat    = "01"
	dayFormat      = "02"
	indentString   = "\t"
)

// LogPath is the directory to store all logs. This should be set externally
var LogPath = ""

// LogString Returns a tab-indented human readable message containing information about
// the given error, request, and response
func LogString(err error, req *http.Request, res http.ResponseWriter) string {
	var status string
	switch cast := err.(type) {
	case Error:
		status = "Status Code: " + strconv.Itoa(cast.Status())
	default:
		status = ""
	}

	var cookies []interface{}
	for _, cookie := range req.Cookies() {
		cookies = append(cookies, cookie.String())
	}

	var responseHeaders []interface{}
	for key, arr := range res.Header() {
		responseHeaders = append(responseHeaders, key)
		vals := []interface{}{}
		for _, val := range arr {
			vals = append(vals, val)
		}
		responseHeaders = append(responseHeaders, vals)
	}

	contentLength := fmt.Sprintf("Content Length: %d", req.ContentLength)

	logFormat := []interface{}{
		time.Now().Format(time.RFC1123),
		[]interface{}{
			"Error Message",
			[]interface{}{
				err.Error(),
				status,
			},
		},
		[]interface{}{
			"Request",
			[]interface{}{
				"Path",
				[]interface{}{
					req.URL.String(),
				},
				"Method",
				[]interface{}{
					req.Method,
				},
				"Remote Address",
				[]interface{}{
					req.RemoteAddr,
				},
				contentLength,
				"Cookies",
				cookies,
			},
		},
		[]interface{}{
			"Response",
			[]interface{}{
				"Current Header Values",
				responseHeaders,
			},
		},
	}

	var result strings.Builder
	buildSections(logFormat, &result, 0)
	return result.String()
}

// Recursive function for converting log format to string format.
// Sections are skipped if their content is empty.
func buildSections(sections []interface{}, result *strings.Builder, indent int) {
	for _, content := range sections {
		switch v := content.(type) {
		case string:
			if len(v) > 0 {
				writeIndent(result, indent)
				result.WriteString(v)
				result.WriteRune('\n')
			}
		case []interface{}:
			buildSections(v, result, indent+1)
		default:
			writeIndent(result, indent)
			result.WriteString("Unrecognized data type")
			result.WriteRune('\n')
		}
	}
}

func writeIndent(builder *strings.Builder, num int) {
	for i := 0; i < num; i++ {
		builder.WriteString(indentString)
	}
}

// WriteToFile writes the given message to the log directory specified by LogPath.
// The directory will have the structure of LogPath/YYYY/MM/DD.log
func WriteToFile(msg string) (err error) {
	if len(LogPath) == 0 {
		return errors.New("no log path specified")
	}
	currTime := time.Now()

	year := currTime.Format(yearFormat)
	month := currTime.Format(monthFormat)
	day := currTime.Format(dayFormat)

	logDirPath := filepath.Join(LogPath, year, month)
	completeLogPath := filepath.Join(LogPath, year, month, day+".log")
	err = os.MkdirAll(logDirPath, os.ModePerm)
	if err != nil {
		return
	}
	file, err := os.OpenFile(completeLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()
	_, err = file.Write([]byte(msg))
	return
}
