package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
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

	// The api only returns the ExpiresIn field, but in later stages we need to know the actual Time when it is expiring.
	expiresIn, _ := strconv.Atoi(tokenBody.ExpiresIn)
	tokenBody.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)

	return tokenBody, nil
}
