package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureLoginCommand(app *kingpin.Application) {
	app.Command("login", "Login").
		Action(executeLoginCommand)
}

func executeLoginCommand(context *kingpin.ParseContext) error {
	fmt.Println("Login here...")
	return nil
}
