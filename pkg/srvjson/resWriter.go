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

package srvjson

// This Package implements a http.ResponseWriter that wraps another
// http.ResponseWriter. When this ResponseWriter is flushed it
// transforms whatever was written to json, and writes the json to the
// wrapped http.ResponseWriter.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseWriter implements http.ResponseWriter and converts responses
// to http
type ResponseWriter struct {
	Internal http.ResponseWriter
	data     map[string]interface{}
}

// NewRW builds a ResponseWriter wrapping the given http.ResponseWriter
func NewRW(in http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		Internal: in,
	}
}

func (rw *ResponseWriter) init() {
	if rw.data == nil {
		rw.data = make(map[string]interface{})
	}
}

// Header implements http.ResponseWriter, calles underlying
// http.ResponseWriter
func (rw *ResponseWriter) Header() http.Header {
	return rw.Internal.Header()
}

// Write implements http.ResponseWriter, preps data to be written to the
// "message" field in the resulting json object when flushed.
func (rw *ResponseWriter) Write(data []byte) (n int, err error) {
	rw.init()
	switch v := rw.data["message"].(type) {
	case io.Writer:
		return v.Write(data)
	case fmt.Stringer:
		rw.data["message"] = fmt.Sprintf("%s%s", v.String(), string(data))
		return len(data), nil
	case string:
		rw.data["message"] = fmt.Sprintf("%s%s", v, string(data))
		return len(data), nil
	case []byte:
		rw.data["message"] = fmt.Sprintf("%s%s", string(v), string(data))
		return len(data), nil
	default:
		rw.data["message"] = string(data)
		return len(data), nil
	}
}

// WriteHeader implements http.ResponseWriter, calls underlying
// http.ResponseWriter
func (rw *ResponseWriter) WriteHeader(sc int) {
	rw.Internal.WriteHeader(sc)
}

// Set assigns additional key=value pairs into the resulting json object
// when flushed. If the key is "message" it will overwrite any previous
// calls to Write with the new value.
func (rw *ResponseWriter) Set(key string, val interface{}) {
	rw.init()
	if val == nil {
		delete(rw.data, key)
		return
	}
	rw.data[key] = val
}

// Flush writes the message and set key/value pairs to the underlying
// http.ResponseWriter in json form.
func (rw *ResponseWriter) Flush() error {
	if stringer, ok := rw.data["message"].(fmt.Stringer); ok {
		rw.data["message"] = stringer.String()
	}
	return json.NewEncoder(rw.Internal).Encode(rw.data)
}
