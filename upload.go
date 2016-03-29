package main

import (
	"bytes"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

// Creates a new file upload http request with optional extra params
// Source: https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func createMultiFileUploadRequest(uri string, files map[string]string) (*http.Request, error) {
    
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    for key, path := range files {
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
      
    err := writer.Close()
    if err != nil {
        return nil, err
    }
    
    request, err := http.NewRequest("POST", uri, body)
    
    if err != nil {
        return nil, err
    }
    request.Header.Set("Content-Type", writer.FormDataContentType())
    
    return request, nil;
}
