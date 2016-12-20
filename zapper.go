package appix

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func PrepareAppUpload(configAppPath string) (appPath string, appName string, manifestPath string, err error) {
	if configAppPath == "" {
		configAppPath = "."
	}

	appPath, err = filepath.Abs(configAppPath)

	if err != nil {
		log.Printf("Invalid App path: %s\n", appPath)
		return "", "", "", err
	}

	manifestPath = appPath + "/app.manifest"

	if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
		log.Printf("App manifest not found: %s\n", manifestPath)
		return "", "", "", err
	}

	type AppManifestName struct {
		Name string `json:"name"`
	}

	var manifestObject AppManifestName

	manifestData, err := ioutil.ReadFile(manifestPath)

	if err != nil {
		log.Println("Couldn't read the app.manifest")
		return "", "", "", err
	}

	err = json.Unmarshal(manifestData, &manifestObject)

	if err != nil {
		log.Println("Couldn't parse the app.manifest")
		return "", "", "", err
	}

	if manifestObject.Name == "" {
		log.Println("The name is missing from the app manifest")
		return "", "", "", errors.New("The name is missing from the app manifest")
	}

	appName = manifestObject.Name

	return appPath, appName, manifestPath, nil
}

func CreateZapPackage(appPath string, verbose bool) (string, error) {
	tempFolder, err := ioutil.TempDir("", "appix")

	if err != nil {
		log.Println("Could not create temp folder!")
		return "", err
	}

	zapFile := tempFolder + "/app.zap"

	if verbose {
		log.Println("Creating ZAP file: " + zapFile)
	}

	err = ZipFolder(appPath, zapFile, func(path string) bool {
		ignored, ignoredFolder := IgnoreFilePath(path)
		if verbose && !ignoredFolder {
			if ignored {
				log.Printf("\tSkipping %s\n", path)
			} else {
				log.Printf("\tAdding %s\n", path)
			}
		}
		return !ignored
	})

	if err != nil {
		log.Println("Could not process App folder!")
		return "", err
	}

	return zapFile, err
}
