package appixLogger

import (
	"log"

	"sync"

	loggy "github.com/Travix-International/logger"
)

type LoggerNotification struct {
	Message string
	Action  string
	Type    string
}

type Logger struct {
	Loggy                   *loggy.Logger
	LoggerNotificationQueue chan LoggerNotification
	Quit                    chan bool
}

const FROGGER_URL = "https://frogger.travix.com/logs/appixlog"

var once sync.Once
var instance *Logger

func createHTTPTransport() *loggy.Transport {
	formatter := loggy.NewJSONFormat()
	transport := loggy.NewHttpTransport(FROGGER_URL, formatter)

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

func (l *Logger) log(notification LoggerNotification, done chan bool) {
	go func(n LoggerNotification) {
		var err error

		if n.Type == "error" {
			err = l.Loggy.ErrorWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
		} else {
			err = l.Loggy.InfoWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
		}

		if err != nil {
			log.Printf("An error occured when trying to log error: %s\n", err.Error())
		}
		done <- true
	}(notification)
}

func (l *Logger) AddMessageToQueue(notification LoggerNotification) {
	if l.Loggy != nil {
		l.LoggerNotificationQueue <- notification
	} else {
		log.Printf("[appix:%s] %s\n", notification.Action, notification.Message)
	}
}

func (l *Logger) Start() {
	go func() {
		for {
			select {
			case notification := <-l.LoggerNotificationQueue:
				done := make(chan bool)
				l.log(notification, done)
				<-done
			case <-l.Quit:
				return
			}
		}
	}()
}

func (l *Logger) Stop() {
	go func() {
		l.Quit <- true
	}()
}

func NewAppixLogger() *Logger {
	once.Do(func() {
		meta := make(map[string]string)
		myLogger, _ := loggy.New(meta)

		if myLogger != nil {
			myLogger.AddTransport(createHTTPTransport())
		}

		instance = &Logger{
			LoggerNotificationQueue: make(chan LoggerNotification),
			Quit:  make(chan bool),
			Loggy: myLogger,
		}
	})
	return instance
}
