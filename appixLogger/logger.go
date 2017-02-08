package appixLogger

import (
	"log"

	"time"

	"sync"

	"github.com/Travix-International/logger"
)

var myLogger *logger.Logger

var once sync.Once

type LoggerNotification struct {
	Message string
	Action  string
	Type    string
}

var LoggerNotificationQueue chan LoggerNotification
var Quit chan bool

const FROGGER_URL = "https://frogger.staging.travix.com/logs/totolog"

func createHttpTransport() *logger.Transport {
	formatter := logger.NewJSONFormat()
	transport := logger.NewHttpTransport(FROGGER_URL, formatter)

	return transport
}

func getDefaultMeta(messageType string, applicationGroup string) map[string]string {
	defaultMeta := make(map[string]string)

	if len(applicationGroup) == 0 {
		applicationGroup = "core"
	}

	defaultMeta["messageType"] = messageType
	defaultMeta["applicationgroup"] = applicationGroup
	defaultMeta["applicationname"] = "appix"

	return defaultMeta
}

func loggy(notification LoggerNotification) chan bool {
	done := make(chan bool)

	go func(n LoggerNotification) {
		var err error

		if n.Type == "error" {
			err = myLogger.ErrorWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
		} else {
			err = myLogger.InfoWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
		}

		if err != nil {
			log.Printf("An error occured when trying to log error: %s\n", err.Error())
		}
		done <- true
	}(notification)

	return done
}

func AddMessageToQueue(notification LoggerNotification) {
	LoggerNotificationQueue <- notification
}

func Start() {
	go func() {
		for {
			select {
			case notification := <-LoggerNotificationQueue:
				<-loggy(notification)
			case <-time.After(500 * time.Millisecond):
				return
			}
		}
	}()
}

func Stop() chan bool {
	done := make(chan bool)

	go func() {
		close(Quit)
		done <- true
	}()
	return done
}

func NewAppixLogger() {
	once.Do(func() {
		LoggerNotificationQueue = make(chan LoggerNotification)
		Quit = make(chan bool)
		meta := make(map[string]string)

		myLogger, _ = logger.New(meta)
		myLogger.AddTransport(createHttpTransport())
	})
}
