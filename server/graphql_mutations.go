package server

import (
	"errors"
	"log"

	gq "github.com/graphql-go/graphql"
)

var rootMutation = gq.NewObject(gq.ObjectConfig{
	Name: "Mutation",
	Fields: gq.Fields{

		"create_image": &gq.Field{
			Type:        imageType,
			Description: "Создать медиа",
			Args: gq.FieldConfigArgument{
				// "id":           &gq.ArgumentConfig{Type: gq.NewNonNull(gq.Int), Description: "Идентификатор изображения"},
				"post_id": &gq.ArgumentConfig{
					Type:        gq.Int,
					Description: "Идентификатор поста",
				},
				"filepath": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "URI изображения",
				},
				"source": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Источник медиа",
				},
				"file_field_name": &gq.ArgumentConfig{
					Type:        gq.String,
					Description: "Имя (name) поля формы для загрузки файла. <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {

				_, ok := params.Args["file_field_name"].(string)
				if ok {
					path, width, height, thumbs, errMsg := SaveUploadedImage(params, "file_field_name")
					if errMsg == "" {
						params.Args["filepath"] = path
						params.Args["thumbs"] = thumbs
						params.Args["width"] = width
						params.Args["height"] = height
					} else {
						msg := "create_image: Resolve(): " + errMsg
						log.Println(msg)
						return nil, errors.New(msg)
					}
				}
				delete(params.Args, "file_field_name")
				return nil, nil

			},
		},
	},
})
