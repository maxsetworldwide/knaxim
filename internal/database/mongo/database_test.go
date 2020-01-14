package mongo

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"
)

var configuration = struct {
	DB *Database
}{
	DB: &Database{
		URI: "mongodb://localhost:27017",
	},
}

func init() {
	conffile, err := os.Open("test/mongoconfig.json")
	if err == nil {
		json.NewDecoder(conffile).Decode(&configuration)
		conffile.Close()
	}
	flag.StringVar(&configuration.DB.URI, "dbpath", configuration.DB.URI, "specify path to mongodb instance to run tests on")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestDatabaseInit(t *testing.T) {
	t.Parallel()
	db := new(Database)
	*db = *configuration.DB
	db.DBName = "TestInit"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := db.Init(ctx, true); err != nil {
		t.Error("Unable to init database", err)
	}
}
