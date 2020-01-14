package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	flag.Parse()
	//setup
	conf = testingconfiguration
	verboseflag = new(bool)
	*verboseflag = true
	db = testingdb
	err := populateDB()
	defer func() {
		if tserver != nil {
			tserver.Stop()
		}
	}()
	if err != nil {
		panic(err)
	}
	server = httptest.NewServer(setupRouter())
	i := m.Run()
	//close
	if tserver != nil {
		tserver.Stop()
	}
	os.Exit(i)
}

func setupRouter() http.Handler {
	mainR := mux.NewRouter()
	mainR.Use(loggingMiddleware)
	mainR.Use(RecoveryMiddleWare)
	mainR.Use(timeoutMiddleware)
	{
		apirouter := mainR.PathPrefix("/api").Subrouter()
		apirouter.Use(databaseMiddleware)
		apirouter.Use(parseMiddleware)
		setupUser(apirouter.PathPrefix("/user").Subrouter())
		setupPerm(apirouter.PathPrefix("/perm").Subrouter())
		setupRecord(apirouter.PathPrefix("/record").Subrouter())
		setupGroup(apirouter.PathPrefix("/group").Subrouter())
		setupDir(apirouter.PathPrefix("/dir").Subrouter())
		setupFile(apirouter.PathPrefix("/file").Subrouter())
		setupPublic(apirouter.PathPrefix("/public").Subrouter())
		setupAcronym(apirouter.PathPrefix("/acronym").Subrouter())
	}
	return mainR
}

var loggingBuffer = new(bytes.Buffer)

func init() {
	vlog = log.New(loggingBuffer, "", log.LstdFlags)
}

func LogBuffer(t *testing.T) {
	t.Helper()
	logs := new(strings.Builder)
	loggingBuffer.WriteTo(logs)
	t.Log(logs.String())
}
