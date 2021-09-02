#!/bin/bash

set -e

# Turn off go modules so that CompileDaemon isn't isn't included in the go 
# mod file and is not included in production builds

GO111MODULE=off go get github.com/githubnemo/CompileDaemon

CompileDaemon --build="go build -o bin/authservice main.go" --command=./bin/authservice