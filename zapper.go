package main

import (
	"strings"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"
)

func prepareAppUpload(configAppPath string) (appPath string, appName string, manifestPath string, err error) {
	if configAppPath == "" {
		configAppPath = "."
	}

	appPath, err = filepath.Abs(configAppPath)

	if err != nil {
		log.Printf("Invalid App path: %s\n", appPath)
		return "", "", "", err
	}

	appName = filepath.Base(appPath)
	manifestPath = appPath + "/app.manifest"
	
	if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
		log.Printf("App manifest not found: %s\n", manifestPath)
		return "", "", "", err
	}
	
	return appPath, appName, manifestPath, nil
}

func createZapPackage(appPath string) (string, error) {
	tempFolder, err := ioutil.TempDir("", "appix")

	if err != nil {
		log.Println("Could not create temp folder!")
		return "", err
	}

	zapFile := tempFolder + "/app.zap"

	if verbose {
		log.Println("Creating ZAP file: " + zapFile)
	}

	err = zipFolder(appPath, zapFile, includePathInZapFile)

	if err != nil {
		log.Println("Could not process App folder!")
		return "", err
	}
	
	return zapFile, err
}

func includePathInZapFile(relPath string, isDir bool) bool {
	path := strings.ToLower(relPath)
	canInclude := strings.HasPrefix(path, "ui/") && // only dirs starting in ui/
		(isDir || strings.Count(path, "/") >= 2) && // only allow subdirs in  ui/
		!strings.Contains(path, "/node_modules/") && // exclude node_modules
		!strings.Contains(path, "/temp/") &&
		!strings.Contains(path, ".git") &&
		!strings.HasSuffix(path, ".idea/") &&
		!strings.HasSuffix(path, ".vscode/") &&
		!strings.HasSuffix(path, ".md") &&
		!strings.HasSuffix(path, ".ds_store") &&
		!strings.HasSuffix(path, "thumbs.db") &&
		!strings.HasSuffix(path, DevFileName) &&
		!strings.HasSuffix(path, "desktop.ini")

	if verbose {
		if canInclude {
			log.Printf("\tAdding %s\n", relPath)
		} else {
			log.Printf("\tSkipping %s\n", relPath)
		}
	}
	return canInclude
}