package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	config "github.com/Travix-International/Travix.Core.Adk/models/config"
	"github.com/Travix-International/Travix.Core.Adk/utils"
)

func registerLogin(app *kingpin.Application, cfg *config.Config) {
	app.Command("login", "Login").
		Action(func(parseContext *kingpin.ParseContext) error {
			url, done := auth.StartServer(cfg)

			fmt.Println("Opening url: " + url)
			utils.OpenUrl(url)
			select {
			case <-done:
			}

			return nil
		})
}
