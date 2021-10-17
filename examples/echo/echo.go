package main

import (
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	server := echo.New()
	server.Use(logger.EchoMiddleware)

	server.GET("/hello/:name", func(ctx echo.Context) error {
		log := logger.GetLogger(ctx.Request().Context())
		name := ctx.Param("name")
		log.AddLog("request name %v", name)
		return ctx.String(200, "Hello "+name)
	})

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
