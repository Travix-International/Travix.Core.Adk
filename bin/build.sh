#!/bin/bash
# Note: This script hasn't been tested yet and probably needs some tweaking

set -euf -o pipefail

git subtree pull --prefix vendor/gopkg.in/alecthomas/kingpin.v2 https://gopkg.in/alecthomas/kingpin.v2.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/template https://github.com/alecthomas/template.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/units https://github.com/alecthomas/units.git master --squash

go build -o appix -i . 
