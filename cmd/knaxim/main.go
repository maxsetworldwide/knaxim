package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.maxset.io/server/knaxim/database"
	"git.maxset.io/server/knaxim/database/mongo"

	"math"

	"github.com/google/go-tika/tika"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type SslCert struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
	HTTPport string `json:"http_port"`
}

type configuration struct {
	Address         string
	StaticPath      string          `json:"static"`
	Server          *http.Server    `json:"server"`
	Cert            *SslCert        `json:"cert"`
	GracefulTimeout time.Duration   `json:"close_time"`
	BasicTimeout    time.Duration   `json:"basic_timeout"`
	FileTimeoutRate int64           `json:"file_timeout_rate"` //nanoseconds per 1 KB
	MaxFileTimeout  time.Duration   `json:"max_file_timeout"`
	MinFileTimeout  time.Duration   `json:"min_file_timeout"`
	DatabaseType    string          `json:"db_type"`
	Database        json.RawMessage `json:"db"`
	DatabaseReset   bool            `json:"db_clear"`
	Tika            tikaconf        `json:"tika"`
	FileLimit       int64           `json:"filelimit"`
	Smtp            struct {
		Active   bool
		Identity string
		Username string
		Password string
		Host     string
		Path     string
		From     string
	}
	// Templates struct {
	//   Files []string
	//   ConfirmEmail string
	// }
	FreeSpace int `json:"total_free_space"`
	AdminKey  string
	GuestUser *guestconf
}

type guestconf struct {
	Name  string
	Pass  string
	Email string
}

type tikaconf struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Port        string `json:"port"`
	MaxFiles    int    `json:"child_max_files"`
	TaskPulse   int    `json:"child_task_pulse"`
	TaskTimeout int    `json:"child_task_timeout"`
	PingPulse   int    `json:"child_ping_pulse"`
	PingTimeout int    `json:"child_ping_timeout"`
}

var tikapath string
var tserver *tika.Server

var conf configuration
var conf_path = flag.String("c", "", "specify configuration file, default is Enviroment Variable KX_SERVER_CONF or if Enviroment Variable is missing or empty, looks for $CWD/conf.json")

var db database.Database

var standardtimeout time.Duration

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	verboseRequest(req, "redirecting to https")
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {
	flag.Parse()
	if len(*conf_path) == 0 {
		econfp := os.Getenv("KX_SERVER_CONF")
		if len(econfp) == 0 {
			*conf_path = "conf.json"
		} else {
			*conf_path = econfp
		}
	}
	fp, err := os.Open(*conf_path)
	if err != nil {
		log.Fatalf("Unable to open configuration file: %v\n", err)
	}
	dec := json.NewDecoder(fp)
	if err = dec.Decode(&conf); err != nil {
		log.Fatalf("Unable to decode configuration file: %v\n", err)
	}
	//log.Printf("Configuration: %v", conf);
	fp.Close()
	if conf.FileLimit == 0 {
		conf.FileLimit = 50 * 1024 * 1024 //Defualt 50MB file limit size
	} else if conf.FileLimit < 0 {
		conf.FileLimit = math.MaxInt64 //-1  = no limit
	}
	switch conf.DatabaseType {
	case "mongo":
		db = new(mongo.Database)
	default:
		log.Fatalln("Unrecognized database type")
	}
	if err := json.Unmarshal(conf.Database, db); err != nil {
		log.Fatalf("Unable to decode Database configuration: %v\n", err)
	}
	//log.Printf("Configuration: %v", conf);
	setupctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := db.Init(setupctx, conf.DatabaseReset); err != nil {
		log.Fatalf("database init error: %v\n", err)
	}
	if conf.Tika.Type == "local" {
		var err error
		tserver, err = tika.NewServer(conf.Tika.Path, conf.Tika.Port)
		if err != nil {
			log.Fatalf(err.Error())
		}
		tserver.ChildMode(&tika.ChildOptions{
			MaxFiles:          conf.Tika.MaxFiles,
			TaskPulseMillis:   conf.Tika.TaskPulse,
			TaskTimeoutMillis: conf.Tika.TaskTimeout,
			PingPulseMillis:   conf.Tika.PingPulse,
			PingTimeoutMillis: conf.Tika.PingTimeout,
		})
		startctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tserver.Start(startctx); err != nil {
			log.Fatalf("tika start error: %s", err.Error())
		}
		tikapath = tserver.URL()
	} else if conf.Tika.Type == "external" {
		if conf.Tika.Port == "" {
			conf.Tika.Port = "9998"
		}
		tikapath = conf.Tika.Path + ":" + conf.Tika.Port
	} else {
		log.Fatalf("unrecognized Tika Type")
	}

	if conf.GuestUser != nil {
		guestUser := database.NewUser(conf.GuestUser.Name, conf.GuestUser.Pass, conf.GuestUser.Email)
		guestUser.SetRole("Guest", true)
		setupContext, close := context.WithTimeout(context.Background(), time.Second*10)
		defer close()
		userbase := db.Owner(setupContext)
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
	if tserver != nil {
		defer tserver.Shutdown(context.Background())
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
