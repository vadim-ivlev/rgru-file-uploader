#!/bin/bash

echo "building ..."
export GO111MODULE=on

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 

# для sqlite CGO_ENABLED=1
# env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .


echo "чистим ./deploy/"

# clean deploy/ directory
rm -rf deploy/rgru-file-uploader
rm -rf deploy/configs_example


echo "копируем файлы в ./deploy/"

cp rgru-file-uploader   deploy/rgru-file-uploader
cp -r configs           deploy/configs_example

echo "осторожно копируем файлы в ./deploy/configs/ "

mkdir -p deploy/configs
cp -f configs/img.yaml  deploy/configs/img.yaml
