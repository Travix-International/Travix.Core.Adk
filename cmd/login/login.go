package login

import (
	"context"
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	appContext "github.com/Travix-International/Travix.Core.Adk/models/context"
	"github.com/Travix-International/Travix.Core.Adk/utils/openUrl"
)

func Register(ctx context.Context) {
	ctxVal, err := ctx.Value(CONTEXTKEY).(appContext.Context)
	if err != nil {
		log.Errorln("General context failure")
	}
	config := ctxVal.Config

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
