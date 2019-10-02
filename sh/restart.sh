#!/bin/bash

# гасим бд
# docker-compose down

# удаляем файлы бд, и чистим загрузки
sudo rm -rf uploads/* 

# поднимаем бд
# docker-compose up -d
# sleep 1

# запускаем приложение
go run main.go -serve 5500 -env=dev


