package server

import (
	"errors"

	gq "github.com/graphql-go/graphql"
)

var rootMutation = gq.NewObject(gq.ObjectConfig{
	Name: "Mutation",
	Fields: gq.Fields{
		"upload_file": &gq.Field{
			Type:        imageType,
			Description: "Загрузить файл на сервер",
			Args: gq.FieldConfigArgument{
				"file_field_name": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Имя (name) поля формы для загрузки файла. <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {

				fileFieldName, _ := params.Args["file_field_name"].(string)
				if fileFieldName == "" {
					return nil, errors.New("file_field_name cannot be empty")
				}

				path, width, height, size, errMsg := SaveUploadedImage(params, "file_field_name")
				if errMsg != "" {
					msg := "Could not save uploaded image: Resolve(): " + errMsg
					return nil, errors.New(msg)
				}

				return map[string]interface{}{
					"filepath": path,
					"width":    width,
					"height":   height,
					"size":     size,
				}, nil

			},
		},
	},
})
