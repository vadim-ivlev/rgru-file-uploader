# RG.RU. File Uploader

## Микросервис загрузки файлов

Используется как часть приложений нуждающихся в сохранении файлов на сервере.

Загруженные файлы сохраняются в директории указанной в настроечном файле `configs/img.yaml`.
Файлы сохраняются в поддиректориях вида:

    YYYY/MM/DD/RANDOM_NUMBER/

Изображения ширина которых превышает 1440px (specified in img.yaml) пропорционально ужимаются до 1440px.

## GraphQL

Конечные точки GraphQL 
- `/schema` 
- `/graphql`


Методы загрузки файлов:
- `upload_local_file (file_field_name)` для загрузки локальных файлов
- `upload_internet_file (file_url)` для загрузки файлов из интернет

Оба метода возвращают 

    {
        filepath        путь к сохраненному файлу 
        width           ширина сохраненного изображения px
        height          высота сохраненного изображения px
        size            размер сохраненного файла bytes
        initial_size    размер загруженного файла bytes
    }

## REST

Загруженные файлы доступны по URI 
`uploads + filepath`. Где `filepath` - то, что вернуло приложение.


### Запуск для разработчика

    go run main.go -serve 5500 -env=dev

### Сборка контейнеров для фронтэнд разработчиков

    sh/build-frontend-container.sh

### Пуш и запуск деплоя на https://git.rgwork.ru

    sh/push.sh
    sh/deploy.sh

### Полное обновление программы

полное обновление программы состоит из следующих этапов

1. Сборка контейнеров для фронтэнда
2. Выгрузка изменений в репозиторий
3. Запуск пайплайна на деплой https://git.rgwork.ru

        git add git add -A .
        git commit -m "fix: description"

        sh/build-frontend-container.sh
        sh/push.sh
        sh/deploy.sh


