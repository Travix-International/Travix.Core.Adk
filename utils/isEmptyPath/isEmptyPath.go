package isEmptyPath

import (
	"io"
	"os"
)

func IsEmptyPath(appPath string) (bool, error) {
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
