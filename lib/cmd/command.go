package cmd

import (
	"github.com/Travix-International/Travix.Core.Adk/lib/context"
)

type Command struct {
	Verbose       bool
	LocalFrontend bool
	TargetEnv     string
}

type Registrable interface {
	Register(context.Context)
}
