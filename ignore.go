package appix

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Travix-International/appix/config"
	"github.com/ryanuber/go-glob"
)

// Supported glob patterns:
// filename          match the filename in the root path, e.g. ./node_modules
// */filename        match filename (depth: exactly 2), e.g. ./ui/node_modules
// */*/filename      match filename (depth: exactly 3), e.g. ./ui/widget/node_modules
// **/filename       match filename (depth: n | n > 0, i.e. does not include root), e.g. /any/depth/greather/than/0/node_modules
// **/filename/*     match any file inside filename (depth > 0).  This is important because of watch

var ignoredFileNames = []string{
	"node_modules",
	"**/node_modules",
	"**/node_modules/*",
	".git",
	".gitignore",
	"**/.gitkeep",
	"temp",
	"npm-debug.log",
	"**/npm-debug.log",
	".idea",
	"**/.idea",
	"**/.idea/*",
	".vscode",
	"**/.vscode",
	"**/.vscode/*",
	".ds_store",
	"**/.ds_store",
	"thumbs.db",
	"**/thumbs.db",
	"desktop.ini",
	"**/desktop.ini",
	config.DevFileName,
	config.IgnoreFileName,
}

func IgnoreFilePath(path string) (ignored bool) {
	fileName := filepath.Base(path)

	if fileName == "." {
		ignored = false
		return
	}

	for _, ignoredFileName := range ignoredFileNames {
		if glob.Glob(ignoredFileName, path) {
			ignored = true
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
		ignoredFileName := strings.Trim(scanner.Text(), " ")
		if ignoredFileName != "" && strings.HasPrefix(ignoredFileName, "#") {
			ignoredFileNames = append(ignoredFileNames, ignoredFileName)
		}
	}
}
