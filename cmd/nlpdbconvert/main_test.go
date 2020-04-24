package main

import (
	"context"
	"flag"
	"testing"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testURI = flag.String("testuri", "mongodb://localhost:27017", "mongodb URI")
var testOldName = flag.String("testoldname", "conversionTestOldDB", "A valid DB name to be read in the test.")
var noclean = flag.Bool("noclean", false, "If true, tests will not clean up databases after finishing, so they can be inspected manually.")
var testGotenPath = flag.String("testgoten", "http://localhost:3000", "gotenberg URI")
var testTikaPath = flag.String("testtika", "http://localhost:9998", "tika URI")

// based on provided testDB.gz
var expectedTagsPerWord = map[string]int{
	"_favorites_": 3,
	"_trash_":     1,
}

func TestConversion(t *testing.T) {
	flag.Parse()
	gotenPath = testGotenPath
	tikaPath = testTikaPath
	testctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	testClient, err := mongo.Connect(testctx, options.Client().ApplyURI(*testURI))
	if err != nil {
		t.Fatalf("Test setup error: %s", err.Error())
	}
	defer testClient.Disconnect(testctx)
	err = testClient.Ping(testctx, nil)
	if err != nil {
		t.Fatalf("Test setup error: %s", err.Error())
	}
	dbList, err := testClient.ListDatabaseNames(testctx, bson.M{"name": *testOldName})
	if len(dbList) == 0 {
		t.Fatalf("Database %s does not exist. This test requires a running instance of a test old database.", *testOldName)
	}
	testNewName := uuid.New().String()
	if testNewName == "" {
		t.Fatalf("Test setup error: uuid generation for new test db name failed")
	}
	t.Logf("New name: %s", testNewName)
	t.Run("Invalid URI", func(t *testing.T) {
		invalidURI := "thisURIShouldError"
		err := convertDB(invalidURI, "a", "b", true)
		if err == nil {
			t.Fatalf("Expected error from invalid URI. Provided '%s'", invalidURI)
		}
	})
	t.Run("Invalid Old Name", func(t *testing.T) {
		invalidName := uuid.New().String()
		err := convertDB(*testURI, invalidName, testNewName, true)
		if err == nil {
			t.Fatalf("Expected error from invalid old DB name. Provided '%s'", invalidName)
		}
	})
	t.Run("Same Name", func(t *testing.T) {
		err := convertDB(*testURI, testNewName, testNewName, true)
		if err == nil {
			t.Fatalf("Expected error when providing the same name for old and new.")
		}
	})
	t.Run("Empty Old Name", func(t *testing.T) {
		err := convertDB(*testURI, "", testNewName, true)
		if err == nil {
			t.Fatalf("Expected error when providing an empty old name.")
		}
	})
	t.Run("Empty New Name", func(t *testing.T) {
		err := convertDB(*testURI, *testOldName, "", true)
		if err == nil {
			t.Fatalf("Expected error when providing an empty new name.")
		}
	})
	t.Run("Existing New Name", func(t *testing.T) {
		existingName := uuid.New().String()
		existingColl := uuid.New().String()
		_, err := testClient.Database(existingName).Collection(existingColl).InsertOne(testctx, bson.M{"placeholder": 5})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		defer testClient.Database(existingName).Drop(testctx)
		t.Run("Overwrite Off", func(t *testing.T) {
			err = convertDB(*testURI, *testOldName, existingName, false)
			if err == nil {
				t.Fatalf("Expected error from providing existing newDB with no overwrite")
			}
		})
		t.Run("Overwrite On", func(t *testing.T) {
			err = convertDB(*testURI, *testOldName, existingName, true)
			if err != nil {
				t.Fatalf("Expected no error from providing existing newDB with overwrite")
			}
		})
	})
	t.Run("Intended Usage", func(t *testing.T) {
		err := convertDB(*testURI, *testOldName, testNewName, false)
		if err != nil {
			t.Fatalf("Expected no error from proper usage.\nProvided:\nURI:%s\nold:%s\nnew:%s\nError:%s", *uri, *testOldName, testNewName, err.Error())
		}
		if !*noclean {
			defer testClient.Database(testNewName).Drop(testctx)
		}
		t.Run("Expected Collections", func(t *testing.T) {
			var expColls = map[string]int{
				"acronym":   0,
				"chunk":     0,
				"file":      0,
				"filetags":  0,
				"group":     0,
				"lines":     0,
				"reset":     0,
				"store":     0,
				"storetags": 0,
				"user":      0,
				"view":      0,
			}
			dbColls, err := testClient.Database(testNewName).ListCollectionNames(testctx, bson.D{})
			if err != nil {
				t.Fatalf("Error getting collection names: %s", err.Error())
			}
			for _, coll := range dbColls {
				expColls[coll]++
			}
			if len(dbColls) != len(expColls) {
				t.Fatalf("Did not receive expected collections in new database.\nExpected:%+v\nReceived:%+v", expColls, dbColls)
			}
			for _, val := range expColls {
				if val != 1 {
					t.Fatalf("Did not receive expected collections in new database.\nExpected:%+v\nReceived:%+v", expColls, dbColls)
				}
			}
		})
		t.Run("Num User Tags", func(t *testing.T) {
			var userTagType tag.Type = 1 << 24
			cursor, err := testClient.Database(testNewName).Collection("filetags").Find(testctx, bson.M{
				"type": userTagType,
			})
			if err != nil {
				t.Fatalf("Error retrieving user tags: %s", err.Error())
			}
			var tags []tag.FileTag
			err = cursor.All(testctx, &tags)
			if err != nil {
				t.Fatalf("Error parsing response: %s", err.Error())
			}
			expectedNumUserTags := 0
			for _, val := range expectedTagsPerWord {
				expectedNumUserTags += val
			}
			if len(tags) != expectedNumUserTags {
				t.Fatalf("Expected %d tags from resulting database. Received %d.", expectedNumUserTags, len(tags))
			}
			tagWords := make(map[string]int)
			for _, tag := range tags {
				tagWords[tag.Word]++
			}
			if len(tagWords) != len(expectedTagsPerWord) {
				t.Fatalf("Expected %d different words in user tags. Received %d. Map: %+v", len(expectedTagsPerWord), len(tagWords), tagWords)
			}
			for key, val := range tagWords {
				expectedNum, ok := expectedTagsPerWord[key]
				if !ok {
					t.Fatalf("Received unexpected user tag: %s, Map: %+v", key, tagWords)
				}
				if expectedNum != val {
					t.Fatalf("Expected %d %s tags. Received %d. Map: %+v", expectedNum, key, val, tagWords)
				}
			}
		})
	})
}
