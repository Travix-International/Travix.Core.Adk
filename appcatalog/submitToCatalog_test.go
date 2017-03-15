package appcatalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

/**
 * - happy flow
 * - retries (1 - 3)
 * - catalog error (authentication)
 * - break prepare
 */

func TestSubmitToCatalog(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(`{ "messages": [ "success" ], "links": { "acc:query": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "./mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:     "mocks/mockAuth.json",
		FirebaseApiKey:   "firebaseapikey",
		MaxRetryAttempts: 3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := SubmitToCatalog(testServer.URL, 10, appManifestFile, "mocks/mock.js", false, conf, logger)

	if err != nil {
		t.Fatalf("An error occured when testing SubmitToCatalog. Details: %s\n", err.Error())
	} else {
		if uploadURI == "http://localhost:3001" {
			t.Logf("The test for SubmitToCatalog went well.")
		} else {
			t.Fatalf("Unexpected returned value: %s\n", uploadURI)
		}
	}
}

func TestSubmitToCatalogWithRetries(t *testing.T) {
	counter := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte

		if counter == 2 {
			w.WriteHeader(http.StatusOK)
			bodyBytes = []byte(`{ "messages": [ "success" ], "links": { "acc:query": "http://localhost:3001" } }`)
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			bodyBytes = []byte(`{ "messages": [ "failed" ], "links": { "acc:query": "http://localhost:3001" } }`)
		}

		w.Write(bodyBytes)
		counter++
	}))

	defer testServer.Close()

	appManifestFile := "mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:     "mocks/mockAuth.json",
		FirebaseApiKey:   "firebaseapikey",
		MaxRetryAttempts: 3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	_, err := SubmitToCatalog(testServer.URL, 10, appManifestFile, "mocks/mock.js", false, conf, logger)

	if err != nil {
		t.Fatalf("An error occured when testing SubmitToCatalog. Details: %s\n", err.Error())
	} else {
		t.Logf("The test for SubmitToCatalog went well.")
	}
}

func TestSubmitToCatalogFailAuthentication(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		bodyBytes := []byte(`{ "messages": [ "failed" ], "links": { "acc:query": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "./mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:     "mocks/mockAuth.json",
		FirebaseApiKey:   "firebaseapikey",
		MaxRetryAttempts: 3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := SubmitToCatalog(testServer.URL, 10, appManifestFile, "mocks/mock.js", false, conf, logger)

	if err != nil {
		t.Logf("An error occured when testing SubmitToCatalog as expected. Details: %s\n", err.Error())
	} else {
		t.Fatalf("Expecting an error. Unexpected returned value: %s\n", uploadURI)
	}
}

func TestSubmitToCatalogFailServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		bodyBytes := []byte(`{ "messages": [ "failed" ], "links": { "acc:query": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "./mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:     "mocks/mockAuth.json",
		FirebaseApiKey:   "firebaseapikey",
		MaxRetryAttempts: 3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := SubmitToCatalog(testServer.URL, 10, appManifestFile, "mocks/mock.js", false, conf, logger)

	if err != nil {
		t.Logf("An error occured when testing SubmitToCatalog as expected. Details: %s\n", err.Error())
	} else {
		t.Fatalf("Expecting an error. Unexpected returned value: %s\n", uploadURI)
	}
}
