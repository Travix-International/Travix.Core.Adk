package appcatalog

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyProgress(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		bodyBytes := []byte(`{ "meta": { "status": "FINISHED", "messages": [ { "widget": "test", "output": "cool" } ] }, "links": { "preview": "http://localhost:3001" } }`)
		_, err := w.Write(bodyBytes)

		if err != nil {
			t.Fatalf("Unable to send the response. Details: %s\n", err.Error())
		}
	}))

	defer testServer.Close()

	quit := make(chan interface{})

	progress := <-verifyProgress(testServer.URL, quit)

	if progress.Meta.Status == pollFinishedStatus {
		t.Log("Application pushed\n")
	} else {
		t.Fatalf("An error occured. The push of the application finished with an unexpected error: %s\n", progress.Meta.Status)
	}
}

func TestVerifyProgressFailedPush(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		bodyBytes := []byte(`{ "meta": { "status": "FAILED", "messages": [ { "widget": "test", "output": "cool" } ] }, "links": { "preview": "http://localhost:3001" } }`)
		_, err := w.Write(bodyBytes)

		if err != nil {
			t.Fatalf("Unable to send the response. Details: %s\n", err.Error())
		}
	}))

	defer testServer.Close()

	quit := make(chan interface{})

	progress := <-verifyProgress(testServer.URL, quit)

	if progress.Meta.Status == pollFailedStatus {
		t.Log("Application push failed\n")
	} else {
		t.Fatalf("An error occured. The push of the application finished with an unexpected error: %s\n", progress.Meta.Status)
	}
}
