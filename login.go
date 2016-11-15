package main

import (
	"fmt"
	"io"
	"net/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureLoginCommand(app *kingpin.Application) {
	app.Command("login", "Login").
		Action(executeLoginCommand)
}

func executeLoginCommand(context *kingpin.ParseContext) error {
	fmt.Println("Login here...")
	ch := make(chan bool)

	go startLoginServer(ch)

	openWebsite("http://localhost:7001")

	select {
	case shouldClose := <-ch:
		if shouldClose {

		} else {

		}
	}

	return nil
}

func loginBaseRoute(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world from login server!")
}

func startLoginServer(c chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world from login server!")
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "closing...")

		c <- true
	})

	http.ListenAndServe(":7001", nil)
}

