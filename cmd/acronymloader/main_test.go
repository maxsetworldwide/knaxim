package main

import (
	"context"
	"os"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/mongo"
)

func TestCli(t *testing.T) {
	db := new(mongo.Acronymbase)
	timeoutctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	db.URI = *databaseURI
	db.DBName = *databaseName
	db.CollNames = make(map[string]string)
	db.CollNames["acronym"] = *collectionName
	err := db.Init(timeoutctx, true)
	if err != nil {
		t.Fatal("unable to init database ", err)
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal("unable to create Pipe ", err)
	}
	os.Stdin = pr
	pw.Write([]byte("ab,acronymbase"))
	pw.Close()
	main()
}
