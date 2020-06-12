package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	defer dbConnection.Close(ctx)
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

func generatePass() string {
	buffer := make([]byte, 9)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(buffer)
}

func newUser(name, email, pass string) (*types.User, error) {
	vPrintf("creating user %s\n", name)
	user := types.NewUser(name, pass, email)
	vPrintf("connecting to database\n")
	ctx, cancel := context.WithTimeout(context.Background(), config.V.BasicTimeout.Duration)
	defer cancel()
	dbConnection, err := config.DB.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to database: %s\n", err)
		return nil, err
	}
	defer dbConnection.Close(ctx)
	vPrintf("transfering user\n")
	if user.ID, err = dbConnection.Owner().Reserve(user.GetID(), user.GetName()); err != nil {
		log.Printf("Failed to reserve user id: %s", err)
		return nil, err
	}
	if err = dbConnection.Owner().Insert(user); err != nil {
		log.Printf("Failed to insert user: %s", err)
		return nil, err
	}
	return user, nil
}
