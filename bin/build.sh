#!/bin/bash

set -euf -o pipefail

echo "Retrieving/updating vendor packages..."
git subtree pull --prefix vendor/gopkg.in/alecthomas/kingpin.v2 https://gopkg.in/alecthomas/kingpin.v2.git master || true
git subtree pull --prefix vendor/github.com/alecthomas/template https://github.com/alecthomas/template.git master || true
git subtree pull --prefix vendor/github.com/alecthomas/units https://github.com/alecthomas/units.git master || true
git subtree pull --prefix vendor/github.com/nu7hatch/gouuid https://github.com/nu7hatch/gouuid.git master || true

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

rm bin/appix && cp bin/appix-linux bin/appix

echo "Done!"