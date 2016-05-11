package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

type InitCommand struct {
	appPath string // Path to where the app will be placed after the scaffold has taken place
}

// Configures the command line on how to deal with the init command
func configureInitCommand(app *kingpin.Application) {
	cmd := &InitCommand{}
	appCmd := app.Command("init", "Scaffold a new application into the specified folder").
		Action(cmd.init).
		Alias("i")
	appCmd.Arg("appPath", "Path to an emty folder. (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
}

func (cmd *InitCommand) init(context *kingpin.ParseContext) error {

	// grab the absolute path
	appPathRelative := cmd.appPath
	appPathAbsolute, err := filepath.Abs(appPathRelative)
	if err != nil {
		log.Printf("Failed to obtain absolute path for %s\n%v", appPathRelative, err)
		return err
	}

	// tell the user what we're planning on doing
	log.Print("Initializing new application")
	if verbose {
		log.Printf("Specified appPath to be %s", appPathRelative)
		log.Printf("Absolute appPath is %s", appPathAbsolute)
	}

	// First we'll check to see if the directory is empty. It's purely for safety purposes, to ensure we don't overwrite
	// anyting special. The command line handling has already validated that the folder actually exists
	isEmptyPath, err := isEmptyAppPath(appPathAbsolute)
	if !isEmptyPath || err != nil {
		log.Printf("The specified appPath '%s' does not appear to be an empty directory\n%v", appPathRelative, err)
		return err
	}

	// Scaffold
	err = scaffoldNewApp(appPathAbsolute)
	if err != nil {
		return err
	}

	log.Print("All done!")
	return nil
}

func isEmptyAppPath(appPath string) (bool, error) {
	// See http://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty

	// Open the directory, which must not fail
	f, err := os.Open(appPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// See if there's anything in the directory at all
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
