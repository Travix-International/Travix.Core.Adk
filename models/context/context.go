package context

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	modelsConfig "github.com/Travix-International/Travix.Core.Adk/models/config"
)

type Context struct {
	App    *kingpin.Application
	Config modelsConfig.Config
}
