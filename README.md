# Travix Fireball ADK 
App Developer Kit for the Travix Fireball infrastructure. The ADK consists of `appix`, a CLI tool to publish Apps to the App Catalog. 

## Installation
The ADK can be used on either Windows, Linux and Mac. There are separate installation scripts for Windows and Linux/Mac.

### Windows
To install the latest version on Windows, execute the following command in a PowerShell window.

```
iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/Travix-International/Travix.Core.Adk/master/appixinstall.ps1'))
```

To install a specific version of the ADK, set the `APPIX_VERSION` environment variable before executing the above script.

```
$env:APPIX_VERSION='1.0.0.5'
iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/Travix-International/Travix.Core.Adk/master/appixinstall.ps1'))
```

### Linux

To install the latest version, run the following command in a terminal.

```
curl -sSL https://raw.githubusercontent.com/Travix-International/Travix.Core.Adk/master/appixinstall.sh | sh
```

To install a specific version of the ADK, run the following command.

```
curl -sSL https://raw.githubusercontent.com/Travix-International/Travix.Core.Adk/master/appixinstall.sh | APPIX_VERSION=1.0.0.5 sh
```

## Getting started

After installation, the ADK can be run by typing `appix` in the terminal. Run `appix --help` to view the available commands and usage.
