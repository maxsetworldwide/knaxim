package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"

	"git.maxset.io/web/knaxim/internal/config"
)

func loadAcronyms(in io.Reader) error {
	parser := csv.NewReader(in)
	parser.FieldsPerRecord = 2
	parser.LazyQuotes = true
	parser.ReuseRecord = true
	parser.TrimLeadingSpace = true

	ctx, cancel := context.WithTimeout(context.Background(), config.V.BasicTimeout.Duration)
	defer cancel()
	dbConnection, err := config.DB.Connect(ctx)
	if err != nil {
		log.Printf("Unable to connect to Database: %s\n", err)
		return err
	}
	defer dbConnection.Close(ctx)
	ab := dbConnection.Acronym()
	var pair []string
	for pair, err = parser.Read(); err == nil; pair, err = parser.Read() {
		err := ab.Put(pair[0], pair[1])
		if err != nil {
			log.Printf("Unable to add acronym (%s, %s): %s\n", pair[0], pair[1], err)
			return err
		}
	}
	if err != io.EOF {
		log.Fatalf("error reading file: %s\n", err)
		return err
	}
	return nil
}
