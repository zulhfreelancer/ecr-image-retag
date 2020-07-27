#!/bin/bash

mkdir -p dist/linux
mkdir -p dist/darwin

export SHORT_VERSION=`cat VERSION`
export GIT_COMMIT_HASH=`git rev-parse --short HEAD`
LONG_VERSION="$SHORT_VERSION-$GIT_COMMIT_HASH"
echo $LONG_VERSION

GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/zulhfreelancer/ecr-image-retag/cmd.cliVersion=$LONG_VERSION" -o dist/linux/ecr-image-retag
GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/zulhfreelancer/ecr-image-retag/cmd.cliVersion=$LONG_VERSION" -o dist/darwin/ecr-image-retag
echo "Done"
