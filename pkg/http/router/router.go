package router

import (
	"github.com/gin-gonic/gin"
)

type RoutersRegistrator interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// example: basePath = /api/v1
func SetupRouter(basePath string, controllers []RoutersRegistrator, middlewares []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	baseGroup := router.Group(basePath)
	for _, controller := range controllers {
		controller.RegisterRoutes(baseGroup)
	}

	return router
}
