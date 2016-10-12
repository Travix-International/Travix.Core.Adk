package main

import (
	"log"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/rjeczalik/notify"
)

// This watcher implements a simple state machine, making sure we handle currently if change events come in while we are executing a push.
//
// NOTE: The file watcher libraries sometimes send two separate events for one file change in quick succession. (Also, some editors, like vim, are doing multiple genuine file modifications for one single file save.)
// To mitigate this we initially wait for a short while befor starting the push, to make sure we are not pushing twice for a single change. That's why we have the initialDelay state.
//
//                              file change event
// initial state                     received
//   -------------> WAITING ------------------------> INITIAL_DELAY
//                     Λ                                    |
//                     |                                    | 100ms passed, executing push
//                     |                                    |
//                     |          push completed            V
//                      -------------------------------- PUSHING
//                                                        Λ   |
//                                         push completed |   | file change event received
//                                     execute a new push |   |
//                                                        |   V
//                                                 PUSHING_AND_GOT_EVENT
//
const (
	waiting            = iota
	initialDelay       = iota
	pushing            = iota
	pushingAndGotEvent = iota
)

var (
	appPath      string
	noBrowser    bool
	watcherState = waiting
)

func configureWatchCommand(app *kingpin.Application) {
	appCmd := app.Command("watch", "Watches the current directory for changes, and pushes on any change.").
		Action(executeWatchCommand)

	appCmd.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&appPath)
	appCmd.Flag("noBrowser", "Appix won't open the frontend in the browser after every push.").
		Default("false").
		BoolVar(&noBrowser)
}

func executeWatchCommand(context *kingpin.ParseContext) error {
	// Channel on which we get file change events.
	fileWatch := make(chan notify.EventInfo)
	// Channel on which we get an event when the initial short delay after a change is passed.
	initialDelayDone := make(chan int)
	// Channel on which we get events when the pushes are done.
	buildDone := make(chan int)

	// NOTE: We need to convert to absolute path, because the file watcher wouldn't accept relative paths on Windows.
	absPath, err := filepath.Abs(appPath)

	if err != nil {
		log.Fatal(err)
	}

	if err := notify.Watch(path.Join(absPath, "..."), fileWatch, notify.All); err != nil {
		log.Fatal(err)
	}

	defer notify.Stop(fileWatch)

	startLivereloadServer()

	// Immediately push once, and then start watching.
	doPush(context, true, nil)

	sendReload()

	// Infinite loop, the user can exit with Ctrl+C.
	for {
		select {
		case ei := <-fileWatch:
			if verbose {
				log.Println("File change event details:", ei)
			}

			if watcherState == waiting {
				watcherState = initialDelay

				go waitForDelay(initialDelayDone)
			} else if watcherState == pushing {
				watcherState = pushingAndGotEvent
			}
		case _ = <-initialDelayDone:
			watcherState = pushing

			log.Println("File change detected, executing appix push.")

			go doPush(context, false, &buildDone)
		case _ = <-buildDone:
			if watcherState == pushingAndGotEvent {
				// A change event arrived while the previous push was happening, we push again.
				watcherState = pushing
				go doPush(context, false, &buildDone)
			} else {
				watcherState = waiting
				log.Println("Push done, watching for file changes.")
			}
		}
	}
}

func waitForDelay(delayDone chan int) {
	time.Sleep(100 * time.Millisecond)
	delayDone <- 0
}

func doPush(context *kingpin.ParseContext, openBrowser bool, buildDone *chan int) {
	pushCmd := &PushCommand{}

	pushCmd.appPath = appPath
	pushCmd.noPolling = false
	pushCmd.waitInSeconds = 180
	pushCmd.noBrowser = !openBrowser

	pushCmd.push(context)

	if !openBrowser {
		sendReload()
	}

	if buildDone != nil {
		*(buildDone) <- 0
	}
}
