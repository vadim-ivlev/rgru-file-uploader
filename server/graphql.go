package server

import (
	"context"
	"encoding/json"
	"errors"

	// "fmt"
	// "os"
	// "path"
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

// SaveUploadedImage - сохраняет первый присоединенный в поле fileFieldName файл во временную директорию,
// оптимизирует его размер. Возвращает путь файла на сервере, ширину и высоту изображения.
func SaveUploadedImage(params gq.ResolveParams, fileFieldName string) (
	serverPath string, width int, height int, size int64, errMsg string) {

	// сохраняем изображение
	filePath, size, err := img.SaveFirstFormFile(params, fileFieldName)
	if err != nil {
		return "", 0, 0, 0, "SaveUploadedImage(): " + err.Error()
	}
	serverPath = img.TrimLocaldir(filePath)

	// проверяем расширение файла. Если это не изображение возвращаем как есть
	if !img.Params.ValidImgExtensions[strings.ToLower(filepath.Ext(filePath))] {
		return serverPath, 0, 0, size, errMsg
	}

	// иначе оптимизируем изображение
	_, width, height = img.OptimizeImage(filePath)

	return serverPath, width, height, size, errMsg
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
