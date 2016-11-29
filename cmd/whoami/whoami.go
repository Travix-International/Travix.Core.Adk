package whoami

import (
	"context"
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	appContext "github.com/Travix-International/Travix.Core.Adk/models/context"
)

const CONTEXTKEY int = 1

func Register(ctx context.Context) {
	ctxVal, err := ctx.Value(CONTEXTKEY).(appContext.Context)
	if err != nil {
		log.Errorln("General context failure")
	}
	config := ctxVal.Config

	context.App.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			// get locally stored auth info
			auth, authErr := auth.GetAuth(config)
			if authErr != nil {
				log.Fatal(authErr)
				return nil
			}

			// fetch refreshed token
			if config.Verbose {
				log.Println("Fetching refreshed token...")
			}
			refreshToken := auth.User.StsTokenManager.RefreshToken
			tokenBody, tokenBodyErr := auth.FetchRefreshedToken(config, refreshToken)
			if tokenBodyErr != nil {
				log.Fatal(tokenBodyErr)
				return nil
			}

			// fetch profile
			if config.Verbose {
				log.Println("Fetching developer profile...")
			}
			body, err := auth.FetchDeveloperProfile(config, tokenBody)
			if err != nil {
				log.Fatal(err)
				return nil
			}

			fmt.Println(body)

			return nil
		})
}
