#!/bin/bash

set -e

echo "*************** Running migrations against authservice database **************"
# Apply migrations to the database
go run cmd/dbmigrate/main.go

# Apply migrations to the  test database and pass testdb name
echo "*************** Running migrations against authservicetest database **************"

go run cmd/dbmigrate/main.go -dbname=authservicetest



# Turn off go modules so that CompileDaemon isn't isn't included in the go 
# mod file and is not included in production builds

GO111MODULE=off go get github.com/githubnemo/CompileDaemon

CompileDaemon --build="go build -o bin/authservice main.go" --command=./bin/authservice