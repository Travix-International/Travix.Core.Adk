package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// PublishCommand used for publishing Apps
type PublishCommand struct {
	appPath     string // path to the App folder
	enviromnent string // environment (default: dev)
}

const publishTemplateURI = "%s/files/%s"

func configurePublishCommand(app *kingpin.Application) {
	cmd := &PublishCommand{}
	appCmd := app.Command("publish", "Publish the App in the specified folder to the specified enviromnent.").
		Action(cmd.publish).
		Alias("p").
		Alias("pub")
	appCmd.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
	appCmd.Arg("environment", "Target environment (dev/acc/prod, default dev)").
		Default("dev").
		EnumVar(&cmd.enviromnent, "dev", "acc", "prod")
}

func (cmd *PublishCommand) publish(context *kingpin.ParseContext) error {
	appPath := cmd.appPath
	environment := cmd.enviromnent
	rootURI := catalogURIs[targetEnv]

	if appPath == "" {
		appPath = "."
	}
	if environment == "" {
		environment = "dev"
	}

	appPath, err := filepath.Abs(appPath)

	if err != nil {
		fmt.Printf("Invalid App path: %s\n", appPath)
		return err
	}

	appName := filepath.Base(appPath)
	appManifestFile := appPath + "/app.manifest"
	tempFolder, err := ioutil.TempDir("", "appix")

	if err != nil {
		fmt.Println("Could not create temp folder!")
		return err
	}

	if _, err = os.Stat(appManifestFile); os.IsNotExist(err) {
		fmt.Printf("App manifest not found: %s\n", appManifestFile)
		return err
	}

	fmt.Printf("Run publish for app %s, env %s, path %s\n", appName, environment, appPath)

	zapFile := tempFolder + "/app.zap"

	fmt.Println("Creating ZAP file: " + zapFile)
	err = zipFolder(appPath, zapFile)

	if err != nil {
		fmt.Println("Could not process App folder!")
		return err
	}

	publishURI := fmt.Sprintf(publishTemplateURI, rootURI, appName)
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	fmt.Println("Posting files to App Catalog: " + publishURI)
	request, err := createMultiFileUploadRequest(publishURI, files)
	if err != nil {
		fmt.Println("Call to App Catalog failed!")
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Call to App Catalog failed!")
		return err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response from App Catalog!")
		return err
	}

	var responseLines []string
	err = json.Unmarshal(responseBody, &responseLines)
	if err != nil {
		fmt.Println(err)
		responseLines[0] = string(responseBody)
	}

	fmt.Printf("App Catalog returned statuscode %v. Response details:\n", response.StatusCode)
	for _, line := range responseLines {
		fmt.Printf("\t%v\n", line)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("App Catalog returns statuscode %v", response.StatusCode)
	}

	return nil
}
