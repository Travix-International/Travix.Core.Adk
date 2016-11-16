package main

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureWhoamiCommand(app *kingpin.Application) {
	app.Command("whoami", "Displays logged in user's email address").
		Action(executeWhoamiCommand)
}

func executeWhoamiCommand(context *kingpin.ParseContext) error {
	config := GetConfig()

	auth, authErr := GetAuth(config)
	if authErr != nil {
		log.Fatal(authErr)
		return nil
	}

	email := auth.User.Email

	// @TODO: verify by pinging server
	fmt.Println(email)

	return nil
}
