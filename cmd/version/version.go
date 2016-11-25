package version

import (
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/models/context"
)

func Register(context context.Context) {
	config := context.Config

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
