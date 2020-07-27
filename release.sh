#!/bin/bash

mkdir -p dist/linux
mkdir -p dist/darwin

export GIT_TAG=`git describe --tags --candidates=1 --dirty`
GOOS=linux GOARCH=amd64 go build -ldflags "-X cmd.Version=$GIT_TAG" -o dist/linux/ecr-image-retag
GOOS=darwin GOARCH=amd64 go build -ldflags "-X cmd.Version=$GIT_TAG" -o dist/darwin/ecr-image-retag
echo "Done"
