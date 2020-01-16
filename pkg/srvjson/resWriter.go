package srvjson

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseWriter struct {
	internal http.ResponseWriter
	data     map[string]interface{}
}

func NewRW(in http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		internal: in,
	}
}

func (rw *ResponseWriter) init() {
	if rw.data == nil {
		rw.data = make(map[string]interface{})
	}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.internal.Header()
}

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
	default:
		rw.data["message"] = string(data)
		return len(data), nil
	}
}

func (rw *ResponseWriter) WriteHeader(sc int) {
	rw.internal.WriteHeader(sc)
}

func (rw *ResponseWriter) Set(key string, val interface{}) {
	rw.init()
	if val == nil {
		delete(rw.data, key)
		return
	}
	rw.data[key] = val
}

func (rw *ResponseWriter) Flush() error {
	return json.NewEncoder(rw.internal).Encode(rw.data)
}
