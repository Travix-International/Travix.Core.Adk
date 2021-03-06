package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix"
	"github.com/Travix-International/appix/appixLogger"
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
	travixFirebaseRefreshTokenUrl   string
	travixDeveloperProfileUrl       string
	travixLoggerUrl                 string
)

// Although these are configuration values, they're not exposed to the public and are therefore kept internally.
var (
	catalogURIs = map[string]string{
		"local":   "http://localhost:5000",
		"dev":     "https://appcatalog.dev.travix.com",
		"staging": "https://appcatalog.stg.travix.com",
		"prod":    "https://appcatalog.travix.com",
	}
	maxRetryAttempts = 5
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
	app.Flag("catalog-url", "Specify the catalog URL to use").
		Short('u').
		StringVar(&args.CatalogURL)
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&args.Verbose)

	app.Parse(os.Args[1:])

	if args.CatalogURL == "" {
		args.CatalogURL = catalogURIs[args.TargetEnv]
	}

	// Context
	config := makeConfig(args.CatalogURL)

	// appixLogger
	logger := appixLogger.NewAppixLogger(config.TravixLoggerUrl)
	logger.Start()
	defer logger.Stop()

	appix.RegisterInit(app, config, &args)
	appix.RegisterLogin(app, config, &args)
	appix.RegisterPush(app, config, &args, logger)
	appix.RegisterSubmit(app, config, &args, logger)
	appix.RegisterVersion(app, config)
	appix.RegisterWatch(app, config, &args, logger)
	appix.RegisterWhoami(app, config, &args, logger)

	app.Parse(os.Args[1:])
}

func makeConfig(catalogUrl string) config.Config {
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
		CatalogURL:      catalogUrl,

		DirectoryPath: directoryPath,
		AuthFilePath:  filepath.Join(directoryPath, "auth.json"),

		DeveloperProfileUrl: travixDeveloperProfileUrl,

		FirebaseApiKey:            travixFirebaseApiKey,
		FirebaseAuthDomain:        travixFirebaseAuthDomain,
		FirebaseDatabaseUrl:       travixFirebaseDatabaseUrl,
		FirebaseStorageBucket:     travixFirebaseStorageBucket,
		FirebaseMessagingSenderId: travixFirebaseMessagingSenderId,
		FirebaseRefreshTokenUrl:   travixFirebaseRefreshTokenUrl,
		TravixLoggerUrl:           travixLoggerUrl,

		AuthServerPort:   "7001",
		MaxRetryAttempts: maxRetryAttempts,
	}

	return config
}
