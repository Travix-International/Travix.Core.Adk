package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func GetAuthData(authFilePath string) (*AuthData, error) {
	content, readErr := ioutil.ReadFile(authFilePath)
	if readErr != nil {
		log.Printf("Failed to read auth file %s", authFilePath)
		return nil, readErr
	}

	auth := AuthData{}
	unmarshalErr := json.Unmarshal(content, &auth)

	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal auth content")
		return nil, unmarshalErr
	}

	return &auth, nil
}
