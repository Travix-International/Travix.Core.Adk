package appix

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ScaffoldNewApp(appPath string, verbose bool) error {
	// Apply the template
	log.Print("Scaffolding the application files...")
	err := ApplyTemplate(appPath)
	if err != nil {
		log.Printf("Failed to apply template")
		return err
	}

	// Write out the development settings
	log.Print("Initializing the development settings...")
	developmentSettings, err := GetDefaultDevelopmentSettings()
	if err != nil {
		log.Printf("Failed to generate default development settings")
		return err
	}
	if verbose {
		LogDevelopmentSettings(developmentSettings)
	}
	err = WriteDevelopmentSettings(appPath, developmentSettings, verbose)
	if err != nil {
		log.Printf("Failed to store the development settings")
		return err
	}

	return nil
}

func ApplyTemplate(appPath string) error {
	helloWorldTemplateURL := "https://raw.githubusercontent.com/Travix-International/travix-fireball-app-templates/master/HelloWorldTemplate.zip"

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

	log.Println("Downloading template from " + helloWorldTemplateURL)

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
	err = ExtractZip(tempFile, appPath)
	if err != nil {
		return err
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
