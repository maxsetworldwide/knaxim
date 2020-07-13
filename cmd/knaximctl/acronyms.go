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
