package appix

import (
	"fmt"

	"github.com/Travix-International/appix/auth"
	"github.com/Travix-International/appix/config"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LoginCommand struct {
	*Command
}

func (cmd *LoginCommand) Register(app *kingpin.Application, config config.Config) {
	app.Command("login", "Login").
		Action(func(parseContext *kingpin.ParseContext) error {
			var url = "http://localhost:" + config.AuthServerPort

			done := make(chan interface{})
			go auth.StartServer(config, done)

			fmt.Println("Opening url: " + url)
			openURL(url)

			select {
			case <-done:
				fmt.Println("Login done!")
				close(done)
			}

			return nil
		})
}
