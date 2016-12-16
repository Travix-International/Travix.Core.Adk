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
func Ignore() func(string, bool) (bool, bool) {
	ignoredSubFolders := make(map[string]struct{})

	// TODO: memoize this function
	return func(path string, isDir bool) (bool, bool) {
		filename := filepath.Base(path)
		dir := filepath.Dir(path)

		if _, ok := ignoredSubFolders[dir]; ok {
			if isDir {
				ignoredSubFolders[path] = struct{}{}
			}
			return true, true
		}

		for _, ignored := range ignoredFileNames {
			if strings.EqualFold(filename, ignored) {
				if isDir {
					ignoredSubFolders[path] = struct{}{}
				}
				return true, false
			}
		}

		return false, false
	}
}
