package submit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/Travix.Core.Adk/lib/cmd"
	"github.com/Travix-International/Travix.Core.Adk/lib/context"
	"github.com/Travix-International/Travix.Core.Adk/lib/upload"
	"github.com/Travix-International/Travix.Core.Adk/lib/zapper"
)

// SubmitCommand used for submitting Apps
type SubmitCommand struct {
	*cmd.Command
	appPath     string // path to the App folder
	environment string // environment (default: dev)
}

type submitResponse struct {
	Messages []string
	Links    map[string]string
}

func (cmd *SubmitCommand) Register(context context.Context) {
	config := context.Config
	const submitTemplateURI = "%s/files/publish/%s"

	command := context.App.Command("submit", "Submits the App for review.").
		Action(func(parseContext *kingpin.ParseContext) error {
			context.LoadAuthToken()

			environment := cmd.environment

			if environment == "" {
				environment = "dev"
			}

			appPath, appName, appManifestFile, err := zapper.PrepareAppUpload(cmd.appPath)

			if err != nil {
				log.Println("Could not prepare the app folder for uploading")
				return err
			}

			zapFile, err := zapper.CreateZapPackage(appPath, config.DevFileName, cmd.Verbose)

			if err != nil {
				log.Println("Could not create zap package!")
				return err
			}

			log.Printf("Run submit for App '%s', env '%s', path '%s'\n", appName, environment, appPath)

			rootURI := config.CatalogURIs[cmd.TargetEnv]
			submitURI := fmt.Sprintf(submitTemplateURI, rootURI, appName)
			files := map[string]string{
				"manifest": appManifestFile,
				"zapfile":  zapFile,
			}

			if cmd.Verbose {
				log.Println("Posting files to App Catalog: " + submitURI)
			}
			request, err := upload.CreateMultiFileUploadRequest(submitURI, files, nil, cmd.Verbose)
			if err != nil {
				log.Println("Call to App Catalog failed!")
				return err
			}

			token, err := context.LoadAuthToken()

			if err == nil {
				request.Header.Set("Authorization", token.TokenType+" "+token.IdToken)
			} else {
				log.Println("WARNING: You are not logged in. In a future version authentication will be mandatory.\nYou can log in using \"appix login\".")
			}

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Println("Call to App Catalog failed!")
				return err
			}

			if response.StatusCode == 401 || response.StatusCode == 403 {
				log.Printf("You are not authorized to submit the application to the App Catalog (status code %v). If you are not signed in, please log in using 'appix login'.", response.StatusCode)
				return fmt.Errorf("Authentication error")
			}

			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println("Error reading response from App Catalog!")
				return err
			}

			var responseObject submitResponse
			err = json.Unmarshal(responseBody, &responseObject)
			if err != nil {
				if cmd.Verbose {
					log.Println(err)
				}

				responseObject = submitResponse{}
				responseObject.Messages = []string{string(responseBody)}
			}

			log.Printf("App Catalog returned statuscode %v. Response details:\n", response.StatusCode)
			for _, line := range responseObject.Messages {
				log.Printf("\t%v\n", line)
			}

			if cmd.Verbose {
				for key, val := range responseObject.Links {
					log.Printf("\tLINK: %s\t\t%s", key, val)
				}
			}

			if response.StatusCode == http.StatusOK {
				log.Println("App has been submitted successfully.")

				if acceptanceQueryUrlPath, ok := responseObject.Links["acc:query"]; ok {
					log.Println("You can use the following query URL to get this particular version of this app:")
					log.Printf("\t%s%s\n", rootURI, acceptanceQueryUrlPath)
				}
			} else {
				return fmt.Errorf("Submit failed, App Catalog returned statuscode %v", response.StatusCode)
			}

			return nil
		}).
		Alias("pub").
		Alias("publ")

	command.Arg("appPath", "path to the App folder (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
}
