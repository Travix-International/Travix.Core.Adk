package config

import (
	"time"
)

type Config struct {
	Version         string
	BuildDate       string
	ParsedBuildDate time.Time
	GitHash         string
	Verbose         bool

	DirectoryPath string
	AuthFilePath  string

	DeveloperProfileUrl string

	FirebaseApiKey            string
	FirebaseAuthDomain        string
	FirebaseDatabaseUrl       string
	FirebaseStorageBucket     string
	FirebaseMessagingSenderId string

	AuthServerPort string
}
