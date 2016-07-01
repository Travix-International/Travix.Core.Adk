#!/bin/bash

set -euf -o pipefail

echo "Retrieving/updating vendor packages using GVT..."

gvt update gopkg.in/alecthomas/kingpin.v2 || true
gvt update github.com/nu7hatch/gouuid || true

set GOARCH=amd64
set GOOS=darwin
echo "Building Mac binary..."
go build -o bin/appix-mac -i .

exit

set GOARCH=amd64
set GOOS=windows
echo "Building Windows binary..."
go build -o bin/appix.exe -i .

set GOOS=darwin
echo "Building Mac binary..."
go build -o bin/appix-mac -i .

set GOOS=linux
echo "Building Linux binary..."
go build -o bin/appix-linux -i . 

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
