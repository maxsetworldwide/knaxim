package mongo

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	testclient, err := mongo.Connect(ctx, options.Client().ApplyURI(db.URI))
	if err != nil {
		t.Fatal("Unable to conntect to mongodb", err)
	}
	defer testclient.Disconnect(ctx)
	if err = testclient.Database(db.DBName).Drop(ctx); err != nil {
		t.Fatal("unable to drop DB", err)
	}
	if err := db.Init(ctx, false); err != nil {
		t.Error("Unable to init database", err)
	}
	dbnames, err := testclient.ListDatabaseNames(ctx, bson.M{
		"name": db.DBName,
	})
	if err != nil {
		t.Error("Unable to List Names", err)
	}
	if len(dbnames) != 1 {
		t.Errorf("db missing %v", dbnames)
	}
}
