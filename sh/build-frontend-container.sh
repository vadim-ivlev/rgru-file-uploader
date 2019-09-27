#!/bin/bash

echo 'building frontend container'

./build.sh
cd deploy

echo building an image

docker build -t vadimivlev/file-uploader:latest .

echo pushing the image 

docker login
docker push vadimivlev/file-uploader:latest
