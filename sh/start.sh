#!/bin/bash

# поднимаем бд
# docker-compose up -d
# sleep 1

# запускаем приложение
go run main.go -serve 5500 -env=dev

