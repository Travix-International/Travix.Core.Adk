@echo off
pushd %~dp0%
cd ..

echo "Retrieving/updating vendor packages..."
git subtree pull --prefix vendor/gopkg.in/alecthomas/kingpin.v2 https://gopkg.in/alecthomas/kingpin.v2.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/template https://github.com/alecthomas/template.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/units https://github.com/alecthomas/units.git master --squash

set GOARCH=amd64
set GOOS=linux
echo "Building Linux binary..."
go build -o bin\appix-linux -i .

set GOOS=darwin
echo "Building Mac binary..."
go build -o bin\appix-mac -i .

set GOOS=windows
echo "Building Windows binary..."
go build -o bin\appix.exe -i .

echo "Done!"
popd