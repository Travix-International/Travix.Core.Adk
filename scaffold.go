package main

import "log"

func scaffoldNewApp(appPath string) error {
	// Apply the template
	log.Print("Scaffolding the application files...")
	err := applyTemplate(appPath)
	if err != nil {
		log.Printf("Failed to apply template")
		return err
	}

	// Write out the development settings
	log.Print("Scaffolding the development settings...")
	developmentSettings, err := getDefaultDevelopmentSettings()
	if err != nil {
		log.Printf("Failed to generate default development settings")
		return err
	}
	if verbose {
		logDevelopmentSettings(developmentSettings)
	}
	err = writeDevelopmentSettings(appPath, developmentSettings)
	if err != nil {
		log.Printf("Failed to set development settings")
		return err
	}

	return nil
}

func applyTemplate(appPath string) error {
	// TODO: apply the actual templating logic.
	return nil
}
