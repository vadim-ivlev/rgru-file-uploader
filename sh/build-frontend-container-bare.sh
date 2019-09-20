#!/bin/bash


# гасим бд
docker-compose down

# удаляем файлы бд, и чистим загрузки
sudo rm -rf pgdata uploads/* uploads_temp/*

# компилируем. линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .


# build a docker image 
docker build -t rgru/rgru-file-uploader:bare -f Dockerfile-frontend-bare . 

# push the docker image 
docker login
docker push rgru/rgru-file-uploader:bare


# копируем docker-compose-frontend.yml и 
cp docker-compose-frontend-bare.yml ../file-uploader/docker-compose.yml
cp readme-frontend.md ../file-uploader/readme.md
