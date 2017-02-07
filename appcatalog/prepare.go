package appcatalog

import (
	"log"
	"net/http"
	"os"

	"github.com/Travix-International/appix/auth"
	"github.com/Travix-International/appix/config"
)

func prepare(uri string, files map[string]string, config config.Config, verbose bool) (req *http.Request, err error) {
	if verbose {
		log.Println("Posting files to the App Catalog: " + uri)
	}

	req, err = CreateMultiFileUploadRequest(uri, files, nil, verbose)
	if err != nil {
		log.Println("Creating the HTTP request failed.")
		return
	}

	err = addAuthenticationHeader(req, config)
	return
}

func addAuthenticationHeader(req *http.Request, config config.Config) error {
	token, err := auth.LoadAuthToken(config)

	if err == nil {
		req.Header.Set("Authorization", token.TokenType+" "+token.IdToken)
		return nil
	}

	log.Println("WARNING: You are not logged in. In a future version authentication will be mandatory.\nYou can log in using \"appix login\".")

	// we can safely ignore path errors (e.g. auth.json doesn't exist)
	if _, ok := err.(*os.PathError); ok {
		return nil
	}

	return err
}
