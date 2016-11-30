package login

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth_server"
	"github.com/Travix-International/Travix.Core.Adk/lib/cmd"
	"github.com/Travix-International/Travix.Core.Adk/lib/context"
	"github.com/Travix-International/Travix.Core.Adk/utils/openUrl"
)

type LoginCommand struct {
	*cmd.Command
}

func (cmd *LoginCommand) Register(context context.Context) {
	config := context.Config

	context.App.Command("login", "Login").
		Action(func(parseContext *kingpin.ParseContext) error {
			var url = "http://localhost:" + config.AuthServerPort

			done := make(chan interface{})
			go auth_server.StartServer(config, done)

			fmt.Println("Opening url: " + url)
			openUrl.OpenUrl(url)

			select {
			case <-done:
				fmt.Println("Closing server...")
				close(done)
			}

			return nil
		})
}
