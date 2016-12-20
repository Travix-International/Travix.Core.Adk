package appcatalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// UploadToFrontend uploads the app package to the frontend for bundling.
func UploadToFrontend(uploadURI string, zapFile string, appName string, sessionID string, verbose bool) (pollURI string, err error) {
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

	request, err := CreateMultiFileUploadRequest(uploadURI, files, params, verbose)

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
