package server

import (
	"net/http"
	"rgru-file-uploader/pkg/img"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// defineRoutes -  Сопоставляет маршруты функциям контроллера
func defineRoutes(r *gin.Engine) {
	r.Handle("OPTIONS", "/graphql", PingHandler)
	r.Handle("POST", "/graphql", GraphQL)
	r.Handle("POST", "/schema", GraphQL)
	r.Handle("GET", "/metrics", gin.WrapH(promhttp.Handler()))
	r.Static("/uploads", "./"+img.Params.Localdir)
}

// PingHandler нужен для фронта, так как сначала отправляется метод с OPTIONS
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

// Setup определяет пути и присоединяет функции middleware.
func Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default() // output to console
	r.Use(CountersMiddleware())
	r.Use(HeadersMiddleware())
	defineRoutes(r)
	return r
}

// Serve запускает сервер на заданном порту.
func Serve(port string) {
	r := Setup()
	_ = r.Run(port)
}
