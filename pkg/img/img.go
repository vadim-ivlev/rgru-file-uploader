// Package img - функции для декодирования, масштабирования и сохранения изображений в файл.
package img

import (
	// "bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	colors "gopkg.in/go-playground/colors.v1"
	"gopkg.in/yaml.v2"
)

// Params - общие параметры хранимые в YAML
type connectionParams struct {
	Localdir           string          `yaml:"localdir"`
	ValidImgExtensions map[string]bool `yaml:"valid_img_extensions"`
	MaxImageWidth      int             `yaml:"max_image_width"`
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
// Возвращает размер файла в байтах.
func saveImageToFile(dst image.Image, filePath string) int64 {
	err := imaging.Save(dst, filePath)
	if err != nil {
		log.Println("saveImageToFile(): ", err.Error())
		return 0
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	// get the size
	size := fi.Size()
	return size
}

// UlidNum - возвращает случайную строку числа в диапазоне [min,max)
func UlidNum(min, max int) string {
	t := time.Now().UnixNano()
	rand.Seed(t)
	return strconv.Itoa(rand.Intn(max-min) + min)
}

// CreateNewDirectory создаем директорию для хранения файлов
func CreateNewDirectory() (path string, err error) {
	saveDir := Params.Localdir + "/" + time.Now().Format("2006/01/02") + "/" + UlidNum(10000, 99999) + "/"
	err = os.MkdirAll(saveDir, 0777)
	return saveDir, err
}

// SaveFirstFile - сохраняет первый загруженный файл поля fileFieldName во временную директорию
// и возвращает путь к сохраненному файлу
func SaveFirstFile(c *gin.Context, fileFieldName string) (string, int64, error) {

	file, header, err := c.Request.FormFile(fileFieldName)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintln("SaveFirstFile 1:", err))

	}
	filename := header.Filename

	// создаем директорию для хранения файлов
	saveDir, err := CreateNewDirectory()
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

// DownloadFile сохраняет файл из интернета в локальный файл, не загружая его в оперативную память.
// Эффективен, поскольку пишет на диск по мере получения данных из интернет.
// Измененный код из: https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(filepath string, url string) (size int64, err error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// Write the body to file
	size, err = io.Copy(out, resp.Body)
	return size, err
}

// OptimizeIfImage Оптимизирует изображение с разрешенными расширениями
func OptimizeIfImage(filePath string) (size int64, width int, height int) {
	// Если это изображение оптимизируем его
	if Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		size, width, height = OptimizeImage(filePath)
		return size, width, height
	}
	return 0, 0, 0
}

// OptimizeImage - уменьшает изображение если нужно,
// Возвращает путь к оптимизированному изображению, его ширину и высоту
func OptimizeImage(filePath string) (size int64, width int, height int) {
	img, err := imaging.Open(filePath)
	if err != nil {
		fmt.Printf("OptimizeImage: failed to open image: %v", err)
		return
	}
	resizedImg, width, height := resizeImage(img)
	size = saveImageToFile(resizedImg, filePath)
	return size, width, height
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

// TrimLocaldir - удаляет префикс временной директории загрузки из пути файла
func TrimLocaldir(path string) string {
	return strings.TrimPrefix(path, Params.Localdir)
}

// GetDominantColor returns dominant color of the image along with brightness flag.
// { hex : "FFFFFF", is_light: true}
// See: https://github.com/EdlinOrg/prominentcolor
func GetDominantColor(filePath string) map[string]interface{} {
	// Если это не изображение возвращаем nil
	if !Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		return nil
	}
	// Если не смогли открыть файл возвращаем nil
	img, err := imaging.Open(filePath)
	if err != nil {
		log.Printf("GetDominantColorIfImage: failed to open image: %v", err)
		return nil
	}

	imageColors, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if err != nil {
		return nil
	}
	dominantColorHexString := imageColors[0].AsString()
	color, err := colors.Parse("#" + dominantColorHexString)
	if err != nil {
		return nil
	}

	return map[string]interface{}{
		"hex":      dominantColorHexString,
		"is_light": color.IsLight(),
	}
}

// CropImage  Обрезает изображение filePath по размерам cropRect,
// и сохраняет его в croppedFilePath.
// Возвращает ширину, высоту в пикселях и размер в байтах.
func CropImage(filePath string, cropRect image.Rectangle, croppedFilePath string) (width, height int, size int64) {

	// Если это не изображение возвращаем 0
	if !Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		return
	}
	// Если не смогли открыть файл возвращаем 0
	img, err := imaging.Open(filePath)
	if err != nil {
		log.Printf("GetDominantColorIfImage: failed to open image: %v", err)
		return
	}

	croppedImage := imaging.Crop(img, cropRect)

	width = croppedImage.Rect.Dx()
	height = croppedImage.Rect.Dy()
	size = saveImageToFile(croppedImage, croppedFilePath)

	return
}

// getImageWidthHeightSize Возвращает ширину, высоту в пикселях изображения и размер в байтах.
func GetImageWidthHeight(filePath string) (width, height int16) {
	// Если это не изображение возвращаем 0
	if !Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		return
	}
	// Если не смогли открыть файл возвращаем 0
	// img, err := imaging.Open(filePath)
	// if err != nil {
	// 	log.Printf("GetDominantColorIfImage: failed to open image: %v", err)
	// 	return
	// }

	// width = img.
	// height = img.Bounds().Dy()
	return

	// // https://stackoverflow.com/questions/21741431/get-image-size-with-golang ------------------------------------
	// if reader, err := os.Open(filepath.Join(dir_to_scan, imgFile.Name())); err == nil {
	// 	defer reader.Close()
	// 	im, _, err := image.DecodeConfig(reader)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
	// 		continue
	// 	}
	// 	fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
	// } else {
	// 	fmt.Println("Impossible to open the file:", err)
	// }

}

// GetFileSize Возвращает размер файла в байтах.
func GetFileSize(filePath string) int64 {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Println(err)
		return 0
	}
	return fi.Size()
}
