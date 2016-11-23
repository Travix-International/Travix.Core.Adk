package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureWhoamiCommand(app *kingpin.Application) {
	app.Command("whoami", "Displays logged in user's information").
		Action(executeWhoamiCommand)
}

func executeWhoamiCommand(context *kingpin.ParseContext) error {
	config := GetConfig()

	auth, authErr := GetAuth(config)
	if authErr != nil {
		log.Fatal(authErr)
		return nil
	}

	refreshToken := auth.User.StsTokenManager.RefreshToken

	// fetch refreshed token
	type TokenBody struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		IdToken      string `json:"id_token"`
		ProjectId    string `json:"project_id"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		UserId       string `json:"user_id"`
	}

	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+config.FirebaseApiKey, bytes.NewBuffer(tokenReqPayload))
	if tokenReqErr != nil {
		log.Fatal(tokenReqErr)
		return nil
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenRes, tokenResErr := tokenClient.Do(tokenReq)
	if tokenResErr != nil {
		log.Fatal(tokenResErr)
		return nil
	}

	tokenBody := TokenBody{}
	json.NewDecoder(tokenRes.Body).Decode(&tokenBody)

	tokenType := tokenBody.TokenType
	tokenValue := tokenBody.IdToken

	// fetch profile
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

	profileClient := &http.Client{}
	profileReq, profileReqErr := http.NewRequest("GET", config.DeveloperProfileUrl+"/profile", nil)
	if profileReqErr != nil {
		log.Fatal(profileReqErr)
		return nil
	}

	profileReq.Header.Set("Content-Type", "application/json")
	profileReq.Header.Set("Authorization", tokenType+" "+tokenValue)
	profileRes, profileResErr := profileClient.Do(profileReq)
	if profileResErr != nil {
		log.Fatal(profileResErr)
		return nil
	}

	profileBody := ProfileBody{}
	json.NewDecoder(profileRes.Body).Decode(&profileBody)

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
}
