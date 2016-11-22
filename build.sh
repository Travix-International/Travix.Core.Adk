#!/bin/bash

set -euf -o pipefail

BUILD_DATE=`LC_ALL=en_US.utf8 date -u +"%a.%B.%d.%Y.%R:%S.%z.%Z"`
: "${TRAVIS_TAG:=0.0.0}"
: "${TRAVIS_COMMIT:=`git rev-parse --short HEAD`}"
APP_LDFLAGS="-s -X main.version=$TRAVIS_TAG -X main.gitHash=$TRAVIS_COMMIT -X main.buildDate=$BUILD_DATE -X main.travixFirebaseApiKey=$TRAVIX_FIREBASE_API_KEY -X main.travixFirebaseAuthDomain=$TRAVIX_FIREBASE_AUTH_DOMAIN -X main.travixFirebaseDatabaseUrl=$TRAVIX_FIREBASE_DATABASE_URL -X main.travixFirebaseStorageBucket=$TRAVIX_FIREBASE_STORAGE_BUCKET -X main.travixFirebaseMessagingSenderId=$TRAVIX_FIREBASE_MESSAGING_SENDER_ID -X main.travixDeveloperProfileUrl=$TRAVIX_DEVELOPER_PROFILE_URL"

echo "Building Windows binary..."
GOARCH=amd64 GOOS=windows go build -ldflags "$APP_LDFLAGS" -o bin/appix.exe -i .

echo "Building Mac binary..."
GOARCH=amd64 GOOS=darwin go build -ldflags "$APP_LDFLAGS" -o bin/appix-mac -i .

echo "Building Linux binary..."
GOARCH=amd64 GOOS=linux go build -ldflags "$APP_LDFLAGS" -o bin/appix-linux -i .

echo "Done!"
