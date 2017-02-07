package appix

import (
	"log"

	"github.com/Travix-International/logger"
)

type AppixLogger struct {
	myLogger *logger.Logger
	action   string
}

const FROGGER_URL = "https://frogger.travix.com/logs/totolog"

func createHttpTransport() *logger.Transport {
	formatter := logger.NewJSONFormat()
	transport := logger.NewHttpTransport(FROGGER_URL, formatter)

	return transport
}

func getDefaultMeta() map[string]string {
	defaultMeta := make(map[string]string)

	defaultMeta["messagetype"] = "AppixLog"
	defaultMeta["applicationgroup"] = "core"
	defaultMeta["applicationname"] = "appix"

	return defaultMeta
}

func (t AppixLogger) Error(message string) {
	err := t.myLogger.ErrorWithMeta(t.action, message, getDefaultMeta())

	if err != nil {
		log.Printf("An error occured when trying to log error: %s\n", err.Error())
	}
}

func (t AppixLogger) Log(message string) {
	err := t.myLogger.InfoWithMeta(t.action, message, getDefaultMeta())

	if err != nil {
		log.Printf("An error occured when trying to log info: %s\n", err.Error())
	}
}

func NewAppixLogger(action string) AppixLogger {
	meta := make(map[string]string)
	myLogger, err := logger.New(meta)

	if err != nil {
		log.Fatalf("An error occured while initialising the logger: %s\n", err.Error())
	}

	myLogger.AddTransport(createHttpTransport())
	myLogger.AddTransport(logger.ConsoleTransport)

	return AppixLogger{
		myLogger: myLogger,
		action:   action,
	}
}
