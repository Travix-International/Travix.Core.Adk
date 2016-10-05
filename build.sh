#!/bin/bash

set -euf -o pipefail

BUILD_DATE=`LC_ALL=en_US.utf8 date -u +"%a.%B.%d.%Y.%R:%S.%z.%Z"`
: "${TRAVIS_TAG:=0.0.0}"
: "${TRAVIS_COMMIT:=`git rev-parse --short HEAD`}"
APP_LDFLAGS="-s -X main.version=$TRAVIS_TAG -X main.gitHash=$TRAVIS_COMMIT -X main.buildDate=$BUILD_DATE"

echo "TRAVIS_BUILD_ID: $TRAVIS_BUILD_ID"
echo "TRAVIS_BUILD_NUMBER: $TRAVIS_BUILD_NUMBER"
echo "TRAVIS_JOB_NUMBER: $TRAVIS_JOB_NUMBER"
echo "TRAVIS_COMMIT: $TRAVIS_COMMIT"

echo "Building Windows binary..."
GOARCH=amd64 GOOS=windows go build -ldflags "$APP_LDFLAGS" -o bin/appix.exe -i .

echo "Building Mac binary..."
GOARCH=amd64 GOOS=darwin go build -ldflags "$APP_LDFLAGS" -o bin/appix-mac -i .

echo "Building Linux binary..."
GOARCH=amd64 GOOS=linux go build -ldflags "$APP_LDFLAGS" -o bin/appix-linux -i . 

ls -la bin

echo "Done!"
