package main

import (
	"context"
	"log"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/types"
)

func adjustUser(username string, update func(types.UserI) (types.UserI, error)) error {
	vPrintf("accessing user %s\n", username)
	ctx, cancel := context.WithTimeout(context.Background(), config.V.BasicTimeout.Duration)
	defer cancel()
	dbConnection, err := config.DB.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to database: %s\n", err)
		return err
	}
	user, err := dbConnection.Owner().FindUserName(username)
	if err != nil {
		log.Printf("Failed to find user %s: %s\n", username, err)
		return err
	}
	vPrintf("updating user %s\n", username)
	user, err = update(user)
	if err != nil {
		log.Printf("unable to update user: %s", err)
		return err
	}
	err = dbConnection.Owner().Update(user)
	if err != nil {
		log.Printf("unable to update database: %s", err)
		return err
	}
	vPrintf("user %s saved\n", username)
	return nil
}
