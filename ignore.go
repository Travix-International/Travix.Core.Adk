package appix

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Travix-International/appix/config"
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

func IgnoreFilePath(path string) (ignored bool, ignoredFolder bool) {
	fileName := filepath.Base(path)
	dir := filepath.Dir(path)

	if fileName == "." {
		ignored = false
		ignoredFolder = false
		return
	}

	for _, ignoredFileName := range ignoredFileNames {
		if strings.EqualFold(fileName, ignoredFileName) {
			ignored = true
			break
		} else if strings.Contains(dir, ignoredFileName + string(os.PathSeparator)) {
			ignoredFolder = true
			break
		}
	}

	return
}

func init() {
	file, err := os.Open(config.IgnoreFileName)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoredFileName := scanner.Text()
		if ignoredFileName != "" {
			ignoredFileNames = append(ignoredFileNames, ignoredFileName)
		}
	}
}
