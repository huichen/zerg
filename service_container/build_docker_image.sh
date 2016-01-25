#!/usr/bin/env bash

set -x

rm -rf docker
mkdir docker
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker/service
cp Dockerfile docker/

docker build -t unmerged/zerg -f docker/Dockerfile docker/
