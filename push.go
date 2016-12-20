package appix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/appcatalog"
	"github.com/Travix-International/appix/auth"
	"github.com/Travix-International/appix/config"
)

// PushCommand used for pushing an app during app development
type PushCommand struct {
	*Command
	AppPath       string // path to the App folder
	NoPolling     bool   // skip polling flag
	NoBrowser     bool   // skip opening the site in the browser
	WaitInSeconds int    // polling timeout
}

type bundleMessage struct {
	Widget string
	Output string
}

type pushPollResponse struct {
	Meta struct {
		Status   string
		Messages []bundleMessage
	}
	Links struct {
		Preview string
	}
}

const (
	pushTemplateURI    = "%s/files/push/%s?sessionid=%s"
	pollClientTimeout  = 5 * time.Second
	pollInterval       = 5 * time.Second // how often to poll status URL
	pollFinishedStatus = "FINISHED"
	pollFailedStatus   = "FAILED"
)

func (cmd *PushCommand) Register(app *kingpin.Application, config config.Config) {
	command := app.Command("push", "Push the App in the specified folder.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return cmd.Push(config)
		})

	command.Arg("appPath", "path to the App folder (default: current folder).").
		Default(".").
		ExistingDirVar(&cmd.AppPath)

	command.Flag("noPolling", "DEPRECATED: Appix won't wait for the bundling of the app to be finished.").
		Default("false").
		BoolVar(&cmd.NoPolling)

	command.Flag("noBrowser", "Appix won't open the frontend in the browser.").
		Default("false").
		BoolVar(&cmd.NoBrowser)

	command.Flag("wait", "The maximum time appix waits for the app bundling to be finished.").
		Short('w').
		Default("180").
		IntVar(&cmd.WaitInSeconds)
}

func (cmd *PushCommand) Push(config config.Config) error {

	appPath := cmd.AppPath
	pollingEnabled := !cmd.NoPolling
	openBrowser := !cmd.NoBrowser
	waitInSeconds := cmd.WaitInSeconds

	appPath, appName, appManifestFile, err := prepareAppUpload(cmd.AppPath)

	if err != nil {
		log.Println("Could not prepare the app folder for uploading")
		return err
	}

	zapFile, err := createZapPackage(appPath, cmd.Verbose)

	if err != nil {
		log.Println("Could not create zap package.")
		return err
	}

	sessionID, err := getSessionID(appPath, cmd.Verbose)

	if err != nil {
		log.Println("Could not get the session id.")
		return err
	}

	log.Printf("Run push for App '%s', path '%s'\n", appName, appPath)

	rootURI := config.CatalogURIs[cmd.TargetEnv]
	pushURI := fmt.Sprintf(pushTemplateURI, rootURI, appName, sessionID)

	uploadURI, err := pushToCatalog(pushURI, appManifestFile, cmd.Verbose, config)

	if err != nil {
		log.Println("Error during pushing the manifest to the App Catalog.")
		return err
	}

	if cmd.LocalFrontend {
		log.Println("Ignoring URL and substituting local front-end URL instead.")
		reg, err := regexp.Compile(`(https?:\/\/.*)(\/.*)`)
		if err != nil {
			log.Println(err)
			return err
		}

		uploadURI = reg.ReplaceAllString(uploadURI, "http://localhost:3001$2")
	}

	log.Println("Frontend upload url:", uploadURI)

	pollURI, err := uploadToFrontend(uploadURI, zapFile, appName, sessionID, cmd.Verbose)

	log.Println("Frontend upload poll uri:", pollURI)

	if err != nil {
		log.Println("Error. during uploading package to the frontend")
		return err
	}

	if pollingEnabled {
		doPolling(pollURI, waitInSeconds, openBrowser, cmd.Verbose)
	} else {
		log.Println("Polling not enabled")
		log.Println("NOTE: The --noPolling will be removed in a future version.")
		log.Println("If you want to prevent appix from opening the frontend in the browser, use the --noBrowser flag.")
	}

	if cmd.Verbose {
		log.Println("Push command has completed")
	}

	return nil
}

func doPolling(pollURI string, waitInSeconds int, openBrowser bool, verbose bool) {
	if verbose {
		log.Println("Entering polling routine")
	}
	quit := make(chan interface{}, 1)
	defer close(quit)

	progressMonitor := verifyProgress(pollURI, quit)
	wait := time.Duration(waitInSeconds) * time.Second

	select {
	case statusResponse, ok := <-progressMonitor:
		if !ok {
			break
		}

		log.Printf("Server output for the app bundling:")
		for _, message := range statusResponse.Meta.Messages {
			log.Printf("Widget: %s", message.Widget)
			log.Printf("Output: %s", message.Output)
		}

		if statusResponse.Meta.Status == pollFinishedStatus {
			log.Printf("App successfully pushed. The frontend for this development session is at %s", statusResponse.Links.Preview)
			if openBrowser {
				openURL(statusResponse.Links.Preview)
			}
		} else {
			log.Printf("App push failed.")
		}

		close(progressMonitor)

	case <-time.After(wait):
		quit <- true // send a cancel signal to progressMonitor
		log.Printf("Operation timed out after %s", wait)
	}
}

func verifyProgress(pollURI string, quit <-chan interface{}) chan pushPollResponse {
	done := make(chan pushPollResponse, 1)
	go func() {
		var statusResponse pushPollResponse
		timeout := time.Duration(pollClientTimeout)
		client := http.Client{Timeout: timeout}

		for {
			// check if operation should be cancelled
			select {
			case <-quit:
				return
			default:
			}

			resp, err := client.Get(pollURI)
			if err != nil {
				log.Println("Error during polling the bundling status.")
				log.Println(err)
				close(done)
				return
			}

			err = json.NewDecoder(resp.Body).Decode(&statusResponse)
			resp.Body.Close()

			if err != nil {
				log.Println("Error. during parsing poll status result")
				bodyData, _ := ioutil.ReadAll(resp.Body)
				if bodyData != nil {
					log.Println(bodyData)
				}
				close(done)
				return
			}

			log.Printf("Pushing to the website to the development environment, status: [%s]", statusResponse.Meta.Status)

			if statusResponse.Meta.Status == pollFinishedStatus || statusResponse.Meta.Status == pollFailedStatus {
				done <- statusResponse
				break
			}

			time.Sleep(pollInterval)
		}
	}()
	return done
}

func pushToCatalog(pushURI string, appManifestFile string, verbose bool, config config.Config) (uploadURI string, err error) {
	// To the App Catalog we have to POST the manifest in a multipart HTTP form.
	// When doing the push, it'll only contain a single file, the manifest.
	files := map[string]string{
		"manifest": appManifestFile,
	}

	if verbose {
		log.Println("Posting the app manifest to the App Catalog overlay: " + pushURI)
	}

	request, err := appcatalog.CreateMultiFileUploadRequest(pushURI, files, nil, verbose)

	if err != nil {
		log.Println("Creating the HTTP request failed.")
		return "", err
	}

	token, err := auth.LoadAuthToken(config)

	if err == nil {
		request.Header.Set("Authorization", token.TokenType+" "+token.IdToken)
	} else {
		log.Println("WARNING: You are not logged in. In a future version authentication will be mandatory.\nYou can log in using \"appix login\".")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Call to App Catalog failed.")
		return "", err
	}

	if response.StatusCode == 401 || response.StatusCode == 403 {
		log.Printf("You are not authorized to push the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", response.StatusCode)
		return "", fmt.Errorf("Authentication error")
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("Error reading response from App Catalog.")
		return "", err
	}

	type PushResponse struct {
		Links    map[string]string `json:"links"`
		Messages []string          `json:"messages"`
	}

	var responseObject PushResponse
	err = json.Unmarshal(responseBody, &responseObject)
	if err != nil {
		if verbose {
			log.Println(err)
		}

		return "", err
	}

	log.Printf("App Catalog returned status code %v. Response details:\n", response.StatusCode)

	for _, line := range responseObject.Messages {
		log.Printf("\t%v\n", line)
	}

	if response.StatusCode == http.StatusOK {
		log.Println("App has been pushed successfully.")
	} else {
		return "", fmt.Errorf("Push failed, App Catalog returned status code %v", response.StatusCode)
	}

	return responseObject.Links["upload"], nil
}

func uploadToFrontend(uploadURI string, zapFile string, appName string, sessionID string, verbose bool) (frontendURI string, err error) {
	files := map[string]string{
		"file": zapFile,
	}

	params := map[string]string{
		"name": appName,
	}

	if verbose {
		log.Println("Uploading the app to the Express frontend: " + uploadURI)
		log.Println("Creating multi-file upload request")
	}

	request, err := appcatalog.CreateMultiFileUploadRequest(uploadURI, files, params, verbose)

	if err != nil {
		log.Println("Creating the HTTP request failed.")
		return "", err
	}

	if verbose {
		log.Println("Multi-file upload request created, proceeding to call front-end")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Call to the Express frontend failed.")
		return "", err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from the fronted.")
		return "", err
	}

	var responseObject map[string]map[string]string
	err = json.Unmarshal(responseBody, &responseObject)
	if err != nil {
		if verbose {
			log.Println(err)
		}

		return "", err
	}

	log.Printf("Express frontend returned status code %v.", response.StatusCode)

	if response.StatusCode == http.StatusOK {
		log.Println("The app has been uploaded to the frontend successfully.")
	} else {
		return "", fmt.Errorf("Uploading failed, the frontend returned status code %v", response.StatusCode)
	}

	// The frontend returns a link which can be used to poll the upload status.
	// {
	//   "links": {
	//     "progress": "https://fireball-dev.travix.com/upload/progress?sessionId=123`"
	//   }
	// }
	return responseObject["links"]["progress"], nil
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
