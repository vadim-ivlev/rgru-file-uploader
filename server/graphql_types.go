package server

import (
	"encoding/json"

	gq "github.com/graphql-go/graphql"
)

// F U N C S ***********************************************

// JSONParamToMap - возвращает параметр paramName в map[string]interface{}.
// Второй параметр возврата - ошибка.
// Применяется для сериализации поля JSON таблицы postgres в map.
func JSONParamToMap(params gq.ResolveParams, paramName string) (interface{}, error) {

	source := params.Source.(map[string]interface{})
	param := source[paramName]

	// TODO: may be it's better to check if it can be converted to map[string]interface{}
	paramBytes, ok := param.([]byte)
	if !ok {
		return param, nil
	}
	var paramMap []map[string]interface{}
	err := json.Unmarshal(paramBytes, &paramMap)
	return paramMap, err
}

// FIELDS **************************************************

var imageFields = gq.Fields{
	"id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор медиа",
	},
	"post_id": &gq.Field{
		Type:        gq.Int,
		Description: "Идентификатор поста",
	},
	"filepath": &gq.Field{
		Type:        gq.String,
		Description: "URI изображения",
	},
	"thumbs": &gq.Field{
		Type:        gq.NewList(thumbType),
		Description: "Превью и изображение для видео - jsonb ",
		Resolve: func(params gq.ResolveParams) (interface{}, error) {
			return JSONParamToMap(params, "thumbs")
		},
	},
	"source": &gq.Field{
		Type:        gq.String,
		Description: "Источник медиа",
	},
	"width": &gq.Field{
		Type:        gq.Int,
		Description: "Ширина в пикселях",
	},
	"height": &gq.Field{
		Type:        gq.Int,
		Description: "Высота в пикселях",
	},
}

// TYPES ****************************************************

var imageType = gq.NewObject(gq.ObjectConfig{
	Name:        "Image",
	Description: "Медиа поста трансляции",
	Fields:      imageFields,
})

var thumbType = gq.NewObject(gq.ObjectConfig{
	Name:        "Thumb",
	Description: "Уменьшенное изображение для видео",
	Fields: gq.Fields{
		"type": &gq.Field{
			Type:        gq.String,
			Description: "Тип (small, middle, large)",
		},
		"filepath": &gq.Field{
			Type:        gq.String,
			Description: "Ссылка на файл на сервере",
		},
		"width": &gq.Field{
			Type:        gq.Int,
			Description: "Ширина изображения в пикселях.",
		},
		"height": &gq.Field{
			Type:        gq.Int,
			Description: "Высота изображения в пикселях.",
		},
	},
})
