package main

import (
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	config "github.com/Travix-International/Travix.Core.Adk/models/config"
)

func registerVersion(app *kingpin.Application, cfg *config.Config) {
	app.Command("version", "Displays version information").
		Action(func(parseContext *kingpin.ParseContext) error {
			log.Printf("Version: %s", cfg.Version)
			log.Printf("Hash: %s", cfg.GitHash)
			log.Printf("Build date: %s", cfg.ParsedBuildDate)

			return nil
		}).
		Alias("ver").
		Alias("v")
}
