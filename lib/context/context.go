package context

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/config"
)

type Context struct {
	App    *kingpin.Application
	Config *config.Config
}
