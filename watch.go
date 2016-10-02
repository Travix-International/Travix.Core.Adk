package main

import (
	"log"
	"path"

	"github.com/rjeczalik/notify"
	"gopkg.in/alecthomas/kingpin.v2"
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

	if err := notify.Watch(path.Join(appPath, "..."), c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Immediately push once, and then start watching.
	doPush(context)

	// Infinite loop, the user can exit with Ctrl+C
	for {
		// Block until an event is received.
		ei := <-c

		log.Println("File change detected, executing appix push.")

		if verbose {
			log.Println("File change event details:", ei)
		}

		doPush(context)

		log.Println("Push done, watching for file changes.")
	}
}

func doPush(context *kingpin.ParseContext) {
	pushCmd := &PushCommand{}

	pushCmd.appPath = appPath
	pushCmd.noPolling = false
	pushCmd.waitInSeconds = 180
	pushCmd.noBrowser = noBrowser

	pushCmd.push(context)
}
