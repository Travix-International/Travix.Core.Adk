package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"
)

// PublishCommand used for publishing Apps
type PublishCommand struct {
	appPath     string // path to the App folder
	environment string // environment (default: dev)
}

const publishTemplateURI = "%s/files/%s"

func configurePublishCommand(app *kingpin.Application) {
	cmd := &PublishCommand{}
	appCmd := app.Command("publish", "Publish the App in the specified folder to the specified enviromnent.").
		Action(cmd.publish).
		Alias("pub").
		Alias("publ")
	appCmd.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
	appCmd.Arg("environment", "Target environment (dev/acc/prod, default dev)").
		Default("dev").
		EnumVar(&cmd.environment, "dev", "acc", "prod")
}

func (cmd *PublishCommand) publish(context *kingpin.ParseContext) error {
	appPath := cmd.appPath
	environment := cmd.environment
	
	appPath, appName, environment, appManifestFile, err := prepareAppUpload(cmd.appPath, cmd.environment)
	
	if err != nil {
		log.Println("Could not prepare the app folder for uploading")
		return err
	}
	
	zapFile, err:= createZapPackage(appPath)

	if err != nil {
		log.Println("Could not create zap package!")
		return err
	}
	
	log.Printf("Run publish for App '%s', env '%s', path '%s'\n", appName, environment, appPath)

	rootURI := catalogURIs[targetEnv]
	publishURI := fmt.Sprintf(publishTemplateURI, rootURI, appName)
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	if verbose {
		log.Println("Posting files to App Catalog: " + publishURI)
	}
	request, err := createMultiFileUploadRequest(publishURI, files)
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
		log.Println("App has been published successfully.")
	} else {
		return fmt.Errorf("Publish failed, App Catalog returned statuscode %v", response.StatusCode)
	}

	return nil
}