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

type submitResponse struct {
	Messages []string
	Links    map[string]string
}

// SubmitToCatalog submits the specified app to the AppCatalog.
func SubmitToCatalog(submitURI string, appManifestFile string, zapFile string, verbose bool, config config.Config) (acceptanceQueryURL string, err error) {
	var req *http.Request
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	if req, err = prepare(submitURI, files, config, verbose); err != nil {
		return "", err
	}

	for attempt := 1; attempt <= config.MaxRetryAttempts; attempt++ {
		log.Printf("Submitting files to App Catalog. Attempt %v of %v\n", attempt, config.MaxRetryAttempts)
		if acceptanceQueryURL, err = doSubmit(req, verbose); err == nil {
			break
		}

		if attempt < config.MaxRetryAttempts {
			wait := math.Pow(2, float64(attempt-1)) * 1000
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

	return
}

func doSubmit(req *http.Request, verbose bool) (acceptanceQueryURL string, err error) {
	client := &http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Call to App Catalog failed!")
		return "", err
	}

	if verbose {
		logServerResponse(res)
	}

	if res.StatusCode == 401 || res.StatusCode == 403 {
		log.Printf("You are not authorized to submit the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", res.StatusCode)
		return "", fmt.Errorf("Authentication error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response from App Catalog!")
		return "", err
	}

	var responseObject submitResponse
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		if verbose {
			log.Println(err)
		}

		responseObject = submitResponse{}
		responseObject.Messages = []string{string(body)}
	}

	log.Printf("App Catalog returned statuscode %v. Response details:\n", res.StatusCode)
	for _, line := range responseObject.Messages {
		log.Printf("\t%v\n", line)
	}

	if verbose {
		for key, val := range responseObject.Links {
			log.Printf("\tLINK: %s\t\t%s", key, val)
		}
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Submit failed, App Catalog returned statuscode %v", res.StatusCode)
	}

	acceptanceQueryURLPath, _ := responseObject.Links["acc:query"]

	return acceptanceQueryURLPath, nil

}
