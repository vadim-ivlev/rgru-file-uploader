# File Uploader


https://file-uploader.rg.ru


## Микросервис загрузки файлов, и обрезки загруженных изображений

<br><br>
<img src="images/uploader.png">
<br><br><br>

Загруженные файлы сохраняются в директории указанной в  `configs/img.yaml`.
Файлы сохраняются в поддиректориях вида:

    YYYY/MM/DD/RANDOM_NUMBER/

Изображения ширина которых превышает  `max_image_width` px (указанной в img.yaml) пропорционально ужимаются до ширины `max_image_width` px.

## Проверка цифровой подписи запросов на загрузку изображений

Для контроля откуда поступают запросы на загрузку изображений, программа проверяет цифровые подписи запросов. Проверка подписи  происходит 
<br>**если**
1. в файле `configs/signature.yaml` проставлены поля: 

    ```
    public_key_file: ./certificates/auth-proxy.key.pub
    keyid: auth-proxy
    ```
    и 

2. файл с публичным ключом присутствует на диске. 

**В противном случае** проверка подписи не производится.

Публичный ключ `auth-proxy` находится здесь 
<https://auth-proxy.rg.ru/publickey>


## GraphQL

Конечные точки GraphQL 
- `/schema` 


GraphQL Методы:

    upload_internet_file (file_name, file_url)    загрузить файл из интернет
    upload_local_file    (file_field_name)        загрузить файлов с компьютера
    crop_image           (file_path, crop_rect)   обрезать загруженный файл


**Замечание:** При вызове метода `upload_local_file()`  наряду со стандартными
GraphQL-полями  `query` и `variables` в HTTP запросе необходимо отправить бинарное поле с
содержимым файла. Имя этого поля передается в качестве
параметра в функцию upload_local_file(...), чтобы сервер мог знать откуда
брать содержимое файла.

Методы возвращают JSON структуры вида:

    {
        filepath        путь к сохраненному файлу               (string) 
        initial_width   ширина оригинального изображения px     (int)
        initial_height  высота оригинального изображения px     (int)
        initial_size    размер оригинального файла bytes        (int)
        width           ширина сохраненного изображения px      (int)
        height          высота сохраненного изображения px      (int)
        size            размер сохраненного файла bytes         (int)
        dominant_color {    
            hex         доминирующий цвет изображения           (string)
            is_light    светлое ли изображение                  (int)
        }
    }

## REST

Загруженные файлы доступны по URL 
https://image-loader.rg.ru/uploads + `filepath`. 
<br>Где `filepath` - то, что вернул метод загрузки/обрезки изображения.





<br><br><br><br><br><br>

---------------

### Запуск для разработчика

    go run main.go -serve 5500 -env=dev

или

    sh/start.sh

### Сборка контейнеров для фронтэнд разработчиков

    sh/build-frontend-container.sh

### Пуш и запуск деплоя на https://git.rgwork.ru

    sh/push.sh
    sh/deploy.sh



### Полное обновление программы
полное обновление программы состоит из следующих этапов

1. Сборка контейнеров для фронтэнда
2. Выгрузка изменений в репозиторий
3. Запуск деплоя https://git.rgwork.ru

```sh
git add git add -A .
git commit -m "fix: description"

sh/build-frontend-container.sh
sh/push.sh
sh/deploy.sh
```

