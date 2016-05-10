#!/usr/bin/env bash

_appixsetup_has() {
    type "$1" > /dev/null 2>&1
    return $?
}

_appixsetup_update_profile() {
    local profile="$1"
    local sourceString="$2"
    if ! grep -qc '.appix' $profile; then
        echo "Adding folder to PATH in $profile"
        echo "" >> "$profile"
        echo $sourceString >> "$profile"
    else
        echo "=> Folder is already added to PATH in $profile"
    fi
}

echo "Starting the Appix ADK installation"

if ! _appixsetup_has "curl"; then
    echo "appixsetup requires curl to be installed"
    exit 1
fi

if [ -z "$APPIX_BIN_URL" ]; then
    APPIX_BIN_URL="https://raw.githubusercontent.com/markvincze/jil-playground/master/appix-linux"
fi

# Downloading to ~/.appix
mkdir -p ~/.appix
if [ -s "~/.appix/appix" ]; then
    echo "appix is already installed in ~/.appix, trying to update"
else
    echo "Downloading appix binary to ~/.appix"
fi

curl -s "$APPIX_BIN_URL" -o ~/.appix/appix || {
    echo >&2 "Failed to download '$APPIX_BIN_URL'."
    exit 1
}

# Adding execute permission
chmod +x ~/.appix/appix

echo

# Detect profile file if not specified as environment variable (eg: PROFILE=~/.myprofile).
if [ -z "$PROFILE" ]; then
    if [ -f "$HOME/.bash_profile" ]; then
        PROFILE="$HOME/.bash_profile"
    elif [ -f "$HOME/.bashrc" ]; then
        PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.profile" ]; then
        PROFILE="$HOME/.profile"
    fi
fi

if [ -z "$ZPROFILE" ]; then
    if [ -f "$HOME/.zshrc" ]; then
        ZPROFILE="$HOME/.zshrc"
    fi
fi

ADD_TO_PATH_STR=$'if [ -d \"$HOME/.appix\" ]; then export PATH=\"$HOME/.appix:$PATH\"; fi # Add Appix binary folder to the path'

if [ -z "$PROFILE" -a -z "$ZPROFILE" ] || [ ! -f "$PROFILE" -a ! -f "$ZPROFILE" ] ; then
    if [ -z "$PROFILE" ]; then
      echo "Profile not found. Tried ~/.bash_profile ~/.zshrc and ~/.profile."
      echo "Create one of them and run this script again"
    elif [ ! -f "$PROFILE" ]; then
      echo "Profile $PROFILE not found"
      echo "Create it (touch $PROFILE) and run this script again"
    else
      echo "Profile $ZPROFILE not found"
      echo "Create it (touch $ZPROFILE) and run this script again"
    fi
    echo "  OR"
    echo "Append the following line to the correct file yourself:"
    echo
    echo " $ADD_TO_PATH_STR"
    echo
else
    [ -n "$PROFILE" ] && _appixsetup_update_profile "$PROFILE" "$ADD_TO_PATH_STR"
    [ -n "$ZPROFILE" ] && _appixsetup_update_profile "$ZPROFILE" "$ADD_TO_PATH_STR"
fi

echo "The Appix ADK has been installed. You can start using it by typing appix."
