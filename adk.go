package main

import (
	"fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    appName := "testApp"
    appRoot := "/Projects/go/src/bitbucket.org/xivart/travix.core.adk/content/TestApp/"
    zapFile := "/Projects/go/src/bitbucket.org/xivart/travix.core.adk/temp/app.zap"
    
    err := zipFolder(appRoot, zapFile)
    if (err != nil) {
        fmt.Println("Could not process App folder!")
        panic(err)
    }    
    
    appCatalogURI := fmt.Sprintf("http://localhost:52426/files/%s", appName)
    files := map[string]string {
        "manifest": appRoot + "app.manifest",
        "zapfile": zapFile,
    }

    fmt.Println("Posting files to App Catalog: " + appCatalogURI)
    request, err := createMultiFileUploadRequest(appCatalogURI, files)
    if (err != nil) {
        fmt.Println("Call to appCatalog failed!")
        panic(err)
    }
    
    client := &http.Client{}
    response, err := client.Do(request)
    if (err != nil) {
        fmt.Println("Call to appCatalog failed!")
        panic(err)
    }
     
    fmt.Println("Response from App Catalog:")
    fmt.Println(response.StatusCode)
    fmt.Println(response.Header)
    responseBody, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(responseBody))
    
    if response.StatusCode != http.StatusOK {
        panic(nil) //TODO
    }
}