package appcatalog

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDoRequestHappyFlow(t *testing.T) {
	testAC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer testAC.Close()

	req, err := http.NewRequest("POST", testAC.URL, &bytes.Buffer{})

	if err != nil {
		t.Fatalf("An error occured when creating the request for testing the DoRequest\n")
	}

	if uploadURI, err := doRequest("Push", "upload", req, time.Duration(2)*time.Second, false); err == nil {
		t.Fatalf("An error occured while performing the doRequest Test. Details: %s\n", err.Error())
	} else {
		t.Logf("successful: %s\n", uploadURI)
	}
}

func TestDoRequestServerError(t *testing.T) {
	testAC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer testAC.Close()

	req, err := http.NewRequest("POST", testAC.URL, &bytes.Buffer{})

	if err != nil {
		t.Fatalf("An error occured when creating the request for testing the DoRequest\n")
	}

	if _, err := doRequest("Push", "upload", req, time.Duration(2)*time.Second, false); err != nil {
		t.Logf("doRequest return an error. Details: %s\n", err.Error())
	} else {
		t.Fatal("doRequest must return an error\n")
	}
}
