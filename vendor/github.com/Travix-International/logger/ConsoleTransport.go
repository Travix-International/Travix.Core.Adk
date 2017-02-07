package logger

import (
	"os"
)

var ConsoleTransport = NewTransport(os.Stdout, DefaultStringFormat)
