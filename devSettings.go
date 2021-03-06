package appix

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/Travix-International/appix/config"
	"github.com/nu7hatch/gouuid"
)

/*
   The development settings are stored into a separate file as part of the app. It must not be part of
   the ZAP file when publishing, and whether this file is version controlled is completely up to the user
   as well.
*/

type DevelopmentSettings struct {
	SessionID         string // This would be the session ID that the user will use the push/preview his changes.
	DevServerOverride string // When pushing an app, the default dev server returned by the app catalog can be overwritten with this.
}

func readDevelopmentSettings(appPath string, verbose bool) (*DevelopmentSettings, error) {
	devSettingsPath := path.Join(appPath, config.DevFileName)
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

func writeDevelopmentSettings(appPath string, settings *DevelopmentSettings, verbose bool) error {
	devSettingsPath := path.Join(appPath, config.DevFileName)
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
