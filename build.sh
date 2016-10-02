#!/bin/bash

set -euf -o pipefail

BUILD_DATE=`LC_ALL=en_US.utf8 date -u +"%a.%B.%d.%Y.%R:%S.%z.%Z"`
: "${APPVEYOR_BUILD_VERSION:=0.0.0}"
: "${APPVEYOR_REPO_COMMIT:=`git rev-parse --short HEAD`}"
APP_LDFLAGS="-s -X main.version=$APPVEYOR_BUILD_VERSION -X main.gitHash=$APPVEYOR_REPO_COMMIT -X main.buildDate=$BUILD_DATE"

set GOARCH=amd64
set GOOS=windows
echo "Building Windows binary..."
go build -ldflags "$APP_LDFLAGS" -o bin/appix.exe -i .

set GOOS=darwin
echo "Building Mac binary..."
go build -ldflags "$APP_LDFLAGS" -o bin/appix-mac -i .

set GOOS=linux
echo "Building Linux binary..."
go build -ldflags "$APP_LDFLAGS" -o bin/appix-linux -i . 

echo "Ensuring correct binary is available as 'appix'..."

if [[ "$OSTYPE" == "darwin"* ]]; then
    # Mac OSX
    rm bin/appix && cp bin/appix-mac bin/appix
elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys"|| "$OSTYPE" == "win32" ]]; then
    # Do nothing (Windows .exe is already named appropriately)
    echo "bin/appix.exe"
else
    rm "bin/appix" && cp bin/appix-linux "bin/appix"
fi

echo "Done!"
