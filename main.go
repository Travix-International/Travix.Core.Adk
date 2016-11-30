package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	config "github.com/Travix-International/Travix.Core.Adk/models/config"
)

// Version numbers passed by build flags
var (
	version         string
	buildDate       string
	parsedBuildDate time.Time
	gitHash         string

	travixFirebaseApiKey            string
	travixFirebaseAuthDomain        string
	travixFirebaseDatabaseUrl       string
	travixFirebaseStorageBucket     string
	travixFirebaseMessagingSenderId string
	travixDeveloperProfileUrl       string
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

func makeConfig() config.Config {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	directoryPath := filepath.Join(user.HomeDir, ".appix")

	config := config.New(
		version,
		buildDate,
		parsedBuildDate,
		gitHash,
		verbose, // @TODO: verify if it works as --verbose
		catalogURIs,
		targetEnv,
		localFrontend,
		".appixDevSettings",

		directoryPath,
		filepath.Join(directoryPath, "auth.json"),

		travixDeveloperProfileUrl,

		travixFirebaseApiKey,
		travixFirebaseAuthDomain,
		travixFirebaseDatabaseUrl,
		travixFirebaseStorageBucket,
		travixFirebaseMessagingSenderId,

		"7001",
	)

	return config
}

func main() {
	var err error
	parsedBuildDate, err = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)
	if err != nil {
		log.Fatal(err)
	}

	// Config
	config := makeConfig()

	// App
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&config.TargetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&config.Verbose)

	app.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&config.LocalFrontend)

	// Commands
	registerInit(app, &config)
	registerLogin(app, &config)
	registerPush(app, &config)
	registerSubmit(app, &config)
	registerVersion(app, &config)
	registerWatch(app, &config)
	registerWhoAmI(app, &config)

	// kingpin config
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
