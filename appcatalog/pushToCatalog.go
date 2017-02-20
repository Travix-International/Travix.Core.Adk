package appcatalog

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

const (
	timeout = time.Duration(10 * time.Second)
)

// PushToCatalog pushes the specified app to the AppCatalog.
func PushToCatalog(pushURI string, timeout int, appManifestFile string, verbose bool, config config.Config, logger *appixLogger.Logger) (uploadURI string, err error) {
	var req *http.Request
	files := map[string]string{
		"manifest": appManifestFile,
	}

	for attempt := 1; attempt <= config.MaxRetryAttempts; attempt++ {
		if req, err = prepare(pushURI, files, config, verbose, logger); err != nil {
			return "", err
		}

		log.Printf("Pushing files to catalog. Attempt %v of %v\n", attempt, config.MaxRetryAttempts)

		if uploadURI, err = doRequest("Push", "upload", req, time.Duration(timeout)*time.Second, verbose); err == nil {
			log.Println("App has been pushed successfully.")
			break
		}

		if err, ok := err.(*catalogError); ok && !err.canRetry() {
			if err.authenticationIssue() {
				log.Printf("You are not authorized to push the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", err.statusCode)
				return "", fmt.Errorf("Authentication error")
			}
			break
		}

		if attempt < config.MaxRetryAttempts {
			wait := math.Pow(2, float64(attempt-1)) * 1000
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

	return
}
