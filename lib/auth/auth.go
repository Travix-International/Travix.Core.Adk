package auth

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

type AuthData struct {
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
