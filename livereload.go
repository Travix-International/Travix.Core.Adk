package main

import (
	"bytes"
	"log"
	"net/http"
)

var (
	hub *Hub
	// Appix always hosts the livereload server at 13221. This is the port the frontend has to try to use to connect to.
	livereloadAddress = ":13221"
)

func startLivereloadServer() {
	hub = newHub()
	go hub.run()
	http.HandleFunc("/appixlivereload", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go runServer()
}

func runServer() {
	err := http.ListenAndServe(livereloadAddress, nil)
	log.Println("Livereload server listening at", livereloadAddress)

	if err != nil {
		log.Println("Failed to start up the Livereload server: ", err)
	}
}

func sendReload() {
	message := bytes.TrimSpace([]byte("reload"))
	hub.broadcast <- message
}
