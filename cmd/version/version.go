package version

import (
	"context"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	appContext "github.com/Travix-International/Travix.Core.Adk/models/context"
)

func Register(ctx context.Context) {
	ctxVal, err := ctx.Value(CONTEXTKEY).(appContext.Context)
	if err != nil {
		log.Errorln("General context failure")
	}
	config := ctxVal.Config

	context.App.Command("version", "Displays version information").
		Action(func(parseContext *kingpin.ParseContext) error {
			log.Printf("Version: %s", config.Version)
			log.Printf("Hash: %s", config.GitHash)
			log.Printf("Build date: %s", config.ParsedBuildDate)

			return nil
		}).
		Alias("ver").
		Alias("v")
}
