package utils

import (
	"os/exec"
)

func OpenUrl(url string) error {
	return exec.Command("xdg-open", url).Start()
}
