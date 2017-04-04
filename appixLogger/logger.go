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
 *			Type:    logger.LevelError,
 *			Message: fmt.Sprintf("Error here is the message: %s", err.Error()),
 *			Action:  "myAction",
 *   })
 *  }
 */

import (
	"log"

	loggy "github.com/Travix-International/logger"
)

const (
	// LevelError is the error level
	LevelError = "error"
	// LevelInfo is the info level
	LevelInfo = "info"
)

// LoggerNotification : The structure describing a notification to log
type LoggerNotification struct {
	Message  string
	LogEvent string
	Level    string
}

// AppixLogger defines the core set of logging calls, to be used by other components
type AppixLogger interface {
	AddMessageToQueue(notification LoggerNotification)
}

// Logger : The structure of the logger singleton
type Logger struct {
	loggerImpl              *loggy.Logger
	loggerNotificationQueue chan LoggerNotification
	quit                    chan bool
	loggerURL               string
}

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

	if n.Level == "error" {
		err = l.loggerImpl.ErrorWithMeta(n.LogEvent, n.Message, getDefaultMeta(n.LogEvent, ""))
	} else {
		err = l.loggerImpl.InfoWithMeta(n.LogEvent, n.Message, getDefaultMeta(n.LogEvent, ""))
	}

	if err != nil {
		log.Printf("An error occured when trying to log error: %s\n", err.Error())
	}
}

// AddMessageToQueue : Add a new LoggerNotification object to the Queue and print on stdout the message
func (l *Logger) AddMessageToQueue(notification LoggerNotification) {
	// log on stdout to kkep the user aware of what's going on
	log.Printf("%s: %s\n", notification.LogEvent, notification.Message)

	if l.loggerImpl != nil {
		l.loggerNotificationQueue <- notification
	}
}

// Start : launch to go routine watching at the queue
func (l *Logger) Start() {
	go func() {
		for {
			select {
			case notification := <-l.loggerNotificationQueue:
				l.log(notification)
			case <-l.quit:
				return
			}
		}
	}()
}

// Stop : kill the logger routine
func (l *Logger) Stop() {
	l.quit <- true
}

// NewAppixLogger : create a new instance of Logger if doesn't exist already. Otherwise return the actual instance
func NewAppixLogger(url string) *Logger {
	meta := make(map[string]string)
	myLogger, _ := loggy.New(meta)

	if myLogger != nil {
		myLogger.AddTransport(createHTTPTransport(url))
	}

	return &Logger{
		// CR JP: Note how we use an unbuffered channel for the notification queue. This has its advantages
		// (easy) at the expense of concurrent processes possibly having to halt due to this. Also, any
		// routine that logs before Start() or after Stop() will deadlock.
		loggerNotificationQueue: make(chan LoggerNotification),
		quit:       make(chan bool),
		loggerImpl: myLogger,
	}
}
