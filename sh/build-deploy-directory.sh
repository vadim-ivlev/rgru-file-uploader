#!/bin/bash

echo 'building /deploy directory'

sh/build-executable.sh

echo "чистим ./deploy/"

rm -rf deploy/rgru-file-uploader
rm -rf deploy/configs_example
rm -rf deploy/certificates



echo "копируем файлы в ./deploy/"

mv rgru-file-uploader   deploy/rgru-file-uploader
cp -r configs           deploy/configs_example
cp -r certificates      deploy/certificates



echo "осторожно копируем файлы в ./deploy/configs/ "

mkdir -p deploy/configs
cp -f configs/img.yaml  deploy/configs/img.yaml
cp -f configs/signature.yaml  deploy/configs/signature.yaml
