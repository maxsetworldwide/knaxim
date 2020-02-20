package config

import (
	"encoding/json"
	"errors"
	"math"
	"os"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/mongo"
	"git.maxset.io/web/knaxim/internal/handlers/spa"
	"github.com/google/go-tika/tika"
)

// V => Configuration of Knaxim
// populate with ParseConfig
var V Configuration

// DB => database instance specified by configuration
var DB database.Database

// T => Tika server and connection information
// if Server is nil, path should hold the http path to the tika server
var T struct {
	Path   string
	Server *tika.Server
}

var StaticHandler spa.SpaHandler

// ParseConfig loads configuration file and populates global vars
func ParseConfig(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(fp).Decode(&V); err != nil {
		return err
	}
	fp.Close()
	if V.FileLimit == 0 {
		V.FileLimit = 5 * 1024 * 1024
	} else if V.FileLimit < 0 {
		V.FileLimit = math.MaxInt64
	}
	StaticHandler.StaticPath = V.StaticPath
	StaticHandler.IndexPath = V.IndexPath
	switch V.DatabaseType {
	case "mongo":
		DB = new(mongo.Database)
	case "memory":
		DB = new(memory.Database)
	default:
		return errors.New("Unrecognized config database type")
	}
	if err = json.Unmarshal(V.Database, DB); err != nil {
		return err
	}
	if V.Tika.Type == "local" {
		T.Server, err = tika.NewServer(V.Tika.Path, V.Tika.Port)
		if err != nil {
			return err
		}
		T.Server.ChildMode(&tika.ChildOptions{
			MaxFiles:          V.Tika.MaxFiles,
			TaskPulseMillis:   V.Tika.TaskPulse,
			TaskTimeoutMillis: V.Tika.TaskTimeout,
			PingPulseMillis:   V.Tika.PingPulse,
			PingTimeoutMillis: V.Tika.PingTimeout,
		})
		T.Path = T.Server.URL()
	} else if V.Tika.Type == "external" {
		if V.Tika.Port == "" {
			V.Tika.Port = "9998"
		}
		T.Path = V.Tika.Path + ":" + V.Tika.Port
	} else {
		return errors.New("unrecognized tika config type")
	}
	return nil
}
