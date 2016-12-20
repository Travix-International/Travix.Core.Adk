package appix

import (
	"log"

	"github.com/Travix-International/appix/config"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// RegisterVersion registers the 'version' command.
func RegisterVersion(app *kingpin.Application, config config.Config) {
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
