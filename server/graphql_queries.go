package server

import (
	gq "github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = gq.NewObject(gq.ObjectConfig{
	Name: "Query",
	Fields: gq.Fields{
		"get_image": &gq.Field{
			Type:        imageType,
			Description: "Показать изображение по идентификатору",
			Args: gq.FieldConfigArgument{
				"id": &gq.ArgumentConfig{
					Type:        gq.NewNonNull(gq.Int),
					Description: "Идентификатор изображения",
				},
			},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
	},
})
