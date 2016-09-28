Write-Output "Building the appix binaries..."

if (Test-Path env:APPVEYOR_BUILD_VERSION) {
    $BUILD_VERSION = $APPVEYOR_BUILD_VERSION
    $GIT_HASH = $APPVEYOR_REPO_COMMIT
} else {
    $BUILD_VERSION = "0.0.0"
    $GIT_HASH = "000"
}
$BUILD_DATE = Get-Date -UFormat "%a.%B.%d.%Y.%R:%S.%z.%Z"
$APP_LDFLAGS="-s -X main.version=$BUILD_VERSION -X main.gitHash=$GIT_HASH -X main.buildDate=$BUILD_DATE"


$env:GOARCH="amd64"
$env:GOOS="linux"
Write-Output "Building Linux binary..."
go build -ldflags "$APP_LDFLAGS" -o bin\appix-linux -i .

$env:GOOS="darwin"
Write-Output "Building Mac binary..."
go build -ldflags "$APP_LDFLAGS" -o bin\appix-mac -i .

$env:GOOS="windows"
Write-Output "Building Windows binary..."
go build -ldflags "$APP_LDFLAGS" -o bin\appix.exe -i .

Write-Output "Done!"
