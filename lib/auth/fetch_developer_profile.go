package auth

import (
	"encoding/json"
	"net/http"
)

func FetchDeveloperProfile(tokenBody *TokenBody, developerProfileUrl string) (ProfileBody, error) {
	tokenType := tokenBody.TokenType
	tokenValue := tokenBody.IdToken

	profileClient := &http.Client{}
	profileReq, profileReqErr := http.NewRequest("GET", developerProfileUrl+"/profile", nil)
	if profileReqErr != nil {
		return ProfileBody{}, profileReqErr
	}

	profileReq.Header.Set("Content-Type", "application/json")
	profileReq.Header.Set("Authorization", tokenType+" "+tokenValue)
	profileRes, profileResErr := profileClient.Do(profileReq)
	if profileResErr != nil {
		return ProfileBody{}, profileResErr
	}

	profileBody := ProfileBody{}
	json.NewDecoder(profileRes.Body).Decode(&profileBody)

	return profileBody, nil
}
