@echo off
pushd %~dp0%
cd ..
git subtree pull --prefix vendor/gopkg.in/alecthomas/kingpin.v2 https://gopkg.in/alecthomas/kingpin.v2.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/template https://github.com/alecthomas/template.git master --squash
git subtree pull --prefix vendor/github.com/alecthomas/units https://github.com/alecthomas/units.git master --squash
go build -o bin\appix.exe -i .
popd