package server

import (
	"errors"
	"image"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rgru-file-uploader/pkg/img"
	"rgru-file-uploader/pkg/signature"
	"rgru-file-uploader/pkg/vutils"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"upload_local_file": &graphql.Field{
			Type:        imageType,
			Description: "Upload a local file",
			Args: graphql.FieldConfigArgument{
				"file_field_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Input name for uploading files. Пример: <input name='fname' type='file' ...>",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Проверяем подпись
				err := checkSignature(params)
				if err != nil {
					return nil, fmt.Errorf("Signature check failed. %v", err)
				}

				// сохраняем файл загруженный с компьютера пользователя
				filePath, initialSize, err := img.SaveFirstFormFile(params, "file_field_name")
				if err != nil {
					return nil, err
				}
				// Оптимизируем его если это изображение
				size, width, height := img.OptimizeIfImage(filePath)

				// Устанавливаем уровень доступа, для возможности удаления файла другими процессами
				err = os.Chmod(filePath, 0777)
				if err != nil {
					log.Println(err)
				}

				return map[string]interface{}{
					"filepath":       img.TrimLocaldir(filePath),
					"ext":            filepath.Ext(filePath),
					"width":          width,
					"height":         height,
					"initial_size":   initialSize,
					"size":           size,
					"dominant_color": img.GetDominantColor(filePath),
				}, nil
			},
		},

		"upload_internet_file": &graphql.Field{
			Type:        imageType,
			Description: "Upload file from Internet",
			Args: graphql.FieldConfigArgument{
				"file_name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "File name for the uploaded file",
				},
				"file_url": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "URL of the file",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Проверяем подпись
				err := checkSignature(params)
				if err != nil {
					return nil, fmt.Errorf("Signature check failed. %v", err)
				}

				// Создаем директорию для хранения файла
				dirName, err := img.CreateNewDirectory()
				if err != nil {
					return nil, err
				}

				fileName := params.Args["file_name"].(string)
				if fileName == "" {
					return nil, errors.New("file_name can't be empty")
				}

				fileUrl := params.Args["file_url"].(string)
				if fileName == "" {
					return nil, errors.New("file_url can't be empty")
				}

				filePath := dirName + fileName

				// сохраняем файл загруженный из интернет
				initialSize, err := img.DownloadFile(filePath, fileUrl)
				if err != nil {
					return nil, err
				}
				// Оптимизируем его если это изображение
				size, width, height := img.OptimizeIfImage(filePath)

				// Устанавливаем уровень доступа, для возможности удаления файла другими процессами
				err = os.Chmod(filePath, 0777)
				if err != nil {
					log.Println(err)
				}

				return map[string]interface{}{
					"filepath":       img.TrimLocaldir(filePath),
					"ext":            filepath.Ext(filePath),
					"width":          width,
					"height":         height,
					"initial_size":   initialSize,
					"size":           size,
					"dominant_color": img.GetDominantColor(filePath),
				}, nil
			},
		},

		"crop_image": &graphql.Field{
			Type:        imageType,
			Description: "Crop image file",
			Args: graphql.FieldConfigArgument{
				"file_path": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "File name for the uploaded file",
				},
				"crop_rect": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(cropRectType),
					Description: "Rectangular area of the image",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Проверяем подпись
				err := checkSignature(params)
				if err != nil {
					return nil, fmt.Errorf("Signature check failed. %v", err)
				}

				// Создаем директорию для хранения файла
				dirName, err := img.CreateNewDirectory()
				if err != nil {
					return nil, err
				}

				filePath := img.Params.Localdir + "/" + params.Args["file_path"].(string)
				fileName := filepath.Base(filePath)
				croppedFilePath := dirName + fileName
				cropRect := params.Args["crop_rect"].(Rect)
				initialSize := getFileSize(filePath)

				// Обрезаем  изображение
				size := img.CropImage(filePath, image.Rect(cropRect.x, cropRect.y, cropRect.x + cropRect.width, cropRect.y+cropRect.height), croppedFilePath)

				// Устанавливаем уровень доступа, для возможности удаления файла другими процессами
				err = os.Chmod(croppedFilePath, 0777)
				if err != nil {
					log.Println(err)
				}

				return map[string]interface{}{
					"filepath":       img.TrimLocaldir(croppedFilePath),
					"ext":            filepath.Ext(croppedFilePath),
					"width":          cropRect.width,
					"height":         cropRect.height,
					"initial_size":   initialSize,
					"size":           size,
					"dominant_color": img.GetDominantColor(croppedFilePath),
				}, nil
			},
		},
	},
})

func getFileSize(filePath string) int64 {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Println(err)
		return 0
	}
	return fi.Size()
}

// Проверяем подпись
func checkSignature(params graphql.ResolveParams) error {
	// Если ключ не предоставлен, значит проверять подпись не нужно
	if signature.PublicKeyText == "" {
		return nil
	}
	c, ok := params.Context.Value("ginContext").(*gin.Context)
	if !ok {
		return errors.New("SaveFirstFormFile(): Cannot get gin context.")
	}
	// Проверяем подпись
	err := signature.Verify(c.Request)
	if err != nil {
		vutils.PrintRequestHeaders(c.Request)
	}
	return err
}
