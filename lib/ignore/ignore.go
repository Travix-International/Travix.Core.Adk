package ignore

import (
	"path/filepath"
	"strings"

	"github.com/Travix-International/Travix.Core.Adk/lib/config"
)

var ignoredFileNames = []string{
	"node_modules",
	"temp",
	".git",
	".idea",
	".vscode",
	".ds_store",
	"thumbs.db",
	"desktop.ini",
	config.DevFileName,
	config.IgnoreFileName,
}

// TODO: read from .appixignore file
func IgnoreFilePath(path string) (ignored bool, ignoredFolder bool) {
	var recurse func(string) (bool, bool)
	recurse = func(path string) (ignored bool, ignoredFolder bool) {
		fileName := filepath.Base(path)

		if fileName == "." {
			ignored = false
			ignoredFolder = false
			return
		}

		for _, ignoredFileName := range ignoredFileNames {
			if strings.EqualFold(fileName, ignoredFileName) {
				ignored = true
				break
			}
		}

		dir := filepath.Dir(path)
		ignoredFolder, _ = recurse(dir)
		ignored = ignored || ignoredFolder
		return
	}

	ignored, ignoredFolder = recurse(path)
	return
}
