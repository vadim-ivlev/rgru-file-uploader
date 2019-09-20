package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	// "go/ast"
	"net/http"
	"strings"

	"rgru-file-uploader/pkg/img"

	"github.com/gin-gonic/gin"
	gq "github.com/graphql-go/graphql"
)

// FUNCTIONS *******************************************************

// jsonStringToMap преобразует строку JSON в map[string]interface{}
func jsonStringToMap(s string) map[string]interface{} {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m)
	return m
}

// getParamsFromBody извлекает параметры запроса из тела запроса
func getParamsFromBody(c *gin.Context) (map[string]interface{}, error) {
	r := c.Request
	mb := make(map[string]interface{})
	if r.ContentLength > 0 {
		errBodyDecode := json.NewDecoder(r.Body).Decode(&mb)
		return mb, errBodyDecode
	}
	return mb, errors.New("No body")
}

// getPayload3 извлекает "query", "variables", "operationName".
// Decoded body has precedence over POST over GET.
func getPayload3(c *gin.Context) (query string, variables map[string]interface{}) {

	// Проверяем на существование данных из Form Data
	query = c.PostForm("query")
	variables = jsonStringToMap(c.PostForm("variables"))

	// если есть тело запроса то берем из Request Payload (для Altair)
	params, errBody := getParamsFromBody(c)
	if errBody == nil {
		query, _ = params["query"].(string)
		variables, _ = params["variables"].(map[string]interface{})
	}

	return
}

// SaveUploadedFile - сохраняет первый присоединенный в поле fileFieldName файл во временную директорию,
// загружает его на сервер и удаляет его из временной директории.
// Возвращает путь файла на сервере, размер файла, сообщение об ошибке.
func SaveUploadedFile(params gq.ResolveParams, fileFieldName string) (finalPath string, size int64, errMsg string) {
	// сохраняем
	filePath, size, err := img.SaveFirstFormFile(params, fileFieldName)
	if err != nil {
		return "", 0, err.Error()
	}
	finalPath = img.TrimLocaldir(filePath)
	return finalPath, size, ""
}

// SaveUploadedImage - сохраняет первый присоединенный в поле fileFieldName файл во временную директорию,
// оптимизирует его размер и порождает иконки разных размеров.
// Загружает полученные файлы на сервер и удаляет их из временной директории.
// Возвращает путь файла на сервере, ширину и высоту изображения, JSON строку иконок,  сообщение об ошибке.
func SaveUploadedImage(params gq.ResolveParams, fileFieldName string) (
	serverPath string, width int, height int, thumbsJSONStr string, errMsg string) {

	// сохраняем изображение
	filePath, _, err := img.SaveFirstFormFile(params, fileFieldName)
	if err != nil {
		return "", 0, 0, "", "SaveUploadedImage(): " + err.Error()
	}

	// проверяем допустимо ли расширение
	if !img.Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		// удаляем файлы вместе с директорией
		dir := path.Dir(filePath)
		if err = os.RemoveAll(dir); err != nil {
			errMsg += "SaveUploadedImage().InvalidExt.RemoveAll(): " + err.Error() + " \n"
		}
		return "", 0, 0, "", "SaveUploadedImage(): Wrong file type. Should be:" + fmt.Sprintf("%v", img.Params.ValidImgExtensions)
	}

	serverPath = img.TrimLocaldir(filePath)

	// оптимизируем изображение
	filePath, width, height = img.OptimizeImage(filePath)

	// Генерируем иконки
	thumbsJSONStr, err = img.GenerateIcons(filePath)
	if err != nil {
		errMsg = "SaveUploadedImage(): " + err.Error()
	}

	return serverPath, width, height, thumbsJSONStr, errMsg
}

// G R A P H Q L ********************************************************************************

var schema, _ = gq.NewSchema(gq.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

// GraphQL исполняет GraphQL запрос
func GraphQL(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100*1024*1024)

	query, variables := getPayload3(c)

	result := gq.Do(gq.Params{
		Schema:         schema,
		RequestString:  query,
		Context:        context.WithValue(context.Background(), "ginContext", c),
		VariableValues: variables,
	})

	c.JSON(200, result)
}
