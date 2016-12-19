package appix

import (
	"log"

	"github.com/Travix-International/appix/config"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type VersionCommand struct {
	*Command
}

func (cmd *VersionCommand) Register(app *kingpin.Application, config config.Config) {
	app.Command("version", "Displays version information").
		Action(func(parseContext *kingpin.ParseContext) error {
			log.Printf("Version: %s", config.Version)
			log.Printf("Hash: %s", config.GitHash)
			log.Printf("Build date: %s", config.ParsedBuildDate)

			return nil
		}).
		Alias("ver").
		Alias("v")
}
