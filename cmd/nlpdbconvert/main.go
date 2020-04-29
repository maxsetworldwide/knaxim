package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"git.maxset.io/web/knaxim/internal/database/mongo"
	"git.maxset.io/web/knaxim/internal/database/types"
	CEErrors "git.maxset.io/web/knaxim/internal/database/types/errors"
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

	perrs, err := findPerrs(ctx, mongoClient, newName)
	if err != nil {
		return err
	}

	var errStrs []string
	for _, perr := range perrs {
		errStrs = append(errStrs, perr.Error())
	}

	if len(errStrs) == 0 {
		return nil
	}
	return fmt.Errorf("%d processing errors occurred during conversion:\n%s", len(errStrs), strings.Join(errStrs, "\n"))
}

func findPerrs(ctx context.Context, client *mongoDB.Client, dbName string) ([]CEErrors.Processing, error) {
	perrs := []CEErrors.Processing{}
	db := client.Database(dbName)
	storeColl := db.Collection("store")
	cursor, err := storeColl.Find(ctx, bson.M{})
	if err != nil {
		return perrs, err
	}
	var stores []types.FileStore
	err = cursor.All(ctx, &stores)
	if err != nil {
		return perrs, err
	}

	for _, fs := range stores {
		if fs.Perr != nil {
			perrs = append(perrs, *fs.Perr)
		}
	}
	return perrs, nil
}
