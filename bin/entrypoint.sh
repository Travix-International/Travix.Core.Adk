#!/bin/bash
#
# Script executed before entering Docker image
#

set -e

echo "--------------------------------------------------------------------------------"
echo "Fireball App Development Kit"
echo ""
echo "Please wait while we're starting up the website..."
echo "--------------------------------------------------------------------------------"

cd /home/docker/rwd
npm start

echo "--------------------------------------------------------------------------------"
echo "How to use this Fireball App Dev Docker image:"
echo ""
echo " * Use 'yo fireball-ui-app' to scaffold a new App and/or a Widget."
echo " * Use 'appix' to verify/publish your App with the App Catalog."
echo " * The NH website is exposed at port 3001"
echo ""
echo "Now you can start fireballin'! ^-^"
echo "--------------------------------------------------------------------------------"

#chown -R docker /home/docker/apps
bash

#exec "$@"