package appixLogger

/**
 * Logger singleton
 *
 * usage: in your file declare a variable of type *Logger
 * example:
 *  var logger = appixLogger.NewAppixLogger()
 *
 * If you use it in command add the following line at the beginning of the code of your command
 * example:
 *  defer logger.Stop()
 *
 * Log something in your code:
 * example:
 *  if err != nil {
 *   logger.AddMessageToQueue(appixLogger.LoggerNotification{
 *			Type:    "error",
 *			Message: fmt.Sprintf("Error here is the message: %s", err.Error()),
 *			Action:  "myAction",
 *   })
 *  }
 */

import (
	"log"

	"sync"

	loggy "github.com/Travix-International/logger"
)

// LoggerNotification : The structure describing a notification to log
type LoggerNotification struct {
	Message string
	Action  string
	Type    string
}

// Logger : The structure of the logger singleton
type Logger struct {
	Loggy                   *loggy.Logger
	LoggerNotificationQueue chan LoggerNotification
	Quit                    chan bool
	loggerURL               string
}

var once sync.Once
var instance *Logger

func createHTTPTransport(url string) *loggy.Transport {
	formatter := loggy.NewJSONFormat()
	transport := loggy.NewHttpTransport(url, formatter)

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

func (l *Logger) log(n LoggerNotification) {
	var err error

	if n.Type == "error" {
		err = l.Loggy.ErrorWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
	} else {
		err = l.Loggy.InfoWithMeta(n.Action, n.Message, getDefaultMeta(n.Action, ""))
	}

	if err != nil {
		log.Printf("An error occured when trying to log error: %s\n", err.Error())
	}
	// done <- true
}

// AddMessageToQueue : Add a new LoggerNotification object to the Queue and print on stdout the message
func (l *Logger) AddMessageToQueue(notification LoggerNotification) {
	// log on stdout to kkep the user aware of what's going on
	log.Printf("[appix:%s] %s\n", notification.Action, notification.Message)

	if l.Loggy != nil {
		l.LoggerNotificationQueue <- notification
	}
}

// Start : launch to go routine watching at the queue
func (l *Logger) Start() {
	go func() {
		for {
			select {
			case notification := <-l.LoggerNotificationQueue:
				// done := make(chan bool)
				l.log(notification)
				// <-done
			case <-l.Quit:
				return
			}
		}
	}()
}

// Stop : kill the logger routine
func (l *Logger) Stop() {
	l.Quit <- true
}

// NewAppixLogger : create a new instance of Logger if doesn't exist already. Otherwise return the actual instance
func NewAppixLogger(url string) *Logger {
	once.Do(func() {
		meta := make(map[string]string)
		myLogger, _ := loggy.New(meta)

		if myLogger != nil {
			myLogger.AddTransport(createHTTPTransport(url))
		}

		instance = &Logger{
			LoggerNotificationQueue: make(chan LoggerNotification),
			Quit:  make(chan bool),
			Loggy: myLogger,
		}
	})
	return instance
}
