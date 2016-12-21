package appcatalog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type bundleMessage struct {
	Widget string
	Output string
}

type pushPollResponse struct {
	Meta struct {
		Status   string
		Messages []bundleMessage
	}
	Links struct {
		Preview string
	}
}

const (
	pollClientTimeout  = 5 * time.Second
	pollInterval       = 5 * time.Second // how often to poll status URL
	pollFinishedStatus = "FINISHED"
	pollFailedStatus   = "FAILED"
)

// PollUntilDone polls the frontend until the bundling is done, failed, or the timeout passed.
func PollUntilDone(pollURI string, waitInSeconds int, openBrowser bool, verbose bool, openURL func(string) error) {
	if verbose {
		log.Println("Entering polling routine")
	}
	quit := make(chan interface{}, 1)
	defer close(quit)

	progressMonitor := verifyProgress(pollURI, quit)
	wait := time.Duration(waitInSeconds) * time.Second

	select {
	case statusResponse, ok := <-progressMonitor:
		if !ok {
			break
		}

		log.Printf("Server output for the app bundling:")
		for _, message := range statusResponse.Meta.Messages {
			log.Printf("Widget: %s", message.Widget)
			log.Printf("Output: %s", message.Output)
		}

		if statusResponse.Meta.Status == pollFinishedStatus {
			log.Printf("App successfully pushed. The frontend for this development session is at %s", statusResponse.Links.Preview)
			if openBrowser {
				openURL(statusResponse.Links.Preview)
			}
		} else {
			log.Printf("App push failed.")
		}

		close(progressMonitor)

	case <-time.After(wait):
		quit <- true // send a cancel signal to progressMonitor
		log.Printf("Operation timed out after %s", wait)
	}
}

func verifyProgress(pollURI string, quit <-chan interface{}) chan pushPollResponse {
	done := make(chan pushPollResponse, 1)
	go func() {
		var statusResponse pushPollResponse
		timeout := time.Duration(pollClientTimeout)
		client := http.Client{Timeout: timeout}

		for {
			// check if operation should be cancelled
			select {
			case <-quit:
				return
			default:
			}

			resp, err := client.Get(pollURI)
			if err != nil {
				log.Println("Error during polling the bundling status.")
				log.Println(err)
				close(done)
				return
			}

			err = json.NewDecoder(resp.Body).Decode(&statusResponse)
			resp.Body.Close()

			if err != nil {
				log.Println("Error. during parsing poll status result")
				bodyData, _ := ioutil.ReadAll(resp.Body)
				if bodyData != nil {
					log.Println(bodyData)
				}
				close(done)
				return
			}

			log.Printf("Pushing to the website to the development environment, status: [%s]", statusResponse.Meta.Status)

			if statusResponse.Meta.Status == pollFinishedStatus || statusResponse.Meta.Status == pollFailedStatus {
				done <- statusResponse
				break
			}

			time.Sleep(pollInterval)
		}
	}()
	return done
}
