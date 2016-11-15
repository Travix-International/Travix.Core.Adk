package main

import (
	"os/user"
	"log"
	"path/filepath"
)

type Config struct {
	DirectoryPath string
	AuthFilePath string

	FirebaseApiKey string
	FirebaseAuthDomain string
	FirebaseDatabaseUrl string
	FirebaseStorageBucket string
	FirebaseMessagingSenderId string

	AuthServerPort string
}

func GetConfig() Config {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	directoryPath := filepath.Join(user.HomeDir, ".appix")

	c := Config{
		DirectoryPath: directoryPath,
		AuthFilePath: filepath.Join(directoryPath, "auth.json"),

		FirebaseApiKey: "",
		FirebaseAuthDomain: "",
		FirebaseDatabaseUrl: "",
		FirebaseStorageBucket: "",
		FirebaseMessagingSenderId: "",

		AuthServerPort: "7001",
	}

	return c
}
