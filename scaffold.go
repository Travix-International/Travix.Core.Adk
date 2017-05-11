package appix

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func scaffoldNewApp(appPath, templateName string, verbose bool) error {
	// Apply the template
	log.Print("Scaffolding the application files...")
	err := applyTemplate(appPath, templateName)
	if err != nil {
		log.Printf("Failed to apply template")
		return err
	}

	// Write out the development settings
	log.Print("Initializing the development settings...")
	developmentSettings, err := getDefaultDevelopmentSettings()
	if err != nil {
		log.Printf("Failed to generate default development settings")
		return err
	}
	if verbose {
		logDevelopmentSettings(developmentSettings)
	}
	err = writeDevelopmentSettings(appPath, developmentSettings, verbose)
	if err != nil {
		log.Printf("Failed to store the development settings")
		return err
	}

	return nil
}

func applyTemplate(appPath, templateName string) error {
	helloWorldTemplateURL := "https://github.com/Travix-International/travix-appix-templates/archive/master.zip"
	zipSubDirectory := "travix-appix-templates-master"
	zipSubDirectory = fmt.Sprintf("%s/%s/", zipSubDirectory, templateName)

	tempFolder, err := ioutil.TempDir("", "appix")
	if err != nil {
		return err
	}

	tempFile := filepath.Join(tempFolder, "template.zip")
	out, err := os.Create(tempFile)
	if err != nil {
		return err
	}

	defer out.Close()

	log.Printf("Downloading template '%s' from %s\n", templateName, helloWorldTemplateURL)

	// Download the template.
	resp, err := http.Get(helloWorldTemplateURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Unzip the template and do some basic replacements.
	var fileCount int
	fileCount, err = extractZip(tempFile, appPath, zipSubDirectory)
	if err != nil {
		return err
	}

	if fileCount <= 0 {
		errorText := fmt.Sprintf("Template '%s' does not exist", templateName)
		log.Print(errorText)
		return errors.New(errorText)
	}

	appName := filepath.Base(appPath)

	// This is just a temporary, proof of concept implementation, we'll need a proper solution for scaffolding and templating.
	if err = ReplaceInFile(filepath.Join(appPath, "app.manifest"), "{APP_NAME}", appName); err != nil {
		return err
	}

	return nil
}

func ReplaceInFile(file, from, to string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	output := bytes.Replace(input, []byte(from), []byte(to), -1)

	err = ioutil.WriteFile(file, output, 0666)

	return err
}
