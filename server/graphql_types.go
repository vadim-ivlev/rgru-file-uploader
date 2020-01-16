package server

import (
	"github.com/graphql-go/graphql"
)

// TYPES ****************************************************

var dominantColorType = graphql.NewObject(graphql.ObjectConfig{
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

var imageType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "File",
	Description: "File.",
	Fields: graphql.Fields{
		"filepath": &graphql.Field{
			Type:        graphql.String,
			Description: "File URI",
		},
		"ext": &graphql.Field{
			Type:        graphql.String,
			Description: "Filename extension",
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
		"dominant_color": &graphql.Field{
			Type:        dominantColorType,
			Description: "Dominant color in the image",
		},
	},
})

var cropRectType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "CropRect",
	Description: "Rectangular area on the image",
	Fields: graphql.Fields{
		"x": &graphql.Field{
			Type:        graphql.Int,
			Description: "distance between LEFT sides of the image and the crop rectangle in pixels",
		},
		"y": &graphql.Field{
			Type:        graphql.Int,
			Description: "distance between TOP sides of the image and the crop rectangle in pixels",
		},
		"width": &graphql.Field{
			Type:        graphql.Int,
			Description: "width of the rectangle in pixels",
		},
		"height": &graphql.Field{
			Type:        graphql.Int,
			Description: "height of the rectangle in pixels",
		},
	},
})

type Rect struct {
	x      int
	y      int
	width  int
	height int
}
