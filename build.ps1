function CheckLastExitCode {
    param ([String]$Platform = "")

    if ($LastExitCode -ne 0) {
        $msg = "The build failed for the platform $Platform"
        throw $msg
    }
}

Write-Output "Building the appix binaries..."

if (Test-Path env:APPVEYOR_BUILD_VERSION) {
    $BUILD_VERSION = $env:APPVEYOR_BUILD_VERSION
    $GIT_HASH = $env:APPVEYOR_REPO_COMMIT
} else {
    $BUILD_VERSION = "0.0.0"
    $GIT_HASH = "000"
}
$BUILD_DATE = (Get-Date -Date ((Get-Date).ToUniversalTime()) -UFormat %a.%B.%d.%Y.%R:%S) + ".+0000.UTC"
$APP_LDFLAGS="-s
-X main.version=$BUILD_VERSION
-X main.gitHash=$GIT_HASH
-X main.buildDate=$BUILD_DATE
-X main.travixFirebaseApiKey=$TRAVIX_FIREBASE_API_KEY
-X main.travixFirebaseAuthDomain=$TRAVIX_FIREBASE_AUTH_DOMAIN
-X main.travixFirebaseDatabaseUrl=$TRAVIX_FIREBASE_DATABASE_URL
-X main.travixFirebaseStorageBucket=$TRAVIX_FIREBASE_STORAGE_BUCKET
-X main.travixFirebaseMessagingSenderId=$TRAVIX_FIREBASE_MESSAGING_SENDER_ID
-X main.travixDeveloperProfileUrl=$TRAVIX_DEVELOPER_PROFILE_URL"
Write-Output "Load flags will be $APP_LDFLAGS"
Write-Output "AV vars are $env:APPVEYOR_BUILD_VERSION / $env:APPVEYOR_REPO_COMMIT"

$env:GOARCH="amd64"
$env:GOOS="linux"
Write-Output "Building Linux binary..."
go build -ldflags "$APP_LDFLAGS" -o bin\appix-linux -i ./cmd/appix

CheckLastExitCode "Linux"

$env:GOOS="darwin"
Write-Output "Building Mac binary..."
# NOTE: If we want to build on Windows and target OSX, we need the -tags kqueue option to make the notify library to compile. Otherwise it would give a build error.
# Details: https://github.com/rjeczalik/notify/issues/108#event-811951351
go build -tags kqueue -ldflags "$APP_LDFLAGS" -o bin\appix-mac -i ./cmd/appix

CheckLastExitCode "OSX"

$env:GOOS="windows"
Write-Output "Building Windows binary..."
go build -ldflags "$APP_LDFLAGS" -o bin\appix.exe -i ./cmd/appix

CheckLastExitCode "Windows"

Write-Output "Done!"
