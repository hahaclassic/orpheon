package swagger

// r := gin.Default()

// projectRoot := "../../"
// authOpenAPIPath := filepath.Join(projectRoot, "api/gateway-msv/auth/v2/openapi.yaml")
// contentOpenAPIPath := filepath.Join(projectRoot, "api/gateway-msv/content/v2/openapi.yaml")
// userOpenAPIPath := filepath.Join(projectRoot, "api/gateway-msv/user/v2/openapi.yaml")
// streamingOpenAPIPath := filepath.Join(projectRoot, "api/gateway-msv/streaming/v2/openapi.yaml")

// // Эндпоинты отдачи YAML
// r.GET("/api/gateway-msv/auth/v2/openapi.yaml", func(c *gin.Context) {
// 	c.File(authOpenAPIPath)
// })
// r.GET("/api/gateway-msv/content/v2/openapi.yaml", func(c *gin.Context) {
// 	c.File(contentOpenAPIPath)
// })
// r.GET("/api/gateway-msv/user/v2/openapi.yaml", func(c *gin.Context) {
// 	c.File(userOpenAPIPath)
// })
// r.GET("/api/gateway-msv/streaming/v2/openapi.yaml", func(c *gin.Context) {
// 	c.File(streamingOpenAPIPath)
// })

// // Swagger UI для каждого сервиса
// r.GET("/swagger/auth/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
// 	ginSwagger.URL("/api/gateway-msv/auth/v2/openapi.yaml")))
// r.GET("/swagger/content/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
// 	ginSwagger.URL("/api/gateway-msv/content/v2/openapi.yaml")))
// r.GET("/swagger/user/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
// 	ginSwagger.URL("/api/gateway-msv/user/v2/openapi.yaml")))
// r.GET("/swagger/streaming/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
// 	ginSwagger.URL("/api/gateway-msv/streaming/v2/openapi.yaml")))

// // Healthcheck
// r.GET("/status", func(c *gin.Context) {
// 	c.JSON(200, gin.H{"status": "ok"})
// })

// r.Run(":8080")
