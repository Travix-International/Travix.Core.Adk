# Travix Fireball ADK

App Developer Kit for the Travix Fireball infrastructure. The ADK consists of `appix`, a CLI tool to publish Apps to the App Catalog.

[![Build Status](https://travis-ci.org/Travix-International/appix.svg?branch=master)](https://travis-ci.org/Travix-International/appix)
[![GitHub release](https://img.shields.io/github/release/Travix-International/appix.svg)](https://github.com/Travix-International/appix/releases/latest)

## Installation

The ADK can be used on either Windows, Linux and Mac. There are separate installation scripts for Windows and Linux/Mac.

### Windows

To install the latest version on Windows, execute the following command in a PowerShell window.

```
iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/Travix-International/appix/master/appixinstall.ps1'))
```

To install a specific version of the ADK, set the `APPIX_VERSION` environment variable before executing the above script.

```
$env:APPIX_VERSION='1.0.0.5'
iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/Travix-International/appix/master/appixinstall.ps1'))
```

### Linux and macOS

To install the latest version, run the following command in a terminal.

```
$ curl -sSL https://raw.githubusercontent.com/Travix-International/appix/master/appixinstall.sh | sh
```

To install a specific version of the ADK, run the following command.

```
$ curl -sSL https://raw.githubusercontent.com/Travix-International/appix/master/appixinstall.sh | APPIX_VERSION=1.0.0.5 sh
```

## Getting started

After installation, the ADK can be run by typing `appix` in the terminal. Run `appix --help` to view the available commands and usage.

## Development

For developing the ADK itself:

### Clone

Clone the repo to your `$GOPATH`:

```
$ git clone git@github.com:Travix-International/appix.git $GOPATH/src/github.com/Travix-International/appix
```

### Dependencies

Install the dependencies with [gvt](https://github.com/FiloSottile/gvt):

```
$ cd $GOPATH/src/github.com/Travix-International/appix
$ gvt restore
```

### Environment variables

Make sure you have these environment variables (usually in your `~/.bash_profile` file) with the correct values:

```
export TRAVIX_FIREBASE_API_KEY=""
export TRAVIX_FIREBASE_AUTH_DOMAIN=""
export TRAVIX_FIREBASE_DATABASE_URL=""
export TRAVIX_FIREBASE_STORAGE_BUCKET=""
export TRAVIX_FIREBASE_MESSAGING_SENDER_ID=""
export TRAVIX_LOGGER_URL=""
export TRAVIX_FIREBASE_REFRESH_TOKEN_URL=""
export TRAVIX_CERT_CONTENT=""
export TRAVIX_KEY_CONTENT=""
export TRAVIX_DEVELOPER_PROFILE_URL=""
```

### Using a custom dev server

If we want to override the dev server to which `appix` is pushing an application under development, we have to add the `"DevServerOverride"` property to the `.appixDevSettings` file in the root folder of our application.

```
{
  ...
  "DevServerOverride": "https://my-dev-server.example.com"
}
```

### Build

Running this would generate `./bin/appix-mac` (if you are on macOS):

```
$ ./build.sh
```

Furthermore, you can copy it to the local appix path for easy access:

```
$ cp ./bin/appix-mac ~/.appix/appix
```

Now you can execute it as `$ appix` as usual.

## License

MIT @ [Travix International](http://travix.com)
