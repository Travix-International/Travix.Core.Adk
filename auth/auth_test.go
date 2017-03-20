package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

func TestSaveAuthData(t *testing.T) {
	authPath := "../mocks/auth.json"
	authData := &AuthData{}

	err := saveAuthData(authPath, authData)

	if err != nil {
		t.Fatalf("An error occured when saving the authentication data. Details: %s\n", err.Error())
	} else {
		if _, err := os.Stat(authPath); os.IsNotExist(err) {
			t.Fatal("The file hasn't been created\n")
		}
		t.Log("The test is passing\n")
	}
}

func TestSaveAuthDataFail(t *testing.T) {
	authPath := "test/auth.json"
	authData := &AuthData{}

	err := saveAuthData(authPath, authData)

	if err == nil {
		t.Fatal("No error occured when saving the authentication data.\n")
	} else {
		t.Log("The test is passing\n")
	}
}

func TestReadAuthData(t *testing.T) {
	authData, err := readAuthData("../mocks/mockAuth.json")

	if err != nil {
		t.Fatalf("An error occured while reading the authentication data. Details: %s\n", err.Error())
	} else {
		if authData.User.DisplayName == "Mock mock" {
			t.Log("Everything went well\n")
		} else {
			t.Fatal("An error occured\n")
		}
	}
}

func TestReadAuthDataFail(t *testing.T) {
	_, err := readAuthData("auth.json")

	if err == nil {
		t.Fatal("No error occured when saving the authentication data.\n")
	} else {
		t.Log("The test is passing\n")
	}
}

func TestRefreshToken(t *testing.T) {
	// "https://securetoken.googleapis.com/v1/token?key="
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(`{ "access_token": "therightone", "expires_in": "3600", "expires_at": "1489675496997", "id_token": "idtoken", "project_id": "projectid", "refresh_token": "refreshtoken", "token_type": "Bearer", "user_id": "toto" }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	authData := &AuthData{
		User: AuthUser{
			StsTokenManager: StsTokenManager{
				RefreshToken: "xxxxxxxxxxxxxxxxxxxxxxx",
			},
		},
	}

	tokenBody, err := authData.refreshToken("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", testServer.URL+"?key=")

	if err != nil {
		t.Fatalf("An error occured when testing refreshToken. Details: %s\n", err.Error())
	} else {
		if tokenBody.AccessToken == "therightone" {
			t.Log("The test succeed in\n")
		} else {
			t.Fatalf("Unexpected values in the response. Details: %s\n", tokenBody.AccessToken)
		}
	}
}

func TestRefreshTokenFail(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := make([]byte, 1)
		w.Write(response)
	}))

	defer testServer.Close()

	authData := &AuthData{
		User: AuthUser{
			StsTokenManager: StsTokenManager{
				RefreshToken: "xxxxxxxxxxxxxxxxxxxxxxx",
			},
		},
	}

	_, err := authData.refreshToken("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", testServer.URL+"?key=")

	if err == nil {
		t.Fatal("It has to throw an error\n")
	} else {
		t.Logf("Everything went well. Error details: %s\n", err.Error())
	}
}

func TestLoadAuthToken(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(`{ "access_token": "therightone", "expires_in": "3600", "expires_at": "1489675496997", "id_token": "idtoken", "project_id": "projectid", "refresh_token": "refreshtoken", "token_type": "Bearer", "user_id": "toto" }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	conf := config.Config{
		AuthFilePath:            "../mocks/mockAuth.json",
		FirebaseApiKey:          "firebaseapikey",
		FirebaseRefreshTokenUrl: testServer.URL + "/firebase?key=",
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	tokenBody, err := LoadAuthToken(conf, logger)

	if err != nil {
		t.Fatalf("An error happened while test LoadAuthToken. Details: %s\n", err.Error())
	} else {
		if tokenBody.AccessToken == "therightone" {
			t.Log("Everything went well\n")
		} else {
			t.Fatalf("An error occured while performing the test. Details: %s\n", tokenBody.AccessToken)
		}
	}
}
