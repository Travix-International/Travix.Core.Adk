package main

import (
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureVersionCommand(app *kingpin.Application) {
	app.Command("version", "Displays version information").
		Action(executeVersionCommand).
		Alias("ver").
		Alias("v")
}

func executeVersionCommand(context *kingpin.ParseContext) error {
	log.Printf("Version: %s", version)
	log.Printf("Hash: %s", gitHash)
	log.Printf("Build date: %s", parsedBuildDate)
	return nil
}
