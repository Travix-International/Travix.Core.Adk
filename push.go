package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

// PushCommand used for pushing an app during app development
type PushCommand struct {
	appPath     string // path to the App folder
}

const pushTemplateURI = "%s/files/%s"

func configurePushCommand(app *kingpin.Application) {
	cmd := &PushCommand{}
	appCmd := app.Command("push", "Push the App in the specified folder.").
		Action(cmd.push)
	appCmd.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
}

func (cmd *PushCommand) push(context *kingpin.ParseContext) error {
	appPath := cmd.appPath
	rootURI := catalogURIs[targetEnv]

	if appPath == "" {
		appPath = "."
	}

	appPath, err := filepath.Abs(appPath)

	if err != nil {
		log.Printf("Invalid App path: %s\n", appPath)
		return err
	}

	appName := filepath.Base(appPath)
	appManifestFile := appPath + "/app.manifest"
	tempFolder, err := ioutil.TempDir("", "appix")

	if err != nil {
		log.Println("Could not create temp folder!")
		return err
	}

	if _, err = os.Stat(appManifestFile); os.IsNotExist(err) {
		log.Printf("App manifest not found: %s\n", appManifestFile)
		return err
	}

	log.Printf("Run push for App '%s', path '%s'\n", appName, appPath)

	zapFile := tempFolder + "/app.zap"

	if verbose {
		log.Println("Creating ZAP file: " + zapFile)
	}

	err = zipFolder(appPath, zapFile, includePathInZapFile)

	if err != nil {
		log.Println("Could not process App folder!")
		return err
	}

	pushURI := fmt.Sprintf(pushTemplateURI, rootURI, appName)
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	if verbose {
		log.Println("Posting files to App Catalog: " + pushURI)
	}
	request, err := createMultiFileUploadRequest(pushURI, files)
	if err != nil {
		log.Println("Call to App Catalog failed!")
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Call to App Catalog failed!")
		return err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from App Catalog!")
		return err
	}

	var responseLines []string
	err = json.Unmarshal(responseBody, &responseLines)
	if err != nil {
		if verbose {
			log.Println(err)
		}
		responseLines[0] = string(responseBody)
	}

	log.Printf("App Catalog returned statuscode %v. Response details:\n", response.StatusCode)
	for _, line := range responseLines {
		log.Printf("\t%v\n", line)
	}

	if response.StatusCode == http.StatusOK {
		log.Println("App has been pushed successfully.")
	} else {
		return fmt.Errorf("Push failed, App Catalog returned statuscode %v", response.StatusCode)
	}

	return nil
}