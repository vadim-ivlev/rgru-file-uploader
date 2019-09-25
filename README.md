# RG.RU. File Uploader

## Микросервис загрузки файлов

Используется как часть приложений нуждающихся в сохранении файлов на сервере.

Загруженные файлы сохраняются в директории указанной в настроечном файле `configs/img.yaml`.
Файлы сохраняются в поддиректориях вида:

    YYYY/MM/DD/RANDOM_NUMBER/

Изображения ширина которых превышает 1440px (specified in img.yaml) пропорционально ужимаются до 1440px.

## GraphQL API

Конечные точки GraphQL 
- `/schema` 
- `/graphql`


Методы для загрузки файлов:
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

Загруженные на сервер файлы доступны по URI 
`uploads + filepath`. Где `filepath` - то, что вернуло приложение.


### Запуск для разработчика

    go run main.go -serve 7700 -env=dev

### Сборка контейнеров

    sh/build-frontend-container-all.sh

### Деплой

    sh/deploy.sh

