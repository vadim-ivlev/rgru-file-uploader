package server

import (
	"github.com/graphql-go/graphql"
)

// TYPES ****************************************************

var imageType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "File",
	Description: "Файл. width и height поля имеют смысл только для файлов изображений",
	Fields: graphql.Fields{
		"filepath": &graphql.Field{
			Type:        graphql.String,
			Description: "URI изображения",
		},
		"width": &graphql.Field{
			Type:        graphql.Int,
			Description: "Ширина в пикселях",
		},
		"height": &graphql.Field{
			Type:        graphql.Int,
			Description: "Высота в пикселях",
		},
		"size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Размер оптимизированного файла в байтах",
		},
		"initial_size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Размер файла в байтах",
		},
	},
})
