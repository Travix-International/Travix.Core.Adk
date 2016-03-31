package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Although these are configuration values, they're not exposed to the public and are therefore kept internally.
var (
	catalogURIs = map[string]string{
		"local":   "http://localhost:52426",
		"dev":     "https://appcatalog.dev.travix.com",
		"staging": "https://appcatalog.staging.travix.com",
		"prod":    "https://appcatalog.travix.com",
	}
	targetEnv = "dev"
	verbose   = false
)

func main() {
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&targetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&verbose)

	configurePublishCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
