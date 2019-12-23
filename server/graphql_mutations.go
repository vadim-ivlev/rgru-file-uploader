package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"rgru-file-uploader/pkg/img"
	"rgru-file-uploader/pkg/signature"
	"rgru-file-uploader/pkg/vutils"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"upload_local_file": &graphql.Field{
			Type:        imageType,
			Description: "Upload a local file",
			Args: graphql.FieldConfigArgument{
				"file_field_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Input name for uploading files. Пример: <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Проверяем подпись
				err := checkSignature(params)
				if err != nil {
					return nil, fmt.Errorf("Signature check failed. %v", err)
				}

				// сохраняем файл загруженный с компьютера пользователя
				filePath, initialSize, err := img.SaveFirstFormFile(params, "file_field_name")
				if err != nil {
					return nil, err
				}
				// Оптимизируем его если это изображение
				size, width, height := img.OptimizeIfImage(filePath)

				// Уровень доступа, для возможности удаления файла.
				err = os.Chmod(filePath, 0777)
				if err != nil {
					log.Println(err)
				}

				return map[string]interface{}{
					"filepath":     img.TrimLocaldir(filePath),
					"width":        width,
					"height":       height,
					"initial_size": initialSize,
					"size":         size,
				}, nil

			},
		},
		"upload_internet_file": &graphql.Field{
			Type:        imageType,
			Description: "Upload file from Internet",
			Args: graphql.FieldConfigArgument{
				"file_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "File name for the uploaded file",
				},
				"file_url": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "URL of the file",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Проверяем подпись
				err := checkSignature(params)
				if err != nil {
					return nil, fmt.Errorf("Signature check failed. %v", err)
				}

				// Создаем директорию для хранения файла
				dirName, err := img.CreateNewDirectory()
				if err != nil {
					return nil, err
				}

				fileName := params.Args["file_name"].(string)
				if fileName == "" {
					return nil, errors.New("file_name can't be empty")
				}

				fileUrl := params.Args["file_url"].(string)
				if fileName == "" {
					return nil, errors.New("file_url can't be empty")
				}

				filePath := dirName + fileName

				// сохраняем файл загруженный из интернет
				initialSize, err := img.DownloadFile(filePath, fileUrl)
				if err != nil {
					return nil, err
				}
				// Оптимизируем его если это изображение
				size, width, height := img.OptimizeIfImage(filePath)

				// Уровень доступа, для возможности удаления файла.
				err = os.Chmod(filePath, 0777)
				if err != nil {
					log.Println(err)
				}

				return map[string]interface{}{
					"filepath":     img.TrimLocaldir(filePath),
					"width":        width,
					"height":       height,
					"initial_size": initialSize,
					"size":         size,
				}, nil

			},
		},
	},
})

// Проверяем подпись
func checkSignature(params graphql.ResolveParams) error {
	// Если ключ не предоставлен, значит проверять подпись не нужно
	if signature.PublicKeyText == "" {
		return nil
	}
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		return errors.New("SaveFirstFormFile(): Cannot get gin context.")
	}
	vutils.PrintRequestHeaders(c.Request)

	// Проверяем подпись
	return signature.Verify(c.Request)
}
