package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

// Auth structs
type StsTokenManager struct {
	ApiKey         string
	RefreshToken   string
	AccessToken    string
	ExpirationTime int
}

type AuthUser struct {
	Uid             string
	DisplayName     string
	Email           string
	EmailVerified   bool
	APIKey          string
	AppName         string
	AuthDomain      string
	StsTokenManager StsTokenManager
}

type AuthCredential struct {
	IdToken     string
	AccessToken string
	Provider    string
}

type AuthData struct {
	User       AuthUser
	Credential AuthCredential
	Token      TokenBody
}

// Token structs
type TokenBody struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    string    `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	IdToken      string    `json:"id_token"`
	ProjectId    string    `json:"project_id"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	UserId       string    `json:"user_id"`
}

// Profile structs
type Profile struct {
	Email          string
	FirebaseUserId string
	Id             int
	IsEnabled      bool
	IsVerified     bool
	Name           string
	PublisherId    string
}

type ProfileBody struct {
	HasProfile bool
	Profile    Profile
}

// LoadAuthToken checks if the user is already logged in, it tries to load the locally stored Authentication token, and refreshes it.
// If the user is not logged in, it returns an error.
func LoadAuthToken(config config.Config, logger *appixLogger.Logger) (TokenBody, error) {
	firebaseAPIKey := config.FirebaseApiKey
	authData, err := readAuthData(config.AuthFilePath)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Could not find authentication data: %s", err.Error()),
			LogEvent: "AppixAuthentication",
		})
		return TokenBody{}, err
	}

	if len(authData.Token.IdToken) > 0 && authData.Token.ExpiresAt.Sub(time.Now().UTC()) > time.Duration(5)*time.Minute {
		// If we already have a token, and there is more then 5 minutes until its expiry, we return it.
		return authData.Token, nil
	}

	// Either we don't have a token yet, or it's expired (or close to expiry), so we get a new token.
	authToken, err := authData.refreshToken(firebaseAPIKey, config.FirebaseRefreshTokenUrl)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Could not retrieve a new token: %s", err.Error()),
			LogEvent: "AppixAuthentication",
		})
		return TokenBody{}, err
	}

	authData.Token = authToken

	err = saveAuthData(config.AuthFilePath, authData)

	if err != nil {
		logger.AddMessageToQueue(appixLogger.LoggerNotification{
			Level:    "error",
			Message:  fmt.Sprintf("Could not save data: %s", err.Error()),
			LogEvent: "AppixAuthentication",
		})
		return TokenBody{}, err
	}

	logger.AddMessageToQueue(appixLogger.LoggerNotification{
		Level:    "info",
		Message:  fmt.Sprintf("User %s successfully connected", authData.User.DisplayName),
		LogEvent: "AppixAuthentication",
	})

	return authToken, nil
}

// refreshToken refreshes the OAuth access token based on the refresh token.
func (auth *AuthData) refreshToken(firebaseAPIKey string, firebaseTokenRefreshURL string) (TokenBody, error) {
	refreshToken := auth.User.StsTokenManager.RefreshToken

	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", firebaseTokenRefreshURL+firebaseAPIKey, bytes.NewBuffer(tokenReqPayload))
	if tokenReqErr != nil {
		return TokenBody{}, tokenReqErr
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenRes, tokenResErr := tokenClient.Do(tokenReq)

	if tokenResErr != nil {
		return TokenBody{}, tokenResErr
	}

	if tokenRes.StatusCode < 200 || tokenRes.StatusCode > 399 {
		return TokenBody{}, fmt.Errorf("An error occured while requesting the Firebase token. Status code: %d\n", tokenRes.StatusCode)
	}

	tokenBody := TokenBody{}
	json.NewDecoder(tokenRes.Body).Decode(&tokenBody)

	// The api only returns the ExpiresIn field, but in later stages we need to know the actual Time when it is expiring.
	expiresIn, _ := strconv.Atoi(tokenBody.ExpiresIn)
	tokenBody.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)

	return tokenBody, nil
}

// readAuthData reads the previously persisted access token from the disk.
func readAuthData(authFilePath string) (*AuthData, error) {
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

// saveAuthData persists the access token to disk.
func saveAuthData(authFilePath string, data *AuthData) error {
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
