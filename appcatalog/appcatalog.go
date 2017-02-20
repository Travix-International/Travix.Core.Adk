package appcatalog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type appCatalogResponse struct {
	Messages []string
	Links    map[string]string
}

func extractAppCatalogResponse(body io.ReadCloser) (*appCatalogResponse, error) {
	b, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var responseObject *appCatalogResponse
	err = json.Unmarshal(b, &responseObject)

	if err != nil {
		return &appCatalogResponse{
			Messages: []string{string(b)},
		}, err
	}
	return responseObject, nil
}

func doRequest(reqType string, link string, req *http.Request, maxTimeoutValue time.Duration, verbose bool) (string, error) {
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

	responseObject, err := extractAppCatalogResponse(res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", &catalogError{
			operation:  reqType,
			statusCode: res.StatusCode,
			response:   responseObject,
		}
	}

	if err != nil {
		log.Println("Error reading response from App Catalog!")
		if verbose {
			log.Print(err)
		}
		return "", err
	}

	log.Print(responseObject.getDetails(res.StatusCode))
	return responseObject.Links[link], nil
}

func (acr *appCatalogResponse) getDetails(statusCode int) string {
	msg := fmt.Sprintf("App Catalog returned statuscode %v. Response details:\n", statusCode)

	for _, line := range acr.Messages {
		msg += fmt.Sprintf("\t%v\n", line)
	}
	return msg
}

func (acr *appCatalogResponse) displayLink() string {
	msg := ""

	for key, val := range acr.Links {
		msg += fmt.Sprintf("\tLINK %s\t\t%s", key, val)
	}
	return msg
}
