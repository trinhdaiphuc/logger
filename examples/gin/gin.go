package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(logger.GinMiddleware(logger.ConfigGin{SkipperGin: func(context *gin.Context) bool {
		if context.Request.RequestURI == "/metrics" {
			return true
		}
		return false
	}}))
	server.GET("/hello/:name", func(ctx *gin.Context) {
		log := logger.GetLogger(ctx)
		name := ctx.Param("name")
		log.AddLog("request name %v", name)
		ctx.String(200, "Hello "+name)
	})

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
