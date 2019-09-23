// Package img - функции для декодирования, масштабирования и сохранения изображений в файл.
package img

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"gopkg.in/yaml.v2"
)

// Thumb - описатель иконки изображения. Используем для resize, move etc.
type Thumb struct {
	Type     string `json:"type"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Filepath string `json:"filepath"`
}

// Params - общие параметры хранимые в YAML
type connectionParams struct {
	Localdir           string          `yaml:"localdir"`
	ValidImgExtensions map[string]bool `yaml:"valid_img_extensions"`
	MaxImageWidth      int             `yaml:"max_image_width"`
	ThumbsTemplate     []Thumb         `yaml:"thumbs_template"`
}

var Params connectionParams

// ReadConfig читает YAML
func ReadConfig(fileName string, env string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	envParams := make(map[string]connectionParams)
	err = yaml.Unmarshal(yamlFile, &envParams)
	if err != nil {
		log.Println(err.Error())
	}
	Params = envParams[env]
	return
}

// AppendToName добавляет строку к имени файла
func AppendToName(fileName string, str string) string {
	ext := filepath.Ext(fileName)
	thumbName := strings.TrimSuffix(fileName, ext) + str + ext
	return thumbName
}

// saveImageToFile - Сохраняет файл.
// Возвращает путь к сохраненному файлу.
func saveImageToFile(dst image.Image, filePath string) string {
	err := imaging.Save(dst, filePath)
	if err != nil {
		log.Println("saveImageToFile(): ", err.Error())
		return ""
	}
	return filePath
}

// thumbImage - генерирует уменьшенное изображение, заданное шириной и высотой.
// Возвращает масштабированное изображение, его ширину и высоту в пикселях.
func thumbImage(im image.Image, thumbWidth int, thumbHeight int) (dst image.Image, width int, height int) {
	b := im.Bounds()
	width = b.Dx()
	height = b.Dy()
	anchor := imaging.Top
	if width/height > thumbWidth/thumbHeight {
		anchor = imaging.Center
	}
	dst = imaging.Fill(im, thumbWidth, thumbHeight, anchor, imaging.Lanczos)
	b = dst.Bounds()
	width = b.Dx()
	height = b.Dy()
	return dst, width, height
}

// UlidNum - возвращает случайную строку числа в диапазоне [min,max)
func UlidNum(min, max int) string {
	t := time.Now().UnixNano()
	rand.Seed(t)
	return strconv.Itoa(rand.Intn(max-min) + min)
}

// SaveFirstFile - сохраняет первый загруженный файл поля fileFieldName во временную директорию
// и возвращает путь к сохраненному файлу
func SaveFirstFile(c *gin.Context, fileFieldName string) (string, int64, error) {

	file, header, err := c.Request.FormFile(fileFieldName)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintln("SaveFirstFile 1:", err))

	}
	filename := header.Filename

	// создаем имя директории для хранения файлов
	saveDir := Params.Localdir + "/file_uploader/" + time.Now().Format("2006/01/02") + "/" + UlidNum(10000, 99999) + "/"
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintln("SaveFirstFile 2:", err))
	}

	// filepath := saveDir + AppendToName(filename, "--"+GetULID())
	filepath := saveDir + filename
	out, err := os.Create(filepath)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintln("SaveFirstFile 3:", err))
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintln("SaveFirstFile 2:", err))
	}
	return filepath, header.Size, nil
}

// SaveFirstFormFile - сохраняет первый загруженный файл поля fileFieldName во временную директорию
// и возвращает путь к сохраненному файлу
func SaveFirstFormFile(p graphql.ResolveParams, fileFieldName string) (string, int64, error) {
	if fileFieldName == "" {
		return "", 0, errors.New("SaveFirstFormFile(): No file field name specified")
	}

	c, ok := p.Context.Value("ginContext").(*gin.Context)
	if !ok {
		return "", 0, errors.New("SaveFirstFormFile(): Cannot get gin context.")
	}

	fieldName, ok := p.Args[fileFieldName].(string)
	if !ok {
		return "", 0, errors.New("SaveFirstFormFile(): There is no '" + fileFieldName + "' field in the form.")
	}

	tempFile, size, err := SaveFirstFile(c, fieldName)
	if err != nil {
		return "", 0, err
	}

	return tempFile, size, nil
}

// OptimizeImage - уменьшает изображение если нужно,
// Возвращает путь к оптимизированному изображению, его ширину и высоту
func OptimizeImage(filePath string) (path string, width int, height int) {
	img, err := imaging.Open(filePath)
	if err != nil {
		fmt.Printf("OptimizeImage: failed to open image: %v", err)
		return
	}
	resizedImg, width, height := resizeImage(img)
	path = saveImageToFile(resizedImg, filePath)
	return path, width, height
}

// resizeImage масштабирует изображение если необходимо.
// Возвращает масштабированное изображение, его ширину и высоту в пикселях.
func resizeImage(im image.Image) (dst image.Image, width int, height int) {
	dst = im
	b := dst.Bounds()
	width = b.Dx()
	height = b.Dy()
	if width > Params.MaxImageWidth {
		dst = imaging.Resize(im, Params.MaxImageWidth, height*Params.MaxImageWidth/width, imaging.Lanczos)
	}
	b = dst.Bounds()
	width = b.Dx()
	height = b.Dy()

	return dst, width, height
}

// GenerateIcons - генерирует иконки заданного в filePath разных размеров,
// сохраняет их рядом с файлом и возвращает JSON строку иконок.
func GenerateIcons(filePath string) (string, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		fmt.Printf("GenerateIcons: failed to open image: %v", err)
		return "[]", err
	}

	thumbs := make([]Thumb, len(Params.ThumbsTemplate))
	copy(thumbs, Params.ThumbsTemplate)

	for i, thumb := range thumbs {
		dst, width, height := thumbImage(img, thumb.Width, thumb.Height)
		thumbFilePath := AppendToName(filePath, "--"+thumb.Type)
		thumbFilePath = saveImageToFile(dst, thumbFilePath)

		thumbs[i].Filepath = TrimLocaldir(thumbFilePath)
		thumbs[i].Width = width
		thumbs[i].Height = height
	}

	jsonBytes, err := json.Marshal(thumbs)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// TrimLocaldir - удаляет префикс временной директории загрузки из пути файла
func TrimLocaldir(path string) string {
	return strings.TrimPrefix(path, Params.Localdir)
}
