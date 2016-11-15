package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// PushCommand used for pushing an app during app development
type PushCommand struct {
	appPath       string // path to the App folder
	noPolling     bool   // skip polling flag
	noBrowser     bool   // skip opening the site in the browser
	waitInSeconds int    // polling timeout
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

func configurePushCommand(app *kingpin.Application) {
	cmd := &PushCommand{}
	appCmd := app.Command("push", "Push the App in the specified folder.").
		Action(cmd.push)
	appCmd.Arg("appPath", "path to the App folder (default: current folder).").
		Default(".").
		ExistingDirVar(&cmd.appPath)
	appCmd.Flag("noPolling", "DEPRECATED: Appix won't wait for the bundling of the app to be finished.").
		Default("false").
		BoolVar(&cmd.noPolling)
	appCmd.Flag("noBrowser", "Appix won't open the frontend in the browser.").
		Default("false").
		BoolVar(&cmd.noBrowser)
	appCmd.Flag("wait", "The maximum time appix waits for the app bundling to be finished.").
		Short('w').
		Default("180").
		IntVar(&cmd.waitInSeconds)
}

func (cmd *PushCommand) push(context *kingpin.ParseContext) error {
	appPath := cmd.appPath
	pollingEnabled := !cmd.noPolling
	openBrowser := !cmd.noBrowser
	waitInSeconds := cmd.waitInSeconds

	appPath, appName, appManifestFile, err := prepareAppUpload(cmd.appPath)

	if err != nil {
		log.Println("Could not prepare the app folder for uploading")
		return err
	}

	zapFile, err := createZapPackage(appPath)

	if err != nil {
		log.Println("Could not create zap package.")
		return err
	}

	sessionID, err := getSessionID(appPath)

	if err != nil {
		log.Println("Could not get the session id.")
		return err
	}

	log.Printf("Run push for App '%s', path '%s'\n", appName, appPath)

	rootURI := catalogURIs[targetEnv]
	pushURI := fmt.Sprintf(pushTemplateURI, rootURI, appName, sessionID)

	uploadURI, err := pushToCatalog(pushURI, appManifestFile)

	if localFrontend {
		log.Println("Ignoring URL and substituting local front-end URL instead.")
		reg, err := regexp.Compile(`(https?:\/\/.*)(\/.*)`)
		if err != nil {
			log.Println(err)
			return err
		}

		uploadURI = reg.ReplaceAllString(uploadURI, "http://localhost:3001$2")
	}

	if err != nil {
		log.Println("Error during pushing the manifest to the App Catalog.")
		return err
	}

	log.Println("Frontend upload url:", uploadURI)

	pollURI, err := uploadToFrontend(uploadURI, zapFile, appName, sessionID)

	log.Println("Frontend upload poll uri:", pollURI)

	if err != nil {
		log.Println("Error. during uploading package to the frontend")
		return err
	}

	if pollingEnabled {
		doPolling(pollURI, waitInSeconds, openBrowser)
	} else {
		log.Println("Polling not enabled")
		log.Println("NOTE: The --noPolling will be removed in a future version.")
		log.Println("If you want to prevent appix from opening the frontend in the browser, use the --noBrowser flag.")
	}

	if verbose {
		log.Println("Push command has completed")
	}

	return nil
}

func doPolling(pollURI string, waitInSeconds int, openBrowser bool) {
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
				openWebsite(statusResponse.Links.Preview)
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
				log.Println("Error. during polling push to the frontend")
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

func pushToCatalog(pushURI string, appManifestFile string) (uploadURI string, err error) {
	// To the App Catalog we have to POST the manifest in a multipart HTTP form.
	// When doing the push, it'll only contain a single file, the manifest.
	files := map[string]string{
		"manifest": appManifestFile,
	}

	if verbose {
		log.Println("Posting the app manifest to the App Catalog overlay: " + pushURI)
	}

	request, err := createMultiFileUploadRequest(pushURI, files, nil)

	if err != nil {
		log.Println("Creating the HTTP request failed.")
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Call to App Catalog failed.")
		return "", err
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

func uploadToFrontend(uploadURI string, zapFile string, appName string, sessionID string) (frontendURI string, err error) {
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

	request, err := createMultiFileUploadRequest(uploadURI, files, params)

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
func getSessionID(appPath string) (string, error) {
	settings, err := readDevelopmentSettings(appPath)

	if err != nil {
		settings, err = getDefaultDevelopmentSettings()

		if err != nil {
			log.Println("Couldn't create new development settings.")
			return "", err
		}

		err = writeDevelopmentSettings(appPath, settings)

		if err != nil {
			log.Println("Could not save new development settings file.")
			return "", err
		}
	}

	return settings.SessionID, nil
}
