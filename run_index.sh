#! /bin/bash +x
cd ./index || exit
go build
./index ../data
