package server

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// HeadersMiddleware добавляет HTTP заголовки к ответу сервера
func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// c.Header("Access-Control-Allow-Origin", "*")
		// if hostIsAllowed(c.Request.Host) {
		// 	c.Header("Access-Control-Allow-Origin", "*")
		// }

		origin := c.GetHeader("Origin")
		// fmt.Println("Origin:", origin, "c.Request.Host:", c.Request.Host)
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Если конечное приложение  не установило Access-Control-Allow-Credentials добавляем его
		if c.GetHeader("Access-Control-Allow-Credentials") == "" {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		// c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Access-Control-Max-Age", "600")
		c.Next()
	}
}

// hostIsAllowed - проверяем пришел ли запрос с разрешенного хоста
func hostIsAllowed(host string) bool {
	if strings.HasPrefix(host, "localhost") ||
		strings.HasPrefix(host, "127.0.0.1") ||
		strings.Contains(host, ".rg.ru:") ||
		strings.HasSuffix(host, ".rg.ru") {
		return true
	}
	return false
}
