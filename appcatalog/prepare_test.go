package appcatalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

func TestPrepare(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()

	files := make(map[string]string)
	files["mock.js"] = "mocks/mock.js"

	conf := config.Config{
		AuthFilePath:   "mocks/mockAuth.json",
		FirebaseApiKey: "firebaseapikey",
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	_, err := prepare(testServer.URL, files, conf, false, logger)

	if err != nil {
		t.Fatalf("An error occured in prepare. Details: %s\n", err.Error())
	} else {
		t.Log("Everything went well\n")
	}
}

func TestPrepareFail(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer testServer.Close()

	files := make(map[string]string)
	files["mock.js"] = "mock.js"

	conf := config.Config{
		AuthFilePath:   "mocks/mockAuth.json",
		FirebaseApiKey: "firebaseapikey",
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	_, err := prepare(testServer.URL, files, conf, false, logger)

	if err != nil {
		t.Logf("An error occured as expected. Details: %s\n", err.Error())
	} else {
		t.Fatalf("Everything went well\n")
	}
}

func TestCanSendWithoutAuthentication(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer testServer.Close()

	files := make(map[string]string)
	files["mock.js"] = "mocks/mock.js"

	conf := config.Config{
		AuthFilePath:   "",
		FirebaseApiKey: "firebaseapikey",
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	_, err := prepare(testServer.URL, files, conf, false, logger)

	if err == nil {
		t.Logf("No error happened\n")
	} else {
		t.Fatalf("Something went wrong during the test. Details: %s\n", err.Error())
	}
}
