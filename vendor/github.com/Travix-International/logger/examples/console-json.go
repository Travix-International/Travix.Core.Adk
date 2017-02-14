package main

import (
	"fmt"
	"os"

	"github.com/Travix-International/logger"
)

func main() {
	log := logger.New()

	// transString := logger.NewTransport(os.Stdout, logger.DefaultStringFormat)

	transJson := logger.NewTransport(os.Stdout, logger.DefaultJSONFormat)

	log.AddTransport(logger.ConsoleTransport, transJson)

	log.Meta["oink"] = "swine"

	_, ok := log.Meta["notHere"]

	if !ok {
		log.Meta["didntfindHere"] = "true"
	}

	e := log.Infof("snorting", "Piggies went snorting %d times over the %s", 3, "yard")

	if e != nil {
		fmt.Println(e)
	}
}
