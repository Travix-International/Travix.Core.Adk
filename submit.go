package appix

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/appcatalog"
	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

// RegisterSubmit registers the 'submit' command.
func RegisterSubmit(app *kingpin.Application, config config.Config, args *GlobalArgs, logger *appixLogger.Logger) {
	const submitTemplateURI = "%s/apps/%s/submit"

	var (
		appPath string // path to the App folder
		timeout int
	)

	command := app.Command("submit", "Submits the App for review.").
		Action(func(parseContext *kingpin.ParseContext) error {
			environment := args.TargetEnv

			if environment == "" {
				environment = "dev"
			}

			appPath, appName, appManifestFile, err := prepareAppUpload(appPath)

			if err != nil {
				logger.AddMessageToQueue(appixLogger.LoggerNotification{
					Level:    "error",
					Message:  fmt.Sprintf("Could not prepare the app folder for uploading: %s", err.Error()),
					LogEvent: "AppixSubmit",
				})
				return err
			}

			zapFile, err := createZapPackage(appPath, args.Verbose)

			if err != nil {
				logger.AddMessageToQueue(appixLogger.LoggerNotification{
					Level:    "error",
					Message:  fmt.Sprintf("Could not create zap package: %s", err.Error()),
					LogEvent: "AppixSubmit",
				})
				return err
			}

			log.Printf("Run submit for App '%s', env '%s', path '%s'\n", appName, environment, appPath)

			rootURI := config.CatalogURL
			submitURI := fmt.Sprintf(submitTemplateURI, rootURI, appName)

			acceptanceQueryURLPath, err := appcatalog.SubmitToCatalog(submitURI, timeout, appManifestFile, zapFile, args.Verbose, config, logger)

			if err != nil {
				logger.AddMessageToQueue(appixLogger.LoggerNotification{
					Level:    "error",
					Message:  fmt.Sprintf("Could not submit manifest to App Catalog: %s", err.Error()),
					LogEvent: "AppixSubmit",
				})
				return err
			}

			logger.AddMessageToQueue(appixLogger.LoggerNotification{
				Level:    "error",
				Message:  "App has been submitted successfully.",
				LogEvent: "AppixSubmit",
			})

			if acceptanceQueryURLPath != "" {
				log.Println("You can use the following query URL to get this particular version of this app:")
				log.Printf("\t%s%s\n", rootURI, acceptanceQueryURLPath)
			}

			return nil
		})

	command.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&appPath)

	command.Flag("timeout", "Set the maximum timeout for the request").
		Default("10").
		IntVar(&timeout)
}
