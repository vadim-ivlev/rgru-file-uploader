#!/bin/bash

echo "building ..."
export GO111MODULE=on

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 

# для sqlite
# env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .


echo "copying stuff ..."

# clean deploy/ directory
rm -rf deploy/rgru-file-uploader
rm -rf deploy/configs_example
# rm -rf deploy/migrations
# rm -rf deploy/templates
# rm -rf deploy/testapps

# careful with configs/
# rm -f deploy/configs/mail.yaml
# rm -f deploy/configs/sqlite.yaml
# rm -f deploy/configs/img.yaml



# copy files to deploy/
cp rgru-file-uploader   deploy/rgru-file-uploader
cp -r configs           deploy/configs_example
# cp -r migrations        deploy/migrations
# cp -r templates         deploy/templates
# cp -r testapps          deploy/testapps

#mv deploy/testapps/node_modules      deploy/testapps/nodemodules


# careful with configs/
mkdir -p deploy/configs
# cp -f configs/mail.yaml  deploy/configs/mail.yaml
# cp -f configs/sqlite.yaml  deploy/configs/sqlite.yaml
cp -f configs/img.yaml  deploy/configs/img.yaml
