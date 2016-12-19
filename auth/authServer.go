package auth

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Travix-International/appix/config"
)

func StartServer(config config.Config, done chan interface{}) {
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

				firebase.auth().getRedirectResult()
					.then(function (result) {
						if (!result.credential) {
							return;
						}

						var formData = new FormData();
						formData.append('content', JSON.stringify(result, null, 2));

						var request = new XMLHttpRequest();
						request.open('POST', '/save');
						request.onload = function () {
							window.location.href = '/success';
						};
						request.send(formData);
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
		var createErr = createAuthFileIfNotExists(config.AuthFilePath)
		if createErr != nil {
			panic(createErr)
		}

		formContent := r.FormValue("content")
		content := []byte(formContent)
		writeErr := ioutil.WriteFile(config.AuthFilePath, content, 0644)
		if writeErr != nil {
			panic(writeErr)
		}

		io.WriteString(w, "File written to disk at: "+config.AuthFilePath)
	})

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Login successful! You can close your browser tab now.")
		// closes the server
		done <- nil
	})

	http.ListenAndServe(":"+config.AuthServerPort, nil)
}

func createAuthFileIfNotExists(authFilePath string) error {
	var _, statErr = os.Stat(authFilePath)
	if os.IsNotExist(statErr) {
		fmt.Println("File does not exist, creating...")
		var _, createErr = os.Create(authFilePath)
		return createErr
	}

	return nil
}
