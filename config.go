package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// ConfigEcho defines a function which is executed just before the middleware.
type ConfigEcho struct {
	// SkipperEcho defines a function to skip middleware.
	SkipperEcho SkipperEcho

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFuncEcho BeforeFuncEcho
}

type (
	// SkipperEcho defines a function to skip middleware. Returning true skips processing
	// the middleware.
	SkipperEcho func(echo.Context) bool

	// BeforeFuncEcho defines a function which is executed just before the middleware.
	BeforeFuncEcho func(echo.Context)
)

// DefaultSkipperEcho returns false which processes the middleware.
func DefaultSkipperEcho(echo.Context) bool {
	return false
}

// ConfigGin defines a function which is executed just before the middleware.
type ConfigGin struct {
	// SkipperGin defines a function to skip middleware.
	SkipperGin SkipperGin

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFuncGin BeforeFuncGin
}

type (
	// SkipperGin defines a function to skip middleware. Returning true skips processing
	// the middleware.
	SkipperGin func(*gin.Context) bool

	// BeforeFuncGin defines a function which is executed just before the middleware.
	BeforeFuncGin func(*gin.Context)
)

// DefaultSkipperGin returns false which processes the middleware.
func DefaultSkipperGin(*gin.Context) bool {
	return false
}

// ConfigFiber defines a function which is executed just before the middleware.
type ConfigFiber struct {
	// SkipperFiber defines a function to skip middleware.
	SkipperFiber SkipperFiber

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFuncFiber BeforeFuncFiber
}

type (
	// SkipperFiber defines a function to skip middleware. Returning true skips processing
	// the middleware.
	SkipperFiber func(ctx *fiber.Ctx) bool

	// BeforeFuncFiber defines a function which is executed just before the middleware.
	BeforeFuncFiber func(ctx *fiber.Ctx)
)

// DefaultSkipperFiber returns false which processes the middleware.
func DefaultSkipperFiber(ctx *fiber.Ctx) bool {
	return false
}

// ConfigGrpc defines a function which is executed just before the middleware.
type ConfigGrpc struct {
	// SkipperGrpc defines a function to skip middleware.
	SkipperGrpc SkipperGrpc

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFuncGrpc BeforeFuncGrpc
}

type (
	// SkipperGrpc defines a function to skip middleware. Returning true skips processing
	// the middleware.
	SkipperGrpc func(context.Context, *grpc.UnaryServerInfo) bool

	// BeforeFuncGrpc defines a function which is executed just before the middleware.
	BeforeFuncGrpc func(context.Context, *grpc.UnaryServerInfo)
)

// DefaultSkipperGrpc returns false which processes the middleware.
func DefaultSkipperGrpc(context.Context, *grpc.UnaryServerInfo) bool {
	return false
}
