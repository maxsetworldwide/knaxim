/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"git.maxset.io/web/knaxim/internal/database/mongo"
)

var databaseURI = flag.String("db", "mongodb://localhost:27017", "specify address of mongodb containing acronyms")
var databaseName = flag.String("dbname", "Knaxim", "specify mongodb database name")
var collectionName = flag.String("collection", "acronym", "specify mongodb acronym collection")

var loadfile = flag.String("f", "", "specify file to upload, default reads from stdin")

var timeout = flag.Duration("dur", time.Minute, "specify the maximum time that the program should take to complete")

var initdb = flag.Bool("init", false, "init database when present")

func main() {
	flag.Parse()
	var err error
	var file io.Reader

	if len(*loadfile) > 0 {
		file, err = os.Open(*loadfile)
		if err != nil {
			log.Fatalln("Unable to read: ", *loadfile, err)
		}
		defer file.(*os.File).Close()
	} else {
		file = os.Stdin
	}

	parser := csv.NewReader(file)
	parser.FieldsPerRecord = 2
	parser.LazyQuotes = true
	parser.ReuseRecord = true
	parser.TrimLeadingSpace = true

	ab := new(mongo.Acronymbase)
	timeoutctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	ab.URI = *databaseURI
	ab.DBName = *databaseName
	ab.CollNames = make(map[string]string)
	ab.CollNames["acronym"] = *collectionName
	err = ab.Init(timeoutctx, *initdb)
	if err != nil {
		log.Fatal("unable to init database ", err)
	}
	{
		db, err := ab.Connect(timeoutctx)
		if err != nil {
			log.Fatal("unable to connect to database ", err)
		}
		defer db.Close(timeoutctx)
		ab = db.Acronym().(*mongo.Acronymbase)
	}
	var pair []string
	for pair, err = parser.Read(); err == nil; pair, err = parser.Read() {
		dbErr := ab.Put(pair[0], pair[1])
		if dbErr != nil {
			log.Fatalln("Database Error: ", dbErr.Error())
		}
	}
	ab.Close(timeoutctx)
	if err != io.EOF {
		log.Fatalln("error reading file: ", err.Error())
	}
}
