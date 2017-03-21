package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchDeveloperProfile(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		response := []byte(`{ "hasProfile": true,  "profile": { "email": "test@test.com", "fireballUserId": 1, "id": 123, "isEnabled": true, "isVerified": true, "name": "Mock mock", "publisherId": "123" }}`)

		w.Write(response)
	}))

	defer testServer.Close()

	tokenBody := TokenBody{
		TokenType: "Bearer",
		IdToken:   "xxxxxxxxxxxxxxxxxxx",
	}

	profile, err := FetchDeveloperProfile(tokenBody, testServer.URL)

	if err != nil {
		t.Fatalf("An error occured while testing FetchDeveloperProfile. Details: %s\n", err.Error())
	} else {
		if profile.HasProfile == true && profile.Profile.Email == "test@test.com" {
			t.Log("The test went well\n")
		} else {
			t.Fatalf("Something went wrong when performing the test. Details: %s\n", err.Error())
		}
	}
}

func TestFetchDeveloperProfileFailBadBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer testServer.Close()

	tokenBody := TokenBody{
		TokenType: "Bearer",
		IdToken:   "xxxxxxxxxxxxxxxxxxx",
	}

	_, err := FetchDeveloperProfile(tokenBody, testServer.URL)

	if err == nil {
		t.Fatal("No error occured during the test\n")
	} else {
		t.Log("The test failed as expected\n")
	}
}
