package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/nu7hatch/gouuid"
)

/*
   The development settings are stored into a separate file as part of the app. It must not be part of
   the ZAP file when publishing, and whether this file is version controlled is completely up to the user
   as well.
*/

type DevelopmentSettings struct {
	SessionID string // This would be the session ID that the user will use the push/preview his changes
}

// DevFileName is the name of the file which contains the appix development settings for this specific application
const DevFileName = ".appixDevSettings" // Exported for blacklisting use during push/submit

func readDevelopmentSettings(appPath string) (*DevelopmentSettings, error) {
	devSettingsPath := path.Join(appPath, DevFileName)
	if verbose {
		log.Printf("Reading development settings from %s", devSettingsPath)
	}

	data, err := ioutil.ReadFile(devSettingsPath)
	if err != nil {
		log.Printf("Failed to read the development settings file %s", devSettingsPath)
		return nil, err
	}

	settings := DevelopmentSettings{}
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Printf("Failed to unmarshal development settings")
		return nil, err
	}

	return &settings, nil
}

func writeDevelopmentSettings(appPath string, settings *DevelopmentSettings) error {
	devSettingsPath := path.Join(appPath, DevFileName)
	if verbose {
		log.Printf("Writing development settings to %s", devSettingsPath)
	}

	data, err := json.Marshal(*settings)
	if err != nil {
		log.Printf("Failed to marshal development settings to JSON")
		return err
	}

	err = ioutil.WriteFile(devSettingsPath, data, 0664)
	if err != nil {
		log.Printf("Failed to write development settings to %s", devSettingsPath)
		return err
	}

	return nil
}

func getDefaultDevelopmentSettings() (*DevelopmentSettings, error) {
	randomGUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	settings := DevelopmentSettings{SessionID: randomGUID.String()}
	return &settings, nil
}

func logDevelopmentSettings(settings *DevelopmentSettings) {
	log.Printf("Session ID: %s", settings.SessionID)
}
