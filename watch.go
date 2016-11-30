package main

import (
	"log"
	"path"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/livereload"
	config "github.com/Travix-International/Travix.Core.Adk/models/config"
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

type pushCommand int

const (
	pushOpenBrowser pushCommand = iota
	pushDontOpenBrowser
	pushDone
)

// TODO: REALLY RAW IMPLEMENTATION! Needs worker pool pattern
func pusherController(cfg *config.Config, pushStatus chan bool) chan pushCommand {
	ch := make(chan int)

	go func() {
	Done:
		for {
			select {
			case cmd := <-command:
				pushCmd := &PushCommand{}

				pushCmd.AppPath = appPath
				pushCmd.NoPolling = false
				pushCmd.WaitInSeconds = 180

				switch cmd {
				case pushOpenBrowser:
					pushCmd.NoBrowser = true
					pushCmd.Push(cfg)

				case pushDonOpenBrowser:
					pushCmd.NoBrowser = false
					pushCmd.Push(cfg)
					livereload.SendReload()

				case pushDone:
					break Done
				}
			}
			pushStatus <- true
		}
	}()

	return ch
}

func registerWatch(app *kingpin.Application, cfg *config.Config) {
	command := app.Command("watch", "Watches the current directory for changes, and pushes on any change.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// Channel on which we get file change events.
			fileWatch := make(chan notify.EventInfo)
			// Channel on which we get an event when the initial short delay after a change is passed.
			initialDelayDone := make(chan int)
			// Channel on which we get events when the pushes are done.
			pushStatus := make(chan bool)
			// Channel to control the pushing
			controller := pusherController(cfg, pushStatus)

			// NOTE: We need to convert to absolute path, because the file watcher wouldn't accept relative paths on Windows.
			absPath, err := filepath.Abs(appPath)

			if err != nil {
				log.Fatal(err)
			}

			if err := notify.Watch(path.Join(absPath, "..."), fileWatch, notify.All); err != nil {
				log.Fatal(err)
			}

			defer notify.Stop(fileWatch)

			livereload.StartServer()

			// Immediately push once, and then start watching.
			controller <- pushOpenBrowser

			livereload.SendReload()

			// Infinite loop, the user can exit with Ctrl+C.
			for {
				select {
				case ei := <-fileWatch:
					if cfg.Verbose {
						log.Println("File change event details:", ei)
					}

					if watcherState == waiting {
						watcherState = initialDelay

						time.AfterFunc(100*time.Millisecond, func() {
							initialDelayDone <- 0
						})
					} else if watcherState == pushing {
						watcherState = pushingAndGotEvent
					}
				case <-initialDelayDone:
					watcherState = pushing

					log.Println("File change detected, executing appix push.")
					controller <- pushDontOpenBrowser
				case <-pushStatus:
					if watcherState == pushingAndGotEvent {
						// A change event arrived while the previous push was happening, we push again.
						watcherState = pushing
						controller <- pushDontOpenBrowser
					} else {
						watcherState = waiting
						log.Println("Push done, watching for file changes.")
					}
				}
			}

			return nil
		})

	command.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&appPath)
	command.Flag("noBrowser", "Appix won't open the frontend in the browser after every push.").
		Default("false").
		BoolVar(&noBrowser)
}
