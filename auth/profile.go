package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// FetchDeveloperProfile gets the profile of the current user from the Profile Api.
func FetchDeveloperProfile(tokenBody TokenBody, developerProfileURL string) (ProfileBody, error) {
	tokenType := tokenBody.TokenType
	tokenValue := tokenBody.IdToken

	baseURL, _ := url.Parse(developerProfileURL)
	relative, _ := url.Parse("profile")

	profileClient := &http.Client{}
	profileReq, profileReqErr := http.NewRequest("GET", baseURL.ResolveReference(relative).String(), nil)
	if profileReqErr != nil {
		return ProfileBody{}, profileReqErr
	}

	profileReq.Header.Set("Content-Type", "application/json")
	profileReq.Header.Set("Authorization", tokenType+" "+tokenValue)
	profileRes, profileResErr := profileClient.Do(profileReq)

	defer profileRes.Body.Close()

	if profileResErr != nil {
		return ProfileBody{}, profileResErr
	}

	profileBody := ProfileBody{}
	err := json.NewDecoder(profileRes.Body).Decode(&profileBody)

	if err != nil {
		log.Println(err)
		return profileBody, err
	}

	return profileBody, nil
}
