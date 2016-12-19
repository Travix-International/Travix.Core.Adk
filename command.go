package appix

import (
	"github.com/Travix-International/appix/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	Verbose       bool
	LocalFrontend bool
	TargetEnv     string
}

type Registrable interface {
	Register(app *kingpin.Application, config config.Config)
}
