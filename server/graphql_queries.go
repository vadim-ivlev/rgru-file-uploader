package server

import (
	"github.com/graphql-go/graphql"
)

// ************************************************************************

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"ping": &graphql.Field{
			Type:        graphql.String,
			Description: "quick test",
			Args:        graphql.FieldConfigArgument{},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return "pong", nil
			},
		},
	},
})
