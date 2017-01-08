package appix

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
	"sync"
)

// Stage : the structure representing the data model for using stages
type Stage struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
}

var pool chan bool
var wg sync.WaitGroup

// startStage : executes the command for the given kind of test
func startStage(name string, cmd string) {
	defer wg.Done()
	// format the command in order to use os/exec/Command
	firstSpace := strings.Index(cmd, " ")

	var command string
	var args []string

	if firstSpace > 0 {
		command = cmd[:firstSpace]
		args = strings.Split(cmd[firstSpace:len(cmd)], " ")
	} else {
		command = cmd
	}

	cm := exec.Command(command)

	if len(args) > 0 {
		cm.Args = args
	}

	// define and set an output buffer.
	// TODO: do we display the logs all the time or do we consider --verbose as a flag?
	var out bytes.Buffer
	cm.Stdout = &out

	// execute the command
	err := cm.Run()

	if err != nil {
		log.Printf("The stage '%s' failed: %s\n", name, err.Error())
		pool <- false
		return
	}
	log.Printf("Stage '%s' done.", name)
	pool <- true
}

// CreateStagePool : create the pool of stage. Each stage must be independent since there is no synchronisation between the routines
func CreateStagePool(stages []Stage) chan bool {
	pool = make(chan bool, len(stages))
	wg.Add(len(stages))

	// the routine taking care of the result set in the pool
	go func() error {
		for status := range pool {
			if !status {
				return errors.New("An error occured")
			}
		}
		return nil
	}()

	// fill the pool
	for _, s := range stages {
		go startStage(s.Name, s.Cmd)
	}

	// close
	wg.Wait()
	close(pool)

	return pool
}
