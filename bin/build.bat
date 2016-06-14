@echo off
pushd %~dp0%
cd ..

echo "Retrieving/updating vendor packages using GVT..."

gvt update gopkg.in/alecthomas/kingpin.v2 || true
gvt update github.com/nu7hatch/gouuid || true

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

rem Since this is a .bat file, it's safe to assume this is run in Windows

echo "Done!"
popd