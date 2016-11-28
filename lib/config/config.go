package config

import (
	"time"
)

type Config struct {
	Version         string
	BuildDate       string
	ParsedBuildDate time.Time
	GitHash         string
	CatalogURIs     map[string]string

	// DevFileName is the name of the file which contains the appix development settings for this specific application
	DevFileName string

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
