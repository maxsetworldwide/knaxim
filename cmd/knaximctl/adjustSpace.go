package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/types"
)

func changeCount() {
	setup(false)
	if flag.NArg() < 3 {
		fmt.Println(helpstrs["updatefilecount"])
		return
	}
	var username = flag.Arg(1)
	var count, err = strconv.ParseInt(flag.Arg(2), 10, 64)
	if err != nil {
		log.Printf("count must be a number: %s\n%s\n", err, helpstrs["updatefilecount"])
		return
	}
	vPrintf("accessing user %s", username)
	ctx, cancel := context.WithTimeout(context.Background(), config.V.BasicTimeout.Duration)
	defer cancel()
	dbConnection, err := config.DB.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to database: %s\n", err)
		return
	}
	useri, err := dbConnection.Owner().FindUserName(username)
	if err != nil {
		log.Printf("Failed to find user %s: %s\n", username, err)
		return
	}
	user, ok := useri.(*types.User)
	if !ok {
		log.Printf("user is not a recognized type")
	}
	user.Max = count
}
