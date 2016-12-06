package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ReadAuthData(authFilePath string) (*AuthData, error) {
	content, readErr := ioutil.ReadFile(authFilePath)
	if readErr != nil {
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

func SaveAuthData(authFilePath string, data *AuthData) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(authFilePath, content, 0644)

	if err != nil {
		return err
	}

	return nil
}
