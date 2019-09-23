package server

import (
	gq "github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = gq.NewObject(gq.ObjectConfig{
	Name: "Query",
	Fields: gq.Fields{
		"ping": &gq.Field{
			Type:        gq.String,
			Description: "Тестовый метод",
			Args:        gq.FieldConfigArgument{},
			Resolve: func(params gq.ResolveParams) (interface{}, error) {
				return "pong", nil
			},
		},
	},
})
