package main

import (
	"log"
	"path"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/rjeczalik/notify"
)

var (
	appPath   string
	noBrowser bool
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

		log.Println("File change detected, executing appix push.")

		if verbose {
			log.Println("File change event details:", ei)
		}

		doPush(context, false)

		sendReload()

		log.Println("Push done, watching for file changes.")
	}
}

func doPush(context *kingpin.ParseContext, openBrowser bool) {
	pushCmd := &PushCommand{}

	pushCmd.appPath = appPath
	pushCmd.noPolling = false
	pushCmd.waitInSeconds = 180
	pushCmd.noBrowser = !openBrowser

	pushCmd.push(context)
}
