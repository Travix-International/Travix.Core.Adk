package appcatalog

import (
	"log"
	"net/http"
)

func logServerResponse(res *http.Response) {
	log.Printf("Server response: %s\n", res.Request.URL)
	logHeaders(res)
}

func logHeaders(res *http.Response) {
	log.Println("\tHeaders:")
	for k, v := range res.Header {
		log.Printf("\t%s: %s\n", k, v)
	}
	log.Println("")
}
