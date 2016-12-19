package appcatalog

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Creates a new file upload http request with optional extra params
// Source: https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func CreateMultiFileUploadRequest(uri string, files map[string]string, rawFields map[string]string, verbose bool) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if rawFields != nil {
		// We add the raw form parameters to the request.
		for key, value := range rawFields {
			field, err := writer.CreateFormField(key)

			if err != nil {
				return nil, err
			}

			if _, err := field.Write([]byte(value)); err != nil {
				return nil, err
			}
		}
	}

	if files != nil {
		// We add the posted files to the request.
		for key, path := range files {
			if verbose {
				log.Println("Trying to add to multi-upload: " + path)
			}
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(key, path)

			if err != nil {
				return nil, err
			}

			_, err = io.Copy(part, file)

			if err != nil {
				return nil, err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)

	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}
