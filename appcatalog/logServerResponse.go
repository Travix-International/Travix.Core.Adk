package appcatalog

import (
	"log"
	"net/http"
	"sort"
)

func logServerResponse(res *http.Response) {
	log.Printf("\t%s %s\n", res.Request.Method, res.Request.URL)
	log.Printf("\t%s %s\n", res.Proto, res.Status)
	keys := make([]string, 0, len(res.Header))
	for k := range res.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		log.Printf("\t%s: %s\n", k, res.Header.Get(k))
	}
	log.Println("")
}
