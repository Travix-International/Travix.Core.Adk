package context

import (
	"log"
	"os"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
)

// RequireUserLoggedIn makes sure that user is logged in and attempts to refresh the auth token.
// If it succeeds, then it places the auth token in the context for client use.
func (context *Context) RequireUserLoggedIn(command string) {
	firebaseApiKey := context.Config.FirebaseApiKey
	authData := getAuthOrExit(context.Config.AuthFilePath, command)
	authToken := refreshTokenOrExit(authData, firebaseApiKey)
	context.AuthToken = &authToken
}

func getAuthOrExit(authFilePath string, command string) *auth.AuthData {
	auth, err := auth.GetAuthData(authFilePath)
	if err != nil {
		log.Printf("You must be logged in in order to run '%s' command\n", command)
		log.Println("Use: appix login")
		os.Exit(1)
	}
	return auth
}

func refreshTokenOrExit(a *auth.AuthData, firebaseApiKey string) auth.TokenBody {
	tokenBody, err := a.RefreshToken(firebaseApiKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return tokenBody
}
