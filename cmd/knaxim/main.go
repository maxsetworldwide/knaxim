package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var conf_path = flag.String("config", "", "specify configuration file, default is Enviroment Variable KNAXIM_SERVER_CONFIG or if Enviroment Variable is missing or empty, looks for /etc/knaxim/conf.json")

var conf_path_short = flag.String("c", "", "see config")

// var standardtimeout time.Duration

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	verboseRequest(req, "redirecting to https")
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func setup() {
	flag.Parse()
	if len(*conf_path_short) > 0 && len(*conf_path) == 0 {
		conf_path = conf_path_short
	}
	if len(*conf_path) == 0 {
		econfp := os.Getenv("KNAXIM_SERVER_CONFIG")
		if len(econfp) == 0 {
			*conf_path = "/etc/knaxim/conf.json"
		} else {
			*conf_path = econfp
		}
	}
	if err := config.ParseConfig(*conf_path); err != nil {
		log.Fatalln("unable to parse config:", err)
	}
	//log.Printf("Configuration: %v", conf);
	setupctx, cancel := context.WithTimeout(context.Background(), config.V.SetupTimeout)
	defer cancel()
	if err := config.DB.Init(setupctx, config.V.DatabaseReset); err != nil {
		log.Fatalf("database init error: %v\n", err)
	}
	if config.V.GuestUser != nil {
		guestUser := database.NewUser(config.V.GuestUser.Name, config.V.GuestUser.Pass, config.V.GuestUser.Email)
		guestUser.SetRole("Guest", true)
		userbase := db.Owner(setupctx)
		if preexisting, err := userbase.FindUserName(conf.GuestUser.Name); preexisting != nil {
			log.Printf("Guest User Already Exists")
		} else if err == database.ErrNotFound {
			if guestUser.ID, err = userbase.Reserve(guestUser.ID, guestUser.Name); err != nil {
				log.Fatalf("unable to reserve guestUser: %v", err)
			}
			if err := userbase.Insert(guestUser); err != nil {
				log.Fatalf("unable to create guestUser: %v", err)
			}
		} else {
			log.Fatalf("Error setting up guest user: %v", err)
		}
		userbase.Close(setupContext)
	}
}

func main() {
	if config.T.Server != nil {
		if err := config.T.Server.Start(context.Background()); err != nil {
			log.Fatalln("Unable to start tika server: ", err)
		}
		defer config.T.Server.Shutdown(context.Background())
	}
	mainR := mux.NewRouter()

	mainR.Use(loggingMiddleware)
	mainR.Use(RecoveryMiddleWare)
	//mainR.Use(handlers.CompressHandler)
	mainR.Use(timeoutMiddleware)
	//mainR.Use(databaseMiddleware)

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
		//setupSearch(apirouter.PathPrefix("/s").Subrouter())
	}
	if len(conf.StaticPath) > 0 {
		staticrouter := mainR.PathPrefix("/").Subrouter()
		staticrouter.Use(handlers.CompressHandler)
		staticrouter.NewRoute().Handler(http.FileServer(http.Dir(conf.StaticPath)))
	}
	//change to safe close with server with time out values
	conf.Server.Handler = mainR
	log.Println("Starting server")
	go func() {
		if conf.Cert == nil {
			if err := conf.Server.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			go http.ListenAndServe(conf.Cert.HTTPport, http.HandlerFunc(redirect))
			if err := conf.Server.ListenAndServeTLS(conf.Cert.CertFile, conf.Cert.KeyFile); err != nil {
				log.Println(err)
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), conf.GracefulTimeout)
	defer cancel()
	conf.Server.Shutdown(ctx)
	log.Println("Shutting down")
}
