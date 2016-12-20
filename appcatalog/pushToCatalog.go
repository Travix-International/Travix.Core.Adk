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

// PushToCatalog pushes the specified app to the AppCatalog.
func PushToCatalog(pushURI string, appManifestFile string, verbose bool, config config.Config) (uploadURI string, err error) {
	// To the App Catalog we have to POST the manifest in a multipart HTTP form.
	// When doing the push, it'll only contain a single file, the manifest.
	files := map[string]string{
		"manifest": appManifestFile,
	}

	if verbose {
		log.Println("Posting the app manifest to the App Catalog overlay: " + pushURI)
	}

	request, err := CreateMultiFileUploadRequest(pushURI, files, nil, verbose)

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

	if verbose {
		logServerResponse(response)
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
