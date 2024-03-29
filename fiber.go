package logger

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

var DefaultConfigFiber = ConfigFiber{
	SkipperFiber: DefaultSkipperFiber,
}

func FiberMiddleware(config ConfigFiber) fiber.Handler {
	if config.SkipperFiber == nil {
		config.SkipperFiber = DefaultSkipperFiber
	}

	return func(ctx *fiber.Ctx) error {
		logger := New(WithFormatter(&JSONFormatter{}))
		if config.SkipperFiber(ctx) {
			return ctx.Next()
		}

		if config.BeforeFuncFiber != nil {
			config.BeforeFuncFiber(ctx)
		}
		var (
			clientIP  = ctx.IP()
			method    = ctx.Method()
			userAgent = string(ctx.Request().Header.UserAgent())
			uri       = string(ctx.Request().Header.RequestURI())
			errs      string
			start     = time.Now()
			isError   bool
		)

		logger.WithFields(map[string]interface{}{
			ClientIPField:      clientIP,
			RequestMethodField: method,
			UserAgentField:     userAgent,
			URIField:           uri,
		})

		ctx.Context().SetUserValue(Key, logger)

		err := ctx.Next()

		var (
			statusCode = ctx.Response().StatusCode()
		)
		logger.WithField(StatusField, statusCode)
		if err != nil {
			isError = true
			errs = err.Error()
		}

		if len(errs) > 0 {
			logger.WithField(ErrorsField, errs)
		}
		end := time.Now()
		logger.WithField(EndField, end)
		msg := fmt.Sprintf("latency: %v", end.Sub(start))
		if isError {
			logger.Error(msg)
		} else {
			logger.Info(msg)
		}
		return err
	}
}
