#!/bin/bash

set -euf -o pipefail

BUILD_DATE=`LC_ALL=en_US.utf8 date -u +"%a.%B.%d.%Y.%R:%S.%z.%Z"`
: "${APPVEYOR_BUILD_VERSION:=0.0.0}"
: "${APPVEYOR_REPO_COMMIT:=`git rev-parse --short HEAD`}"
APP_LDFLAGS="-s -X main.version=$APPVEYOR_BUILD_VERSION -X main.gitHash=$APPVEYOR_REPO_COMMIT -X main.buildDate=$BUILD_DATE"

echo "Building Windows binary..."
GOARCH=amd64 GOOS=windows go build -ldflags "$APP_LDFLAGS" -o bin/appix.exe -i .

echo "Building Mac binary..."
GOARCH=amd64 GOOS=darwin go build -ldflags "$APP_LDFLAGS" -o bin/appix-mac -i .

echo "Building Linux binary..."
GOARCH=amd64 GOOS=linux go build -ldflags "$APP_LDFLAGS" -o bin/appix-linux -i . 

ls -la bin

echo "Done!"
