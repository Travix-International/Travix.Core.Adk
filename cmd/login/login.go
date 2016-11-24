package login

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	"github.com/Travix-International/Travix.Core.Adk/models/context"
	"github.com/Travix-International/Travix.Core.Adk/utils/openUrl"
)

func Register(context context.Context) {
	config := context.Config

	context.App.Command("login", "Login").
		Action(func(parseContext *kingpin.ParseContext) error {
			var url = "http://localhost:" + config.AuthServerPort

			ch := make(chan interface{})
			go auth.StartServer(ch, config)

			fmt.Println("Opening url: " + url)
			openUrl.OpenUrl(url)

			select {
			case <-ch:
				fmt.Println("Closing server...")
				close(ch)
			}

			return nil
		})
}
