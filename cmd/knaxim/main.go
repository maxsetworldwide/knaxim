package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/handlers"
	"git.maxset.io/web/knaxim/internal/util"
	"git.maxset.io/web/knaxim/pkg/srverror"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var confPath = flag.String("config", "", "specify configuration file, default is Enviroment Variable KNAXIM_SERVER_CONFIG or if Enviroment Variable is missing or empty, looks for /etc/knaxim/conf.json")

var confPathShort = flag.String("c", "", "see config")

// var standardtimeout time.Duration

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	util.VerboseRequest(req, "redirecting to https")
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func setup() {
	flag.Parse()
	if len(*confPathShort) > 0 && len(*confPath) == 0 {
		confPath = confPathShort
	}
	if len(*confPath) == 0 {
		econfp := os.Getenv("KNAXIM_SERVER_CONFIG")
		if len(econfp) == 0 {
			*confPath = "/etc/knaxim/conf.json"
		} else {
			*confPath = econfp
		}
	}
	if err := config.ParseConfig(*confPath); err != nil {
		log.Fatalln("unable to parse config:", err)
	}
	log.Printf("Configuration: %+v", config.V)
	setupctx, cancel := context.WithTimeout(context.Background(), config.V.SetupTimeout.Duration)
	defer cancel()
	if err := config.DB.Init(setupctx, config.V.DatabaseReset); err != nil {
		log.Fatalf("database init error: %v\n", err)
	}
	if config.V.GuestUser != nil {
		guestUser := types.NewUser(config.V.GuestUser.Name, config.V.GuestUser.Pass, config.V.GuestUser.Email)
		guestUser.SetRole("Guest", true)
		db, err := config.DB.Connect(setupctx)
		if err != nil {
			log.Fatalf("unable to connect to database: %v", err)
		}
		userbase := db.Owner()
		if preexisting, err := userbase.FindUserName(config.V.GuestUser.Name); preexisting != nil {
			log.Printf("Guest User Already Exists")
		} else if se, ok := err.(srverror.Error); !ok || se.Status() == errors.ErrNotFound.Status() {
			if guestUser.ID, err = userbase.Reserve(guestUser.ID, guestUser.Name); err != nil {
				log.Fatalf("unable to reserve guestUser: %v", err)
			}
			if err := userbase.Insert(guestUser); err != nil {
				log.Fatalf("unable to create guestUser: %v", err)
			}
		} else {
			log.Fatalf("Error setting up guest user: %v", err)
		}
		db.Close(setupctx)
	}
}

func main() {
	setup()
	if config.T.Server != nil {
		if err := config.T.Server.Start(context.Background()); err != nil {
			log.Fatalln("Unable to start tika server: ", err)
		}
		defer config.T.Server.Shutdown(context.Background())
	}
	mainR := mux.NewRouter()

	mainR.Use(handlers.Logging)
	mainR.Use(handlers.Recovery)
	//mainR.Use(handlers.CompressHandler)
	mainR.Use(handlers.Timeout)

	{
		apirouter := mainR.PathPrefix("/api").Subrouter()
		handlers.AttachUser(apirouter.PathPrefix("/user").Subrouter())
		handlers.AttachPerm(apirouter.PathPrefix("/perm").Subrouter())
		handlers.AttachRecord(apirouter.PathPrefix("/record").Subrouter())
		handlers.AttachGroup(apirouter.PathPrefix("/group").Subrouter())
		handlers.AttachDir(apirouter.PathPrefix("/dir").Subrouter())
		handlers.AttachFile(apirouter.PathPrefix("/file").Subrouter())
		handlers.AttachPublic(apirouter.PathPrefix("/public").Subrouter())
		handlers.AttachAcronym(apirouter.PathPrefix("/acronym").Subrouter())
		handlers.AttachNLP(apirouter.PathPrefix("/nlp").Subrouter())
		handlers.AttachSearch(apirouter.PathPrefix("/search").Subrouter())
	}
	if len(config.V.StaticPath) > 0 {
		staticrouter := mainR.PathPrefix("/").Subrouter()
		staticrouter.Use(muxhandlers.CompressHandler)
		staticrouter.NewRoute().Handler(config.StaticHandler)
	}
	//change to safe close with server with time out values
	config.V.Server.Handler = mainR
	log.Println("Starting server")
	go func() {
		if config.V.Cert == nil {
			if err := config.V.Server.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			go http.ListenAndServe(config.V.Cert.HTTPport, http.HandlerFunc(redirect))
			if err := config.V.Server.ListenAndServeTLS(config.V.Cert.CertFile, config.V.Cert.KeyFile); err != nil {
				log.Println(err)
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), config.V.GracefulTimeout.Duration)
	defer cancel()
	config.V.Server.Shutdown(ctx)
	log.Println("Shutting down")
}
