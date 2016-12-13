package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/config"
	"github.com/Travix-International/Travix.Core.Adk/lib/context"

	cmd "github.com/Travix-International/Travix.Core.Adk/lib/cmd"
	cmdInit "github.com/Travix-International/Travix.Core.Adk/lib/cmd/init"
	cmdLogin "github.com/Travix-International/Travix.Core.Adk/lib/cmd/login"
	cmdPush "github.com/Travix-International/Travix.Core.Adk/lib/cmd/push"
	cmdSubmit "github.com/Travix-International/Travix.Core.Adk/lib/cmd/submit"
	cmdVersion "github.com/Travix-International/Travix.Core.Adk/lib/cmd/version"
	cmdWatch "github.com/Travix-International/Travix.Core.Adk/lib/cmd/watch"
	cmdWhoami "github.com/Travix-International/Travix.Core.Adk/lib/cmd/whoami"
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
	var err error
	parsedBuildDate, err = time.Parse("Mon.January.2.2006.15:04:05.-0700.MST", buildDate)
	if err != nil {
		log.Fatal(err)
	}

	command := &cmd.Command{}

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
	context := context.Context{
		App:    app,
		Config: config,
	}

	commands := [...]cmd.Registrable{
		&cmdInit.InitCommand{Command: command},
		&cmdLogin.LoginCommand{Command: command},
		&cmdPush.PushCommand{Command: command},
		&cmdSubmit.SubmitCommand{Command: command},
		&cmdVersion.VersionCommand{Command: command},
		&cmdWatch.WatchCommand{Command: command},
		&cmdWhoami.WhoamiCommand{Command: command},
	}

	for _, c := range commands {
		c.Register(context)
	}

	// kingpin config
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func makeConfig() *config.Config {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	directoryPath := filepath.Join(user.HomeDir, ".appix")

	config := &config.Config{
		Version:         version,
		BuildDate:       buildDate,
		ParsedBuildDate: parsedBuildDate,
		GitHash:         gitHash,
		DevFileName:     ".appixDevSettings",
		IgnoreFileName:  ".appixignore",
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
