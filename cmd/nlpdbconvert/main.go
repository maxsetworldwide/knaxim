/*
* nlpdbconvert: takes a pre-nlpgraphs update database for Knaxim/CE and updates
* the tagbase to include new nlp tags in the tagbase structure.
* This cmd takes and mongoDB URI and two DB names - the old db name and a new
* one.
* THE NEW DB NAME WILL BE OVERWRITTEN!
 */
package main

import (
	"context"
	"errors"
	"flag"

	"git.maxset.io/web/knaxim/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
 * params:
 *   mongo URI
 *   old db name
 *   new db name
 */
var uri = flag.String("uri", "mongodb://localhost:27017", "mongodb URI")
var oldDBName = flag.String("oldname", "", "DB name to read from")
var newDBName = flag.String("newname", "", "New DB name to write to")
var overwrite = flag.Bool("overwrite", false, "Overwrite newname if it already exists")

func main() {
	flag.Parse()

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
	ctx := context.TODO()
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

	err = copyColls(ctx, mongoClient, oldName, newName)
	if err != nil {
		return err
	}

	userTags, err := convertUserTags(ctx, mongoClient, oldName)
	if err != nil {
		return err
	}
	err = newDB.Tag(ctx).Upsert(userTags...)
	if err != nil {
		return err
	}

	_, err = createNLPTags(ctx, mongoClient, oldName)
	if err != nil {
		return err
	}

	return nil
}
