package server

import (
	"github.com/graphql-go/graphql"
)

// GraphQL OUTPUT types ****************************************************

var dominantColorObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "DominantColor",
	Description: "Dominant color in the image",
	Fields: graphql.Fields{
		"hex": &graphql.Field{
			Type:        graphql.String,
			Description: "Hex representation of the color",
		},
		"is_light": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "If the color looks light",
		},
	},
})

var imageObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ImageFile",
	Description: "Can be a file of any type. Not only image.",
	Fields: graphql.Fields{
		"filepath": &graphql.Field{
			Type:        graphql.String,
			Description: "File URI",
		},
		"ext": &graphql.Field{
			Type:        graphql.String,
			Description: "Filename extension",
		},
		"initial_width": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width in pixels of the original image (for images)",
		},
		"initial_height": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height in pixels of the original image (for images)",
		},
		"initial_size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Initial file size in bytes",
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
		"dominant_color": &graphql.Field{
			Type:        dominantColorObject,
			Description: "Dominant color of the image",
		},
	},
})

// GraphQL INPUT types ****************************************************

var inputCropRectObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "InputRect",
	Description: "Rectangular area on the image",
	Fields: graphql.InputObjectConfigFieldMap{
		"x": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "distance between LEFT sides of the image and the crop rectangle in pixels",
		},
		"y": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "distance between TOP sides of the image and the crop rectangle in pixels",
		},
		"width": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "width of the rectangle in pixels",
		},
		"height": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "height of the rectangle in pixels",
		},
	},
})
