// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
