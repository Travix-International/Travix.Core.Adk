package appcatalog

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadToFrontend(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		bodyBytes := []byte(`{ "links": { "progress":  "https://fireball-dev.travix.com/upload/progress?sessionId=123" } }`)

		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	pollURI, err := UploadToFrontend(testServer.URL, "mocks/mock.js", "test", "132-321", false)

	if err != nil {
		t.Fatalf("TestUpladToFrontend failed. Details: %s\n", err.Error())
	} else {
		if pollURI == "https://fireball-dev.travix.com/upload/progress?sessionId=123" {
			t.Log("The test was successful")
		} else {
			t.Fatalf("Something wrong happened. Output: %s\n", pollURI)
		}
	}
}

func TestUploadToFrontendFail(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		bodyBytes := []byte(`{ "links": { "progress":  "https://fireball-dev.travix.com/upload/progress?sessionId=123" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	pollURI, err := UploadToFrontend(testServer.URL, "mocks/mock.js", "test", "132-321", false)

	if err != nil {
		t.Logf("TestUpladToFrontend failed. Details: %s\n", err.Error())
	} else {
		t.Fatalf("UploadToFrontend must throw an error. Output: %s\n", pollURI)
	}
}

func TestUploadToFrontendBadPath(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		bodyBytes := []byte(`{ "links": { "progress":  "https://fireball-dev.travix.com/upload/progress?sessionId=123" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	pollURI, err := UploadToFrontend(testServer.URL, "mock.js", "test", "132-321", false)

	if err != nil {
		t.Logf("TestUpladToFrontend failed. Details: %s\n", err.Error())
	} else {
		t.Fatalf("UploadToFrontend must throw an error. Output: %s\n", pollURI)
	}
}
