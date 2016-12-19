package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

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
func LoadAuthToken(config config.Config) (TokenBody, error) {
	firebaseAPIKey := config.FirebaseApiKey
	authData, err := getAuthData(config.AuthFilePath)

	if err != nil {
		return TokenBody{}, err
	}

	if len(authData.Token.IdToken) > 0 && authData.Token.ExpiresAt.Sub(time.Now().UTC()) > time.Duration(5)*time.Minute {
		// If we already have a token, and there is more then 5 minutes until its expiry, we return it.
		return authData.Token, nil
	}

	// Either we don't have a token yet, or it's expired (or close to expiry), so we get a new token.
	authToken, err := refreshToken(authData, firebaseAPIKey)

	if err != nil {
		return TokenBody{}, err
	}

	authData.Token = authToken

	err = SaveAuthData(config.AuthFilePath, authData)

	if err != nil {
		return TokenBody{}, err
	}

	return authToken, nil
}

func getAuthData(authFilePath string) (authResult *AuthData, err error) {
	authResult, err = ReadAuthData(authFilePath)
	return authResult, err
}

func refreshToken(a *AuthData, firebaseAPIKey string) (TokenBody, error) {
	tokenBody, err := a.RefreshToken(firebaseAPIKey)
	return tokenBody, err
}

// RefreshToken refreshes the OAuth access token based on the refresh token.
func (auth *AuthData) RefreshToken(firebaseAPIKey string) (TokenBody, error) {
	refreshToken := auth.User.StsTokenManager.RefreshToken

	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+firebaseAPIKey, bytes.NewBuffer(tokenReqPayload))
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

// ReadAuthData reads the previously persisted access token from the disk.
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

// SaveAuthData persists the access token to disk.
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
