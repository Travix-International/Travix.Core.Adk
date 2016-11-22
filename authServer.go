package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func createAuthFileIfNotExists(c Config) error {
	var _, statErr = os.Stat(c.AuthFilePath)
	if os.IsNotExist(statErr) {
		fmt.Println("File does not exist, creating...")
		var _, createErr = os.Create(c.AuthFilePath)

		return createErr
	}

	return nil
}

type StsTokenManager struct {
	ApiKey         string
	RefreshToken   string
	AccessToken    string
	ExpirationTime int
}

type AuthUser struct {
	Uid             string
	DisplayName     string
	Email           string
	EmailVerified   bool
	ApiKey          string
	AppName         string
	AuthDomain      string
	StsTokenManager StsTokenManager
}

type AuthCredential struct {
	IdToken     string
	AccessToken string
	Provider    string
}

type Auth struct {
	User       AuthUser
	Credential AuthCredential
}

func GetAuth(c Config) (*Auth, error) {
	content, readErr := ioutil.ReadFile(c.AuthFilePath)
	if readErr != nil {
		log.Printf("Failed to read auth file %s", c.AuthFilePath)
		return nil, readErr
	}

	auth := Auth{}
	unmarshalErr := json.Unmarshal(content, &auth)

	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal auth content")
		return nil, unmarshalErr
	}

	return &auth, nil
}

func startAuthServer(c chan bool, config Config) {
	firebaseConfig := `
		<script src="https://www.gstatic.com/firebasejs/3.6.0/firebase.js"></script>
		<script>
		// Initialize Firebase
		var config = {
			apiKey: "` + config.FirebaseApiKey + `",
			authDomain: "` + config.FirebaseAuthDomain + `",
			databaseURL: "` + config.FirebaseDatabaseUrl + `",
			storageBucket: "` + config.FirebaseStorageBucket + `",
			messagingSenderId: "` + config.FirebaseMessagingSenderId + `"
		};
		firebase.initializeApp(config);

		var provider = new firebase.auth.GoogleAuthProvider();
		provider.addScope('https://www.googleapis.com/auth/plus.login');
		provider.setCustomParameters({
			'login_hint': 'user@travix.com'
		});
		</script>`

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
			<head>
				` + firebaseConfig + `

				<script>
				function loginWithGoogle() {
					firebase.auth().signInWithRedirect(provider);
				}

				function postData(path, params) {
					var method = "post";

					var form = document.createElement("form");
					form.setAttribute("method", method);
					form.setAttribute("action", path);

					for (var key in params) {
						if (params.hasOwnProperty(key)) {
							var hiddenField = document.createElement("input");

							hiddenField.setAttribute("type", "hidden");
							hiddenField.setAttribute("name", key);
							hiddenField.setAttribute("value", params[key]);

							form.appendChild(hiddenField);
						}
					}

					document.body.appendChild(form);
					form.submit();
				}

				firebase.auth().getRedirectResult()
					.then(function (result) {
						if (result.credential) {
							postData('/save', {
								content: JSON.stringify(result, null, 2)
							});
						}
					})
					.catch(function (error) {
						console.log('error', error);

						// Handle Errors here.
						var errorCode = error.code;
						var errorMessage = error.message;

						// The email of the user's account used.
						var email = error.email;

						// The firebase.auth.AuthCredential type that was used.
						var credential = error.credential;

						// ...
					});
				</script>
			</head>

			<body>
				<a href="#" onClick="javascript: loginWithGoogle()">
					Login with Google
				</a>
			</body>
		</html>`

		io.WriteString(w, html)
	})

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		var createErr = createAuthFileIfNotExists(config)
		if createErr != nil {
			panic(createErr)
		}

		formContent := r.FormValue("content")
		content := []byte(formContent)
		writeErr := ioutil.WriteFile(config.AuthFilePath, content, 0644)
		if writeErr != nil {
			panic(writeErr)
		}

		io.WriteString(w, "File written at: "+config.AuthFilePath)
		c <- true
	})

	http.ListenAndServe(":"+config.AuthServerPort, nil)
}
