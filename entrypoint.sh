#!/bin/bash

set -euf -o pipefail

BUILD_DATE=`LC_ALL=en_US.utf8 date -u +"%a.%B.%d.%Y.%R:%S.%z.%Z"`
: "${TRAVIS_TAG:=0.0.0}"
: "${TRAVIS_COMMIT:=`git rev-parse --short HEAD`}"
APP_LDFLAGS="-s
-X main.version=$TRAVIS_TAG
-X main.gitHash=$TRAVIS_COMMIT
-X main.buildDate=$BUILD_DATE
-X main.travixFirebaseApiKey=$TRAVIX_FIREBASE_API_KEY
-X main.travixFirebaseAuthDomain=$TRAVIX_FIREBASE_AUTH_DOMAIN
-X main.travixFirebaseDatabaseUrl=$TRAVIX_FIREBASE_DATABASE_URL
-X main.travixFirebaseStorageBucket=$TRAVIX_FIREBASE_STORAGE_BUCKET
-X main.travixFirebaseMessagingSenderId=$TRAVIX_FIREBASE_MESSAGING_SENDER_ID
-X main.travixFirebaseRefreshTokenUrl=$TRAVIX_FIREBASE_REFRESH_TOKEN_URL
-X main.travixDeveloperProfileUrl=$TRAVIX_DEVELOPER_PROFILE_URL
-X main.travixLoggerUrl=$TRAVIX_LOGGER_URL
-X main.uploadUrl=$TRAVIX_UPLOAD_URL
-X github.com/Travix-International/appix/livereload.certContent=$TRAVIX_CERT_CONTENT
-X github.com/Travix-International/appix/livereload.keyContent=$TRAVIX_KEY_CONTENT"

echo "Building Linux binary..."
GOARCH=amd64 GOOS=linux go build -ldflags "$APP_LDFLAGS" -o $HOME/.appix/appix -i ./cmd/appix

echo "Done!"

echo "Setting credentials..."
cp ./auth.json $HOME/.appix/auth.json

echo "Done!"
