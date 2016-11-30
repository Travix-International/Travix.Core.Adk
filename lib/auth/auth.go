package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	modelsConfig "github.com/Travix-International/Travix.Core.Adk/models/config"
)

func createAuthFileIfNotExists(c *modelsConfig.Config) error {
	var _, statErr = os.Stat(c.AuthFilePath)
	if os.IsNotExist(statErr) {
		fmt.Println("File does not exist, creating...")
		var _, createErr = os.Create(c.AuthFilePath)

		return createErr
	}

	return nil
}

// Auth structs
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

// Token structs
type TokenBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	IdToken      string `json:"id_token"`
	ProjectId    string `json:"project_id"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserId       string `json:"user_id"`
}

// Profile structs
type Profile struct {
	Email          string
	FirebaseUserId string
	Id             int
	IsEnabled      bool
	IsVerified     bool
	Name           string
	PublisherId    string
}

type ProfileBody struct {
	HasProfile bool
	Profile    Profile
}

func (p ProfileBody) String() string {
	if p.HasProfile {
		var buf bytes.Buffer

		buf.WriteString(fmt.Sprintln("Email: " + p.Profile.Email))
		buf.WriteString(fmt.Sprintln("Name: " + p.Profile.Name))

		if p.Profile.IsEnabled == true {
			buf.WriteString(fmt.Sprintln("Enabled: Yes"))
		} else {
			buf.WriteString(fmt.Sprintln("Enabled: No"))
		}

		if p.Profile.IsVerified == true {
			buf.WriteString(fmt.Sprintln("Verified: Yes"))
		} else {
			buf.WriteString(fmt.Sprintln("Verified: No"))
		}

		buf.WriteString(fmt.Sprintln("Publisher ID: " + p.Profile.PublisherId))
		return buf.String()
	}
	return "No profile found."
}

func GetAuth(c *modelsConfig.Config) (*Auth, error) {
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

func FetchRefreshedToken(config *modelsConfig.Config, refreshToken string) (TokenBody, error) {
	tokenClient := &http.Client{}
	var tokenReqPayload = []byte(`{"grant_type":"refresh_token","refresh_token": "` + refreshToken + `"}`)
	tokenReq, tokenReqErr := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+config.FirebaseApiKey, bytes.NewBuffer(tokenReqPayload))
	if tokenReqErr != nil {
		return TokenBody{}, tokenReqErr
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenRes, tokenResErr := tokenClient.Do(tokenReq)
	if tokenResErr != nil {
		return TokenBody{}, tokenResErr
	}

	tokenBody := TokenBody{}
	json.NewDecoder(tokenRes.Body).Decode(&tokenBody)

	return tokenBody, nil
}

func FetchDeveloperProfile(config *modelsConfig.Config, tokenBody TokenBody) (ProfileBody, error) {
	tokenType := tokenBody.TokenType
	tokenValue := tokenBody.IdToken

	profileClient := &http.Client{}
	profileReq, profileReqErr := http.NewRequest("GET", config.DeveloperProfileUrl+"/profile", nil)
	if profileReqErr != nil {
		return ProfileBody{}, profileReqErr
	}

	profileReq.Header.Set("Content-Type", "application/json")
	profileReq.Header.Set("Authorization", tokenType+" "+tokenValue)
	profileRes, profileResErr := profileClient.Do(profileReq)
	if profileResErr != nil {
		return ProfileBody{}, profileResErr
	}

	profileBody := ProfileBody{}
	json.NewDecoder(profileRes.Body).Decode(&profileBody)

	return profileBody, nil
}

func StartServer(config *modelsConfig.Config) (url string, ch <-chan struct{}) {
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

	url = "http://localhost:" + config.AuthServerPort

	ctx, cancel := context.WithCancel(context.Background())
	ch = ctx.Done()

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

		fmt.Fprint(w, html)
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

		fmt.Fprintf(w, "File written to disk at: %s", config.AuthFilePath)
	})

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Login successful! You can close your browser tab now.")
		fmt.Println("Closing server...")
		cancel()
	})

	go func() {
		http.ListenAndServe(":"+config.AuthServerPort, nil)
	}()

	return
}
