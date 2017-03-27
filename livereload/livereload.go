package livereload

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	hub *Hub
	// Appix always hosts the livereload server at 13221. This is the port the frontend has to try to use to connect to.
	livereloadAddress    = ":13221"
	livereloadAddressTLS = ":13222"
)

var (
	certContent string
	keyContent  string
)

func base64Decode(src string) []byte {
	val, err := base64.StdEncoding.DecodeString(certContent)

	if err != nil {
		log.Fatalf("An error occured while decoding the value. Details: %s\n", err.Error())
	}
	return val
}

func createCertFiles() (cert string, key string) {
	tempFolder, _ := ioutil.TempDir("", "appix")

	cert = tempFolder + "/livereload-cert.pem"
	key = tempFolder + "/livereload-key.pem"

	ioutil.WriteFile(cert, base64Decode(certContent), 0644)
	ioutil.WriteFile(key, base64Decode(keyContent), 0644)

	return cert, key
}

// StartServer starts the Websocket server listening for the websites that want to connect.
func StartServer() {
	hub = newHub()
	go hub.run()
	http.HandleFunc("/appixlivereload", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go startLocalServer()
	go startLocalServerTLS()
	log.Println("Livereload server listening at", livereloadAddress)
}

func startLocalServer() {
	err := http.ListenAndServe(livereloadAddress, nil)

	if err != nil {
		log.Println("Failed to start up the Livereload server: ", err)
		return
	}
}

func startLocalServerTLS() {
	cert, key := createCertFiles()
	err := http.ListenAndServeTLS(livereloadAddressTLS, cert, key, nil)

	if err != nil {
		log.Println("Failed to start up the Livereload server with TLS: ", err)
		return
	}
}

// SendReload sends a message to the Websocket listeners to refresh the page.
func SendReload() {
	message := bytes.TrimSpace([]byte("reload"))
	hub.Broadcast <- message
}
