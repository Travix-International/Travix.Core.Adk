package whoami

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	libAuth "github.com/Travix-International/Travix.Core.Adk/lib/auth"
	modelsContext "github.com/Travix-International/Travix.Core.Adk/models/context"
)

func Register(context modelsContext.Context) {
	config := context.Config

	context.App.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			// get locally stored auth info
			auth, authErr := libAuth.GetAuth(config)
			if authErr != nil {
				log.Fatal(authErr)
				return nil
			}

			// fetch refreshed token
			refreshToken := auth.User.StsTokenManager.RefreshToken
			tokenBody, tokenBodyErr := libAuth.FetchRefreshedToken(config, refreshToken)
			if tokenBodyErr != nil {
				log.Fatal(tokenBodyErr)
				return nil
			}

			// fetch profile
			profileBody, profileBodyErr := libAuth.FetchDeveloperProfile(config, tokenBody)
			if profileBodyErr != nil {
				log.Fatal(profileBodyErr)
				return nil
			}

			if profileBody.HasProfile {
				fmt.Println("Email: " + profileBody.Profile.Email)
				fmt.Println("Name: " + profileBody.Profile.Name)

				if profileBody.Profile.IsEnabled == true {
					fmt.Println("Enabled: Yes")
				} else {
					fmt.Println("Enabled: No")
				}

				if profileBody.Profile.IsVerified == true {
					fmt.Println("Verified: Yes")
				} else {
					fmt.Println("Verified: No")
				}

				fmt.Println("Publisher ID: " + profileBody.Profile.PublisherId)
			} else {
				fmt.Println("No profile found.")
			}

			return nil
		})
}
