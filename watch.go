package main

import (
	"log"
	"path"
	"path/filepath"
	"sync"
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
	mutex        = &sync.Mutex{}
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
	// NOTE: The second argument controls the buffer length.
	// 1 is ideal, because if any number of change happen during the push, we want to push one more time afterwards.
	// But we don't want to push for each intermediate change. (So let's say there were 3 more file changes during the push. Afterwards we want to push only once, and not three more times.)
	c := make(chan notify.EventInfo, 1)

	// NOTE: We need to convert to absolute path, because the file watcher wouldn't accept relative paths on Windows.
	absPath, err := filepath.Abs(appPath)

	if err != nil {
		log.Fatal(err)
	}

	if err := notify.Watch(path.Join(absPath, "..."), c, notify.All); err != nil {
		log.Fatal(err)
	}

	defer notify.Stop(c)

	startLivereloadServer()

	// Immediately push once, and then start watching.
	doPush(context, true)

	sendReload()

	// Infinite loop, the user can exit with Ctrl+C
	for {
		// Block until an event is received.
		ei := <-c

		if verbose {
			log.Println("File change event details:", ei)
		}

		mutex.Lock()
		if watcherState == waiting {
			watcherState = initialDelay
			log.Println("File change detected, executing appix push.")
			go doPush(context, false)
		} else if watcherState == pushing {
			watcherState = pushingAndGotEvent
		}
		mutex.Unlock()
	}
}

func doPush(context *kingpin.ParseContext, openBrowser bool) {
	time.Sleep(100 * time.Millisecond)

	watcherState = pushing

	for watcherState == pushing {
		pushCmd := &PushCommand{}

		pushCmd.appPath = appPath
		pushCmd.noPolling = false
		pushCmd.waitInSeconds = 180
		pushCmd.noBrowser = !openBrowser

		pushCmd.push(context)

		if !openBrowser {
			sendReload()
		}

		mutex.Lock()
		if watcherState == pushingAndGotEvent {
			// A change event arrived while the previous push was happening, we push again.
			watcherState = pushing
			mutex.Unlock()
		} else {
			watcherState = waiting
			mutex.Unlock()
			log.Println("Push done, watching for file changes.")
		}
	}
}
