package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Version numbers passed by build flags
var (
	version         string
	buildDate       string
	parsedBuildDate time.Time
	gitHash         string
)

// Although these are configuration values, they're not exposed to the public and are therefore kept internally.
var (
	catalogURIs = map[string]string{
		"local":   "http://localhost:5000",
		"dev":     "https://appcatalog.development.travix.com",
		"staging": "https://appcatalog.staging.travix.com",
		"prod":    "https://appcatalog.travix.com",
	}
	targetEnv     = "prod"
	verbose       = false
	localFrontend = false
)

func main() {
	var err error
	parsedBuildDate, err = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)
	if err != nil {
		log.Fatal(err)
	}

	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&targetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&verbose)

	app.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&localFrontend)

	configureInitCommand(app)
	configureLoginCommand(app)
	configurePushCommand(app)
	configureSubmitCommand(app)
	configureVersionCommand(app)
	configureWatchCommand(app)
	configureWhoamiCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
