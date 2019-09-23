package server

import (
	gq "github.com/graphql-go/graphql"
)

// TYPES ****************************************************

var imageType = gq.NewObject(gq.ObjectConfig{
	Name:        "File",
	Description: "Файл. width и height поля имеют смысл только для файлов изображений",
	Fields: gq.Fields{
		"filepath": &gq.Field{
			Type:        gq.String,
			Description: "URI изображения",
		},
		"width": &gq.Field{
			Type:        gq.Int,
			Description: "Ширина в пикселях",
		},
		"height": &gq.Field{
			Type:        gq.Int,
			Description: "Высота в пикселях",
		},
		"size": &gq.Field{
			Type:        gq.Int,
			Description: "Размер файла в байтах",
		},
	},
})
