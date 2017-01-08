package appix

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/appcatalog"
	"github.com/Travix-International/appix/config"
)

// RegisterSubmit registers the 'submit' command.
func RegisterSubmit(app *kingpin.Application, config config.Config, args *GlobalArgs) {
	const submitTemplateURI = "%s/files/publish/%s"

	var (
		appPath  string // path to the App folder
		noVerify bool
	)

	command := app.Command("submit", "Submits the App for review.").
		Action(func(parseContext *kingpin.ParseContext) error {
			environment := args.TargetEnv

			if environment == "" {
				environment = "dev"
			}

			appPath, appName, appManifestFile, err := prepareAppUpload(appPath, noVerify)

			if err != nil {
				log.Println("Could not prepare the app folder for uploading")
				return err
			}

			zapFile, err := createZapPackage(appPath, args.Verbose)

			if err != nil {
				log.Println("Could not create zap package!")
				return err
			}

			log.Printf("Run submit for App '%s', env '%s', path '%s'\n", appName, environment, appPath)

			rootURI := config.CatalogURIs[environment]
			submitURI := fmt.Sprintf(submitTemplateURI, rootURI, appName)

			acceptanceQueryURLPath, err := appcatalog.SubmitToCatalog(submitURI, appManifestFile, zapFile, args.Verbose, config)

			if err != nil {
				return err
			}

			log.Println("App has been submitted successfully.")

			if acceptanceQueryURLPath != "" {
				log.Println("You can use the following query URL to get this particular version of this app:")
				log.Printf("\t%s%s\n", rootURI, acceptanceQueryURLPath)
			}

			return nil
		})

	command.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&appPath)

	command.Flag("noVerify", "Appix won't run the tests.").
		Default("false").
		BoolVar(&noVerify)
}
