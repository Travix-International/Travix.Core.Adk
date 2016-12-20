package appcatalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Travix-International/appix/auth"
	"github.com/Travix-International/appix/config"
)

type submitResponse struct {
	Messages []string
	Links    map[string]string
}

// SubmitToCatalog submits the specified app to the AppCatalog.
func SubmitToCatalog(submitURI string, appManifestFile string, zapFile string, verbose bool, config config.Config) (acceptanceQueryURL string, err error) {
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	if verbose {
		log.Println("Posting files to App Catalog: " + submitURI)
	}
	request, err := CreateMultiFileUploadRequest(submitURI, files, nil, verbose)
	if err != nil {
		log.Println("Call to App Catalog failed!")
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
		log.Println("Call to App Catalog failed!")
		return "", err
	}

	if response.StatusCode == 401 || response.StatusCode == 403 {
		log.Printf("You are not authorized to submit the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", response.StatusCode)
		return "", fmt.Errorf("Authentication error")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from App Catalog!")
		return "", err
	}

	var responseObject submitResponse
	err = json.Unmarshal(responseBody, &responseObject)
	if err != nil {
		if verbose {
			log.Println(err)
		}

		responseObject = submitResponse{}
		responseObject.Messages = []string{string(responseBody)}
	}

	log.Printf("App Catalog returned statuscode %v. Response details:\n", response.StatusCode)
	for _, line := range responseObject.Messages {
		log.Printf("\t%v\n", line)
	}

	if verbose {
		for key, val := range responseObject.Links {
			log.Printf("\tLINK: %s\t\t%s", key, val)
		}
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Submit failed, App Catalog returned statuscode %v", response.StatusCode)
	}

	acceptanceQueryURLPath, _ := responseObject.Links["acc:query"]

	return acceptanceQueryURLPath, nil
}
