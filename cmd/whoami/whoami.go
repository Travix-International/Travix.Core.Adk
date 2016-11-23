package whoami

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	libAuth "github.com/Travix-International/Travix.Core.Adk/lib/auth"
	modelsConfig "github.com/Travix-International/Travix.Core.Adk/models/config"
	modelsContext "github.com/Travix-International/Travix.Core.Adk/models/context"
)

type TokenBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	IdToken      string `json:"id_token"`
	ProjectId    string `json:"project_id"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserId       string `json:"user_id"`
}

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

func fetchRefreshedToken(config modelsConfig.Config, refreshToken string) (TokenBody, error) {
	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+config.FirebaseApiKey, bytes.NewBuffer(tokenReqPayload))
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

func fetchDeveloperProfile(config modelsConfig.Config, tokenBody TokenBody) (ProfileBody, error) {
	tokenType := tokenBody.TokenType
	tokenValue := tokenBody.IdToken

	profileClient := &http.Client{}
	profileReq, profileReqErr := http.NewRequest("GET", config.DeveloperProfileUrl+"/profile", nil)
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

func Register(context modelsContext.Context) {
	config := context.Config

	context.App.Command("whoami", "Displays logged in user's information").
		Action(func(parseContext *kingpin.ParseContext) error {
			// get locally stored auth info
			auth, authErr := libAuth.GetAuth(config)
			if authErr != nil {
				log.Fatal(authErr)
				return nil
			}

			// fetch refreshed token
			refreshToken := auth.User.StsTokenManager.RefreshToken
			tokenBody, tokenBodyErr := fetchRefreshedToken(config, refreshToken)
			if tokenBodyErr != nil {
				log.Fatal(tokenBodyErr)
				return nil
			}

			// fetch profile
			profileBody, profileBodyErr := fetchDeveloperProfile(config, tokenBody)
			if profileBodyErr != nil {
				log.Fatal(profileBodyErr)
				return nil
			}

			if profileBody.HasProfile {
				fmt.Println("Email: " + profileBody.Profile.Email)
				fmt.Println("Name: " + profileBody.Profile.Name)

				if profileBody.Profile.IsEnabled == true {
					fmt.Println("Enabled: Yes")
				} else {
					fmt.Println("Enabled: No")
				}

				if profileBody.Profile.IsVerified == true {
					fmt.Println("Verified: Yes")
				} else {
					fmt.Println("Verified: No")
				}

				fmt.Println("Publisher ID: " + profileBody.Profile.PublisherId)
			} else {
				fmt.Println("No profile found.")
			}

			return nil
		})
}
