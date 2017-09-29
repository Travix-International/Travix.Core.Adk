package appix

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/auth"
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
		timeout       int    // request timeout
		localFrontend bool   // true if we open the local frontend instead of the dev server
	)

	command := app.Command("push", "Push the App in the specified folder.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return push(config, appPath, noBrowser, waitInSeconds, timeout, localFrontend, args, logger)
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

	command.Flag("timeout", "Set the maximum timeout for the request").
		Default("10").
		IntVar(&timeout)
}

func push(config config.Config, appPath string, noBrowser bool, wait int, timeout int, localFrontend bool, args *GlobalArgs, logger *appixLogger.Logger) error {
	appPath, appName, _, err := prepareAppUpload(appPath)

	if err != nil {
		log.Println("Could not prepare the app folder for uploading")
		return err
	}

	zapFile, err := createZapPackage(appPath, args.Verbose)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Could not create zap package: %s", err.Error()),
			LogEvent: "AppixPush",
		})
		return err
	}

	devSettings, err := getDevSettings(appPath, args.Verbose)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Could not get the session id: %s", err.Error()),
			LogEvent: "AppixPush",
		})
		return err
	}

	log.Printf("Run push for App '%s', path '%s'\n", appName, appPath)

	// load auth token
	tb, err := auth.LoadAuthToken(config, logger)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("You are not logged in.\nYou can sign in by using 'appix login'. Error message: %s", err.Error()),
			LogEvent: "AppixPush",
		})
		return err
	}

	// request the upload url
	var uploadObject *SignedUploadURL
	for attempt := 1; attempt <= config.MaxRetryAttempts; attempt++ {
		uploadObject, err = RetrieveUploadURL(config.TravixUploadUrl, tb.IdToken, appName, devSettings.SessionID)
		if err == nil {
			break
		}
		log.Printf("An error ocurred while retrieving an upload URL for your app. Retry attempt %d of %d \n", attempt, config.MaxRetryAttempts)
		if attempt < config.MaxRetryAttempts {
			wait := math.Pow(2, float64(attempt-1)) * 1000
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Unable to retrieve an upload URL for your app. Error message: %s", err.Error()),
			LogEvent: "AppixPush",
		})
		return err
	}

	if err = uploadObject.UploadResource(zapFile); err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Unable to upload your app. Error message: %s", err.Error()),
			LogEvent: "AppixPush",
		})
		return err
	}

	logger.AddMessageToQueue(appixLogger.LoggerNotification{
		Level:    "log",
		Message:  fmt.Sprintf("The widget was succesfully uploaded"),
		LogEvent: "AppixPush",
	})

	// if localFrontend {
	// 	uploadURI, err = replaceURLSchemaAndDomain(uploadURI, "http://localhost:3001")
	// } else if devSettings.DevServerOverride != "" {
	// 	uploadURI, err = replaceURLSchemaAndDomain(uploadURI, devSettings.DevServerOverride)
	// }

	// if err != nil {
	// 	logger.AddMessageToQueue(appixLogger.LoggerNotification{
	// 		Level:    "error",
	// 		Message:  fmt.Sprintf("Error during trying to overriding the Dev server URL: %s", err.Error()),
	// 		LogEvent: "AppixPush",
	// 	})
	// 	return err
	// }

	// log.Println("Frontend upload url:", uploadURI)

	// pollURI, err := appcatalog.UploadToFrontend(uploadURI, zapFile, appName, devSettings.SessionID, args.Verbose)

	// if err != nil {
	// 	logger.AddMessageToQueue(appixLogger.LoggerNotification{
	// 		Level:    "error",
	// 		Message:  fmt.Sprintf("Error during uploading package to the frontend: %s", err.Error()),
	// 		LogEvent: "AppixPush",
	// 	})
	// 	return err
	// }

	// log.Println("Frontend upload poll uri:", pollURI)

	// appcatalog.PollUntilDone(pollURI, wait, !noBrowser, args.Verbose, openURL)

	if args.Verbose {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "log",
			Message:  "Push command has completed",
			LogEvent: "AppixPush",
		})
	}

	return nil
}

// replaceURLSchemaAndDomain replaces the schema and domain part in the uri, with this we can override which dev server we are pushing to.
func replaceURLSchemaAndDomain(uri string, replace string) (string, error) {
	log.Printf("Ignoring URL and substituting it with %s.", replace)
	reg, err := regexp.Compile(`(https?:\/\/.*)(\/.*)`)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return reg.ReplaceAllString(uri, replace+"$2"), nil
}

// getSessionID gets the current session id. If there is an existing one in the folder, it uses that, otherwise it creates a new one.
func getDevSettings(appPath string, verbose bool) (*DevelopmentSettings, error) {
	s, err := readDevelopmentSettings(appPath, verbose)

	if err != nil {
		s, err = getDefaultDevelopmentSettings()

		if err != nil {
			log.Println("Couldn't create new development settings.")
			return nil, err
		}

		err = writeDevelopmentSettings(appPath, s, verbose)

		if err != nil {
			log.Println("Could not save new development settings file.")
			return nil, err
		}
	}

	return s, nil
}
