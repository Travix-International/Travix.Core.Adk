package ignore

import (
	"strings"
)

func Includes(appixIgnoreFileName string, devFileName string) func(string, bool) bool {
	return func(filePath string, isDir bool) bool {
		if isDir {
			filePath += "/"
		}

		path := strings.ToLower(filePath)
		canInclude := !strings.Contains(path, "/node_modules/") &&
			!strings.Contains(path, "/temp/") &&
			!strings.Contains(path, ".git/") &&
			!strings.HasSuffix(path, ".idea/") &&
			!strings.HasSuffix(path, ".vscode/") &&
			!strings.HasSuffix(path, ".ds_store") &&
			!strings.HasSuffix(path, "thumbs.db") &&
			!strings.HasSuffix(path, appixIgnoreFileName) &&
			!strings.HasSuffix(path, devFileName) &&
			!strings.HasSuffix(path, "desktop.ini")

		return canInclude
	}
}

func Ignores(appixIgnoreFileName string, devFileName string) func(string, bool) bool {
	shouldInclude := Includes(appixIgnoreFileName, devFileName)
	return func(path string, isDir bool) bool {
		return !shouldInclude(path, isDir)
	}
}
