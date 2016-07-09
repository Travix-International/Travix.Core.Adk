Write-Output "Retrieving/updating vendor packages using GVT..."

gvt update gopkg.in/alecthomas/kingpin.v2
gvt update github.com/nu7hatch/gouuid

$env:GOARCH=amd64
$env:GOOS=linux
Write-Output "Building Linux binary..."
go build -o bin\appix-linux -i .

$env:GOOS=darwin
Write-Output "Building Mac binary..."
go build -o bin\appix-mac -i .

$env:GOOS=windows
Write-Output "Building Windows binary..."
go build -o bin\appix.exe -i .

Write-Output "Done!"
