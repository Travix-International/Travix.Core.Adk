package appix

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/appcatalog"
	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

const (
	pushTemplateURI    = "%s/apps/%s/push?sessionid=%s"
	pollClientTimeout  = 5 * time.Second
	pollInterval       = 5 * time.Second // how often to poll status URL
	pollFinishedStatus = "FINISHED"
	pollFailedStatus   = "FAILED"
)

// RegisterPush registers the 'push' command.
func RegisterPush(app *kingpin.Application, config config.Config, args *GlobalArgs, logger *appixLogger.Logger) {
	var (
		appPath       string // path to the App folder
		noBrowser     bool   // skip opening the site in the browser
		waitInSeconds int    // polling timeout
		localFrontend bool   // true if we open the local frontend instead of the dev server
	)

	command := app.Command("push", "Push the App in the specified folder.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return push(config, appPath, noBrowser, waitInSeconds, localFrontend, args, logger)
		})

	command.Arg("appPath", "path to the App folder (default: current folder).").
		Default(".").
		ExistingDirVar(&appPath)

	command.Flag("noBrowser", "Appix won't open the frontend in the browser.").
		Default("false").
		BoolVar(&noBrowser)

	command.Flag("wait", "The maximum time appix waits for the app bundling to be finished.").
		Short('w').
		Default("180").
		IntVar(&waitInSeconds)

	command.Flag("local", "Upload to the local RWD frontend instead of the one returned by the catalog.").
		BoolVar(&localFrontend)
}

func push(config config.Config, appPath string, noBrowser bool, wait int, localFrontend bool, args *GlobalArgs, logger *appixLogger.Logger) error {
	appPath, appName, appManifestFile, err := prepareAppUpload(appPath)

	if err != nil {
		log.Println("Could not prepare the app folder for uploading")
		return err
	}

	zapFile, err := createZapPackage(appPath, args.Verbose)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Type:    "error",
			Message: fmt.Sprintf("Could not create zap package: %s", err.Error()),
			Action:  "AppixPush",
		})
		return err
	}

	sessionID, err := getSessionID(appPath, args.Verbose)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Type:    "error",
			Message: fmt.Sprintf("Could not get the session id: %s", err.Error()),
			Action:  "AppixPush",
		})
		return err
	}

	log.Printf("Run push for App '%s', path '%s'\n", appName, appPath)

	rootURI := config.CatalogURIs[args.TargetEnv]
	pushURI := fmt.Sprintf(pushTemplateURI, rootURI, appName, sessionID)

	uploadURI, err := appcatalog.PushToCatalog(pushURI, appManifestFile, args.Verbose, config, logger)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Type:    "error",
			Message: fmt.Sprintf("Error during pushing the manifest to the App Catalog: %s", err.Error()),
			Action:  "AppixPush",
		})
		return err
	}

	if localFrontend {
		log.Println("Ignoring URL and substituting local front-end URL instead.")
		reg, err := regexp.Compile(`(https?:\/\/.*)(\/.*)`)
		if err != nil {
			log.Println(err)
			return err
		}

		uploadURI = reg.ReplaceAllString(uploadURI, "http://localhost:3001$2")
	}

	log.Println("Frontend upload url:", uploadURI)

	pollURI, err := appcatalog.UploadToFrontend(uploadURI, zapFile, appName, sessionID, args.Verbose)

	log.Println("Frontend upload poll uri:", pollURI)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Type:    "error",
			Message: fmt.Sprintf("Error during uploading package to the frontend: %s", err.Error()),
			Action:  "AppixPush",
		})
		return err
	}

	appcatalog.PollUntilDone(pollURI, wait, !noBrowser, args.Verbose, openURL)

	if args.Verbose {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Type:    "log",
			Message: "Push command has completed",
			Action:  "AppixPush",
		})
	}

	return nil
}

// getSessionID gets the current session id. If there is an existing one in the folder, it uses that, otherwise it creates a new one.
func getSessionID(appPath string, verbose bool) (string, error) {
	s, err := readDevelopmentSettings(appPath, verbose)

	if err != nil {
		s, err = getDefaultDevelopmentSettings()

		if err != nil {
			log.Println("Couldn't create new development settings.")
			return "", err
		}

		err = writeDevelopmentSettings(appPath, s, verbose)

		if err != nil {
			log.Println("Could not save new development settings file.")
			return "", err
		}
	}

	return s.SessionID, nil
}
