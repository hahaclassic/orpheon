package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	// Определяем абсолютные пути к YAML файлам
	projectRoot := "../../" // путь от cmd/gateway/main.go до корня проекта
	authOpenAPIPath := filepath.Join(projectRoot, "api/gateway/auth/v2/openapi.yaml")
	contentOpenAPIPath := filepath.Join(projectRoot, "api/gateway/content/v2/openapi.yaml")
	userOpenAPIPath := filepath.Join(projectRoot, "api/gateway/user/v2/openapi.yaml")
	streamingOpenAPIPath := filepath.Join(projectRoot, "api/gateway/streaming/v2/openapi.yaml")

	// Эндпоинты отдачи YAML
	r.GET("/api/gateway/auth/v2/openapi.yaml", func(c *gin.Context) {
		c.File(authOpenAPIPath)
	})
	r.GET("/api/gateway/content/v2/openapi.yaml", func(c *gin.Context) {
		c.File(contentOpenAPIPath)
	})
	r.GET("/api/gateway/user/v2/openapi.yaml", func(c *gin.Context) {
		c.File(userOpenAPIPath)
	})
	r.GET("/api/gateway/streaming/v2/openapi.yaml", func(c *gin.Context) {
		c.File(streamingOpenAPIPath)
	})

	// Swagger UI для каждого сервиса
	r.GET("/swagger/auth/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/api/gateway/auth/v2/openapi.yaml")))
	r.GET("/swagger/content/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/api/gateway/content/v2/openapi.yaml")))
	r.GET("/swagger/user/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/api/gateway/user/v2/openapi.yaml")))
	r.GET("/swagger/streaming/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/api/gateway/streaming/v2/openapi.yaml")))

	// Healthcheck
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
