package router

import (
	"github.com/gin-gonic/gin"
)

type RoutersRegistrator interface {
	RegisterRoutes(router *gin.RouterGroup)
}

func SetupRouter(controllers []RoutersRegistrator, middlewares []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	v2 := router.Group("/api/v2")
	for _, controller := range controllers {
		controller.RegisterRoutes(v2)
	}

	return router
}
