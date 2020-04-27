package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"git.maxset.io/web/knaxim/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri = flag.String("uri", "mongodb://localhost:27017", "mongodb URI")
var oldDBName = flag.String("oldname", "", "DB name to read from")
var newDBName = flag.String("newname", "", "New DB name to write to")
var overwrite = flag.Bool("overwrite", false, "Overwrite newname if it already exists")
var gotenPath = flag.String("g", "http://localhost:3000", "gotenberg URI")
var tikaPath = flag.String("t", "http://localhost:9998", "tika URI")
var quiet = flag.Bool("q", false, "Suppress console output")

func main() {
	flag.Parse()
	err := convertDB(*uri, *oldDBName, *newDBName, *overwrite)
	if !*quiet {
		if err != nil {
			fmt.Println(err.Error())
		} else {
			var verb string
			if *overwrite {
				verb = "Overwrote"
			} else {
				verb = "Created"
			}
			fmt.Printf("Done! %s database '%s'\n", verb, *newDBName)
		}
	}
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func getMongoClient(ctx context.Context, uri string) (*mongoDB.Client, error) {
	client, err := mongoDB.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func convertDB(uri, oldName, newName string, overwrite bool) error {
	if oldName == newName {
		return errors.New("old name and new name should not be the same")
	}
	if len(oldName) == 0 || len(newName) == 0 {
		return errors.New("oldName and newName should not be empty")
	}
	ctx := context.Background()
	if !*quiet {
		fmt.Println("Connecting to MongoDB...")
	}
	mongoClient, err := getMongoClient(ctx, uri)
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(ctx)
	databaseNames, err := mongoClient.ListDatabaseNames(ctx, bson.M{"name": oldName})
	if err != nil {
		return err
	}
	if len(databaseNames) == 0 {
		return errors.New("Src database does not exist")
	}
	if !overwrite {
		databaseNames, err = mongoClient.ListDatabaseNames(ctx, bson.M{"name": newName})
		if err != nil {
			return err
		}
		if len(databaseNames) != 0 {
			return errors.New("new database name already exists")
		}
	}

	newDB := new(mongo.Database) // Knaxim mongo DB
	defer newDB.Close(ctx)
	newDB.URI = uri
	newDB.DBName = newName
	if err := newDB.Init(ctx, true); err != nil {
		return err
	}

	if !*quiet {
		fmt.Println("Copying unchanged collections...")
	}
	err = copyColls(ctx, mongoClient, oldName, newName)
	if err != nil {
		return err
	}

	if !*quiet {
		fmt.Println("Converting user tags...")
	}
	userTags, err := convertUserTags(ctx, mongoClient, oldName)
	if err != nil {
		return err
	}
	err = newDB.Tag(ctx).Upsert(userTags...)
	if err != nil {
		return err
	}

	if !*quiet {
		fmt.Println("Creating and inserting views, content, and NLP tags...")
	}
	err = insertNLPTags(ctx, mongoClient, newDB)
	if err != nil {
		return err
	}

	return nil
}
