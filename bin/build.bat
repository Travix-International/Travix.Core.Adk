@echo off
pushd %~dp0%
go build -o appix.exe -i ../
popd