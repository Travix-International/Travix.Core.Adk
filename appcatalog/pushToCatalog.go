package appcatalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Travix-International/appix/config"
)

const (
	timeout = time.Duration(10 * time.Second)
)

// PushToCatalog pushes the specified app to the AppCatalog.
func PushToCatalog(pushURI string, timeout int, appManifestFile string, verbose bool, config config.Config) (uploadURI string, err error) {
	var req *http.Request
	files := map[string]string{
		"manifest": appManifestFile,
	}

	for attempt := 1; attempt <= config.MaxRetryAttempts; attempt++ {
		if req, err = prepare(pushURI, files, config, verbose); err != nil {
			return "", err
		}

		log.Printf("Pushing files to catalog. Attempt %v of %v\n", attempt, config.MaxRetryAttempts)

		if uploadURI, err = doPush(req, time.Duration(timeout)*time.Second, verbose); err == nil {
			break
		}

		log.Printf("An error occured when trying to push the application: %s\n", err.Error())

		if attempt < config.MaxRetryAttempts {
			wait := math.Pow(2, float64(attempt-1)) * 1000
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

	return
}

func doPush(req *http.Request, maxTimeoutValue time.Duration, verbose bool) (uploadURI string, err error) {
	client := &http.Client{
		Timeout: maxTimeoutValue,
	}

	res, err := client.Do(req)

	if err != nil {
		log.Println("Call to App Catalog failed.")
		return "", err
	}

	if verbose {
		logServerResponse(res)
	}

	if res.StatusCode == 401 || res.StatusCode == 403 {
		log.Printf("You are not authorized to push the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", res.StatusCode)
		return "", fmt.Errorf("Authentication error")
	}

	if res.StatusCode == 504 || res.StatusCode == 408 {
		log.Printf("The AppCatalog request timed out (status code %v)", res.StatusCode)
		return "", fmt.Errorf("Timeout error")
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Println("Error reading response from App Catalog.")
		return "", err
	}

	type PushResponse struct {
		Links    map[string]string `json:"links"`
		Messages []string          `json:"messages"`
	}

	var responseObject PushResponse
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		if verbose {
			log.Println(err)
		}

		return "", err
	}

	log.Printf("App Catalog returned status code %v. Response details:\n", res.StatusCode)

	for _, line := range responseObject.Messages {
		log.Printf("\t%v\n", line)
	}

	if res.StatusCode == http.StatusOK {
		log.Println("App has been pushed successfully.")
	} else {
		return "", fmt.Errorf("Push failed, App Catalog returned status code %v", res.StatusCode)
	}

	return responseObject.Links["upload"], nil
}
