package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (auth *AuthData) RefreshToken(firebaseApiKey string) (TokenBody, error) {
	refreshToken := auth.User.StsTokenManager.RefreshToken

	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+firebaseApiKey, bytes.NewBuffer(tokenReqPayload))
	if tokenReqErr != nil {
		return TokenBody{}, tokenReqErr
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenRes, tokenResErr := tokenClient.Do(tokenReq)
	if tokenResErr != nil {
		return TokenBody{}, tokenResErr
	}

	tokenBody := TokenBody{}
	json.NewDecoder(tokenRes.Body).Decode(&tokenBody)

	return tokenBody, nil
}
