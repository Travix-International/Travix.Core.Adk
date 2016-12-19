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
	targetEnv     = "prod"
	verbose       = false
	localFrontend = false
)

func main() {
	parsedBuildDate, _ = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)

	command := &appix.Command{}

	// App
	app := kingpin.New("appix", "App Developer Kit for the Travix Fireball infrastructure.")

	app.Flag("cat", "Specify the catalog to use (local, dev, staging, prod)").
		Default("prod").
		EnumVar(&command.TargetEnv, "local", "dev", "staging", "prod")
	app.Flag("verbose", "Verbose mode.").
		Short('v').
		BoolVar(&command.Verbose)

	app.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&command.LocalFrontend)

	// Context
	config := makeConfig()

	commands := [...]appix.Registrable{
		&appix.InitCommand{Command: command},
		&appix.LoginCommand{Command: command},
		&appix.PushCommand{Command: command},
		&appix.SubmitCommand{Command: command},
		&appix.VersionCommand{Command: command},
		&appix.WatchCommand{Command: command},
		&appix.WhoamiCommand{Command: command},
	}

	for _, c := range commands {
		c.Register(app, config)
	}

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
		DevFileName:     ".appixDevSettings",
		CatalogURIs:     catalogURIs,

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
