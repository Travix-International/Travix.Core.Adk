package context

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
	modelsConfig "github.com/Travix-International/Travix.Core.Adk/models/config"
)

type Context struct {
	App    *kingpin.Application
	Config *modelsConfig.Config
	Auth   *auth.Auth
}

func (context *Context) RequireUserLoggedIn(command string) {
	auth, authErr := auth.GetAuth(context.Config)
	if authErr != nil {
		fmt.Println("You must be logged in order to run '" + command + "' command")
		fmt.Println("Use: appix login")
		os.Exit(1)
	}
	context.Auth = auth
}
