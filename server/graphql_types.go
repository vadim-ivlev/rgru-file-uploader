package server

import (
	"github.com/graphql-go/graphql"
)

// TYPES ****************************************************

var imageType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "File",
	Description: "File.",
	Fields: graphql.Fields{
		"filepath": &graphql.Field{
			Type:        graphql.String,
			Description: "File URI",
		},
		"width": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width in pixels (for images)",
		},
		"height": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height in pixels (for images)",
		},
		"size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Size of the optimized image in bytes",
		},
		"initial_size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Initial file size in bytes",
		},
	},
})
