package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Although these are configuration values, they're not exposed to the public and are therefore kept internally.
var (
	catalogURIs = map[string]string{
		"local":   "http://localhost:5000",
		"dev":     "https://appcatalog.development.travix.com",
		"staging": "https://appcatalog.staging.travix.com",
		"prod":    "https://appcatalog.travix.com",
	}
	targetEnv     = "staging"
	verbose       = false
	localFrontend = false
)

func main() {
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("dev").
		EnumVar(&targetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&verbose)

	app.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&localFrontend)

	configureInitCommand(app)
	configurePushCommand(app)
	configurePublishCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
