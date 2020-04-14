package main

/*
 * This is an integration test that requires an external mongodb session to be
 * running. The URI will be grabbed from the command line when running the test.
 * A name for an existing and valid old DB should be provided as well. JSON
 * files should be provided with this package so they can be easily created via
 * Compass.
 */

import (
	"context"
	"flag"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testUri = flag.String("testuri", "mongodb://localhost:27017", "mongodb URI")
var testOldName = flag.String("testoldname", "", "A valid DB name to be read in the test.")
var testNewName = "testNewDBName"

// var noCleanup = flag.Bool("noclean", false, "If true, tests will not clean up databases after finishing, so they can be inspected manually.")

//wanted to do a test that checks that old database is not changed, but doesn't
//seem entirely possible
func TestConversion(t *testing.T) {
	flag.Parse()
	testctx := context.TODO()
	testClient, err := mongo.Connect(testctx, options.Client().ApplyURI(*testUri))
	if err != nil {
		t.Fatalf("Test setup error: %s", err.Error())
	}
	defer testClient.Disconnect(testctx)
	if err != nil {
		t.Fatalf("Test setup error: %s", err.Error())
	}
	t.Run("Invalid URI", func(t *testing.T) {
		invalidUri := "thisURIShouldError"
		err := convertDB(invalidUri, "a", "b")
		if err == nil {
			t.Fatalf("Expected error from invalid URI. Provided '%s'", invalidUri)
		}
	})
	t.Run("Invalid Old Name", func(t *testing.T) {
		invalidName := "thisDBNameShouldNotExist"
		err := convertDB(*testUri, invalidName, testNewName)
		t.Logf("err: %s", err)
		if err == nil {
			t.Fatalf("Expected error from invalid old DB name. Provided '%s'", invalidName)
		}
	})
	t.Run("Same Name", func(t *testing.T) {
		err := convertDB(*testUri, testNewName, testNewName)
		if err == nil {
			t.Fatalf("Expected error when providing the same name for old and new.")
		}
	})
	t.Run("Empty Old Name", func(t *testing.T) {
		err := convertDB(*testUri, "", testNewName)
		if err == nil {
			t.Fatalf("Expected error when providing an empty old name.")
		}
	})
	t.Run("Empty New Name", func(t *testing.T) {
		err := convertDB(*testUri, *testOldName, "")
		if err == nil {
			t.Fatalf("Expected error when providing an empty new name.")
		}
	})
	t.Run("Intended Usage", func(t *testing.T) {
		err := convertDB(*testUri, *testOldName, testNewName)
		if err != nil {
			t.Fatalf("Expected no error from proper usage.\nProvided:\nURI:%s\nold:%s\nnew:%s\nError:%s", *uri, *testOldName, testNewName, err.Error())
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
			if len(dbColls) != len(expColls) {
				t.Fatalf("Did not receive expected collections in new database.\nExpected:%+v\nReceived:%+v", expColls, dbColls)
			}
			for _, coll := range dbColls {
				expColls[coll] += 1
			}
			for _, val := range expColls {
				if val != 1 {
					t.Fatalf("Did not receive expected collections in new database.\nExpected:%+v\nReceived:%+v", expColls, dbColls)
				}
			}
		})
	})
}
