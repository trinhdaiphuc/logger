package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.FiberMiddleware())

	app.Get("/hello/:name", func(ctx *fiber.Ctx) error {
		log := logger.GetLogger(ctx.Context())
		name := ctx.Params("name")
		log.AddLog("request name %v", name)
		return ctx.Status(200).SendString("Hello " + name)
	})

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
