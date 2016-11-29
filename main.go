package main

import (
	"context"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	config "github.com/Travix-International/Travix.Core.Adk/models/config"
	modelsContext "github.com/Travix-International/Travix.Core.Adk/models/context"

	cmdInit "github.com/Travix-International/Travix.Core.Adk/cmd/init"
	cmdLogin "github.com/Travix-International/Travix.Core.Adk/cmd/login"
	cmdPush "github.com/Travix-International/Travix.Core.Adk/cmd/push"
	cmdSubmit "github.com/Travix-International/Travix.Core.Adk/cmd/submit"
	cmdVersion "github.com/Travix-International/Travix.Core.Adk/cmd/version"
	cmdWatch "github.com/Travix-International/Travix.Core.Adk/cmd/watch"
	cmdWhoami "github.com/Travix-International/Travix.Core.Adk/cmd/whoami"
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
		".appixDevSettings",
		catalogURIs,
		targetEnv,
		localFrontend,

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

const CONTEXTKEY int = 1

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

	// Context
	ctx := context.WithValue(context.Background(), CONTEXTKEY, modelsContext.Context{app, &config})

	// modelsContext.Context{
	// 	App:    app,
	// 	Config: &config,
	// }

	// Commands
	cmdInit.Register(ctx)
	cmdLogin.Register(ctx)
	cmdPush.Register(ctx)
	cmdSubmit.Register(ctx)
	cmdVersion.Register(ctx)
	cmdWatch.Register(ctx)
	cmdWhoami.Register(ctx)

	// kingpin config
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
