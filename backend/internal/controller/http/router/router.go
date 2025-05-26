package router

import (
	"github.com/gin-gonic/gin"
)

type RoutersRegistrator interface {
	RegisterRoutes(router *gin.RouterGroup)
}

func SetupRouter(controllers ...RoutersRegistrator) *gin.Engine {
	router := gin.Default()
	//router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/api/v1")
	for _, controller := range controllers {
		controller.RegisterRoutes(v1)
	}

	return router
}
