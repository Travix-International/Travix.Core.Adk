package context

import (
	"time"

	"github.com/Travix-International/Travix.Core.Adk/lib/auth"
)

// LoadAuthToken checks if the user is already logged in, it tries to load the locally stored Authentication token, and refreshes it.
// If the user is not logged in, it returns an error.
func (context *Context) LoadAuthToken() (auth.TokenBody, error) {
	firebaseAPIKey := context.Config.FirebaseApiKey
	authData, err := getAuthData(context.Config.AuthFilePath)

	if err != nil {
		return auth.TokenBody{}, err
	}

	if len(authData.Token.IdToken) > 0 && authData.Token.ExpiresAt.Sub(time.Now().UTC()) > time.Duration(5)*time.Minute {
		// If we already have a token, and there is more then 5 minutes until its expiry, we return it.
		return authData.Token, nil
	}

	// Either we don't have a token yet, or it's expired (or close to expiry), so we get a new token.
	authToken, err := refreshToken(authData, firebaseAPIKey)

	if err != nil {
		return auth.TokenBody{}, err
	}

	authData.Token = authToken

	err = auth.SaveAuthData(context.Config.AuthFilePath, authData)

	if err != nil {
		return auth.TokenBody{}, err
	}

	return authToken, nil
}

func getAuthData(authFilePath string) (authResult *auth.AuthData, err error) {
	authResult, err = auth.ReadAuthData(authFilePath)
	return authResult, err
}

func refreshToken(a *auth.AuthData, firebaseAPIKey string) (auth.TokenBody, error) {
	tokenBody, err := a.RefreshToken(firebaseAPIKey)
	return tokenBody, err
}
