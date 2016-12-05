package context

import "github.com/Travix-International/Travix.Core.Adk/lib/auth"

// LoadAuthToken checks if the user is already logged in, it tries to load the locally stored Authentication token, and refreshes it.
// If the user is not logged in, it returns an error.
func (context *Context) LoadAuthToken() (auth.TokenBody, error) {
	firebaseAPIKey := context.Config.FirebaseApiKey
	authData, err := getAuthData(context.Config.AuthFilePath)

	if err != nil {
		return auth.TokenBody{}, err
	}

	authToken, err := refreshTokenOrExit(authData, firebaseAPIKey)

	if err != nil {
		return auth.TokenBody{}, err
	}

	return authToken, nil
}

func getAuthData(authFilePath string) (authResult *auth.AuthData, err error) {
	authResult, err = auth.GetAuthData(authFilePath)
	return authResult, err
}

func refreshTokenOrExit(a *auth.AuthData, firebaseAPIKey string) (auth.TokenBody, error) {
	tokenBody, err := a.RefreshToken(firebaseAPIKey)
	return tokenBody, err
}
