package main

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	config "github.com/Travix-International/Travix.Core.Adk/models/config"
)

func registerWhoAmI(app *kingpin.Application, cfg *config.Config) {
	app.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			// get locally stored auth info
			authInfo, authErr := auth.GetAuth(cfg)
			if authErr != nil {
				log.Fatal(authErr)
				return nil
			}

			// fetch refreshed token
			if cfg.Verbose {
				log.Println("Fetching refreshed token...")
			}
			refreshToken := authInfo.User.StsTokenManager.RefreshToken
			tokenBody, tokenBodyErr := auth.FetchRefreshedToken(cfg, refreshToken)
			if tokenBodyErr != nil {
				log.Fatal(tokenBodyErr)
				return nil
			}

			// fetch profile
			if cfg.Verbose {
				log.Println("Fetching developer profile...")
			}
			body, err := auth.FetchDeveloperProfile(cfg, tokenBody)
			if err != nil {
				log.Fatal(err)
				return nil
			}

			fmt.Println(body)

			return nil
		})
}
