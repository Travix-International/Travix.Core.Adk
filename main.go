package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	modelsConfig "github.com/Travix-International/Travix.Core.Adk/models/config"
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

func makeConfig() modelsConfig.Config {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	directoryPath := filepath.Join(user.HomeDir, ".appix")

	config := modelsConfig.Config{
		Version:         version,
		BuildDate:       buildDate,
		ParsedBuildDate: parsedBuildDate,
		GitHash:         gitHash,
		Verbose:         verbose, // @TODO: verify if it works as --verbose
		DevFileName:     ".appixDevSettings",
		CatalogURIs:     catalogURIs,
		TargetEnv:       targetEnv,
		LocalFrontend:   localFrontend,

		DirectoryPath: directoryPath,
		AuthFilePath:  filepath.Join(directoryPath, "auth.json"),

		DeveloperProfileUrl: travixDeveloperProfileUrl,

		FirebaseApiKey:            travixFirebaseApiKey,
		FirebaseAuthDomain:        travixFirebaseAuthDomain,
		FirebaseDatabaseUrl:       travixFirebaseDatabaseUrl,
		FirebaseStorageBucket:     travixFirebaseStorageBucket,
		FirebaseMessagingSenderId: travixFirebaseMessagingSenderId,

		AuthServerPort: "7001",
	}

	return config
}

func main() {
	var err error
	parsedBuildDate, err = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)
	if err != nil {
		log.Fatal(err)
	}

	// App
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&targetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&verbose)

	app.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&localFrontend)

	// Context
	config := makeConfig()
	context := modelsContext.Context{
		App:    app,
		Config: config,
	}

	// Commands
	cmdVersion.Register(context)
	cmdLogin.Register(context)
	cmdWhoami.Register(context)
	cmdInit.Register(context)
	cmdPush.Register(context)
	cmdSubmit.Register(context)
	cmdWatch.Register(context)

	// kingpin config
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
