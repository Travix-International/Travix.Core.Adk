package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix"
	"github.com/Travix-International/appix/config"
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
	targetEnv        = "prod"
	verbose          = false
	localFrontend    = false
	maxRetryAttempts = 5
	maxTimeoutValue  = 5 // seconds
)

func main() {
	parsedBuildDate, _ = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)

	args := appix.GlobalArgs{}

	// App
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("catalog", "Specify the catalog to use (local, dev, staging, prod)").
		Short('c').
		Default("prod").
		EnumVar(&args.TargetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&args.Verbose)

	// Context
	config := makeConfig()

	appix.RegisterInit(app, config, &args)
	appix.RegisterLogin(app, config, &args)
	appix.RegisterPush(app, config, &args)
	appix.RegisterSubmit(app, config, &args)
	appix.RegisterVersion(app, config)
	appix.RegisterWatch(app, config, &args)
	appix.RegisterWhoami(app, config, &args)

	// kingpin config
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func makeConfig() config.Config {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	directoryPath := filepath.Join(user.HomeDir, ".appix")

	config := config.Config{
		Version:         version,
		BuildDate:       buildDate,
		ParsedBuildDate: parsedBuildDate,
		GitHash:         gitHash,
		CatalogURIs:     catalogURIs,

		DirectoryPath: directoryPath,
		AuthFilePath:  filepath.Join(directoryPath, "auth.json"),

		DeveloperProfileUrl: travixDeveloperProfileUrl,

		FirebaseApiKey:            travixFirebaseApiKey,
		FirebaseAuthDomain:        travixFirebaseAuthDomain,
		FirebaseDatabaseUrl:       travixFirebaseDatabaseUrl,
		FirebaseStorageBucket:     travixFirebaseStorageBucket,
		FirebaseMessagingSenderId: travixFirebaseMessagingSenderId,

		AuthServerPort:   "7001",
		MaxRetryAttempts: maxRetryAttempts,
		MaxTimeoutValue:  maxTimeoutValue,
	}

	return config
}
