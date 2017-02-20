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

// SubmitToCatalog submits the specified app to the AppCatalog.
func SubmitToCatalog(submitURI string, timeout int, appManifestFile string, zapFile string, verbose bool, config config.Config, logger *appixLogger.Logger) (acceptanceQueryURL string, err error) {
	var req *http.Request
	files := map[string]string{
		"manifest": appManifestFile,
		"zapfile":  zapFile,
	}

	for attempt := 1; attempt <= config.MaxRetryAttempts; attempt++ {
		if req, err = prepare(submitURI, files, config, verbose, logger); err != nil {
			return "", err
		}

		log.Printf("Submitting files to App Catalog. Attempt %v of %v\n", attempt, config.MaxRetryAttempts)

		if acceptanceQueryURL, err = doRequest("Submit", "acc:query", req, time.Duration(timeout)*time.Second, verbose); err == nil {
			break
		}

		if err, ok := err.(*catalogError); ok && !err.canRetry() {
			if err.authenticationIssue() {
				log.Printf("You are not authorized to submit the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", err.statusCode)
				return "", fmt.Errorf("Authentication error")
			}
			// log.Print(err.Error())
			break
		}

		if attempt < config.MaxRetryAttempts {
			wait := math.Pow(2, float64(attempt-1)) * 1000
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

	return
}
