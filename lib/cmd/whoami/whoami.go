package whoami

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	libAuth "github.com/Travix-International/Travix.Core.Adk/lib/auth"
	"github.com/Travix-International/Travix.Core.Adk/lib/cmd"
	"github.com/Travix-International/Travix.Core.Adk/lib/context"
)

type WhoamiCommand struct {
	*cmd.Command
}

func (cmd *WhoamiCommand) Register(context context.Context) {
	config := context.Config

	context.App.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			// get locally stored auth info
			auth, authErr := libAuth.GetAuthData(config.AuthFilePath)
			if authErr != nil {
				log.Fatal(authErr)
				return nil
			}

			// fetch refreshed token
			if cmd.Verbose {
				log.Println("Fetching refreshed token...")
			}
			tokenBody, tokenBodyErr := auth.RefreshToken(config.FirebaseApiKey)
			if tokenBodyErr != nil {
				log.Fatal(tokenBodyErr)
				return nil
			}

			// fetch profile
			if cmd.Verbose {
				log.Println("Fetching developer profile...")
			}
			profileBody, profileBodyErr := tokenBody.FetchDeveloperProfile(config.DeveloperProfileUrl)
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
