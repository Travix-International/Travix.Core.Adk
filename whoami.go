package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureWhoamiCommand(app *kingpin.Application) {
	app.Command("whoami", "Displays logged in user's email address").
		Action(executeWhoamiCommand)
}

func executeWhoamiCommand(context *kingpin.ParseContext) error {
	config := GetConfig()

	auth, authErr := GetAuth(config)
	if authErr != nil {
		log.Fatal(authErr)
		return nil
	}

	accessToken := auth.Credential.AccessToken

	type UserInfo struct {
		Email string
		FamilyName string
		Gender string
		GivenName string
		Hd string
		Id string
		Link string
		Locale string
		Name string
		Picture string
		VerifiedEmail string
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	req.Header.Set("Authorization", "Bearer " + accessToken)
	res, _ := client.Do(req)

	userInfo := UserInfo{}
	json.NewDecoder(res.Body).Decode(&userInfo)
	fmt.Println(userInfo.Email)

	return nil
}
