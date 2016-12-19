package appix

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/auth"
	"github.com/Travix-International/appix/config"
)

type WhoamiCommand struct {
	*Command
}

func (cmd *WhoamiCommand) Register(app *kingpin.Application, config config.Config) {
	app.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			authToken, err := auth.LoadAuthToken(config)

			if err != nil {
				fmt.Println("You are not logged in.\nYou can sign in by using 'appix login'.")
				return nil
			}

			// fetch profile
			if cmd.Verbose {
				log.Println("Fetching developer profile...")
			}
			profileBody, profileBodyErr := auth.FetchDeveloperProfile(authToken, config.DeveloperProfileUrl)
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
				fmt.Println("Please log in at https://developerportal.travix.com/ to create a developer profile.")
			}

			return nil
		})
}
