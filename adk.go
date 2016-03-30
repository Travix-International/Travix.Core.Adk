package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	catalogURIs = map[string]string{
		"local":   "http://localhost:52426",
		"dev":     "http://appcatalog.dev.travix.com",
		"staging": "https://appcatalog.staging.travix.com",
		"prod":    "https://appcatalog.travix.com",
	}
	targetEnv = "dev"
)

func main() {
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&targetEnv, "local", "dev", "staging", "prod")

	configurePublishCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
