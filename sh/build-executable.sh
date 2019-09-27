#!/bin/bash

echo "building executable"

export GO111MODULE=on

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 

# для sqlite CGO_ENABLED=1
# env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .
