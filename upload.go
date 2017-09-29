package appix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const headerAuthKey = "FIREBASE-TOKEN"

// SignedUploadURL : an object describing the answer from the service generating the upload URLs
type SignedUploadURL struct {
	SignedUploadURL string `json:"signedUploadUrl"`
	ZipFileName     string `json:"zipFileName"`
	TokenID         string `json:"tokenId"`
}

// RetrieveUploadURL : request a signed URL for pushing an app.
func RetrieveUploadURL(uploadURL string, authToken string, appName string, sessionID string) (*SignedUploadURL, error) {
	reqBody := []byte("{ \"appName\": \"" + appName + "\", \"sessionId\": \"" + sessionID + "\" }")
	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set(headerAuthKey, authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var respBody = &SignedUploadURL{}

	err = json.Unmarshal(body, respBody)

	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// UploadResource : upload a resource with the given upload URL
func (suu *SignedUploadURL) UploadResource(appZipPath string) error {
	reqBody, err := os.Open(appZipPath)

	if err != nil {
		return err
	}
	defer reqBody.Close()

	req, err := http.NewRequest("PUT", suu.SignedUploadURL, reqBody)

	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("An error occured while uploading the zap file. Status code: %d", resp.StatusCode)
}
