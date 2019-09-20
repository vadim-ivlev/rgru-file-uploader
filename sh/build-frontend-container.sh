#!/bin/bash


# гасим бд
docker-compose down

# удаляем файлы бд, и чистим загрузки
sudo rm -rf pgdata uploads/* uploads_temp/*


# build a docker image 
docker build -t rgru/rgru-file-uploader:latest -f Dockerfile-frontend . 

# push the docker image 
docker login
docker push rgru/rgru-file-uploader:latest


# копируем docker-compose-frontend.yml и 
cp docker-compose-frontend.yml ../file-uploader/docker-compose.yml
cp readme-frontend.md ../file-uploader/readme.md
