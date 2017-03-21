package appcatalog

import "testing"

func TestCreateMultiFileUploadRequest(t *testing.T) {
	files := make(map[string]string)
	files["application"] = "../mocks/mock.js"

	rawFields := make(map[string]string)
	rawFields["field"] = "test"

	url := "https://test.com"

	req, err := CreateMultiFileUploadRequest(url, files, rawFields, false)

	if err != nil {
		t.Fatalf("An error occured while running the happy test for CreateMultiFileUploadRequest. Details: %s\n", err.Error())
	} else {
		if req.URL.String() == url && req.Method == "POST" {
			t.Log("The test succeed in")
		} else {
			t.Fatalf("An error occured while running the happy test for CreateMultiFileUploadRequest. Unexpected URL and method: %s - %s\n", req.URL, req.Method)
		}
	}
}

func TestCreateMultiFileUploadRequestWithoutFiles(t *testing.T) {
	files := make(map[string]string)

	rawFields := make(map[string]string)

	url := "https://test.com"

	req, err := CreateMultiFileUploadRequest(url, files, rawFields, false)

	if err != nil {
		t.Fatalf("An error occured while running the happy test for CreateMultiFileUploadRequest. Details: %s\n", err.Error())
	} else {
		if req.URL.String() == url && req.Method == "POST" {
			t.Log("The test succeed in")
		} else {
			t.Fatalf("An error occured while running the happy test for CreateMultiFileUploadRequest. Unexpected URL and method: %s - %s\n", req.URL, req.Method)
		}
	}
}
