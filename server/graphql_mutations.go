package server

import (
	"errors"
	"rgru-file-uploader/pkg/img"

	"github.com/graphql-go/graphql"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"upload_file": &graphql.Field{
			Type:        imageType,
			Description: "Загрузить файл c локального компьютера",
			Args: graphql.FieldConfigArgument{
				"file_field_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Имя (name) поля формы для загрузки файла. Пример: <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// сохраняем файл загруженный с компьютера пользователя
				filePath, initialSize, err := img.SaveFirstFormFile(params, "file_field_name")
				if err != nil {
					return nil, err
				}
				// Оптимизируем его если это изображение
				size, width, height := img.OptimizeIfImage(filePath)

				return map[string]interface{}{
					"filepath":     filePath,
					"width":        width,
					"height":       height,
					"initial_size": initialSize,
					"size":         size,
				}, nil

			},
		},
		"download_file": &graphql.Field{
			Type:        imageType,
			Description: "Загрузить файл из интернет",
			Args: graphql.FieldConfigArgument{
				"file_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Имя под которым нужно сохранить файл.",
				},
				"file_url": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "URL откуда загрузить файл.",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
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

				return map[string]interface{}{
					"filepath":     filePath,
					"width":        width,
					"height":       height,
					"initial_size": initialSize,
					"size":         size,
				}, nil

			},
		},
	},
})
