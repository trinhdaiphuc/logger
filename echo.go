package logger

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

var DefaultConfigEcho = ConfigEcho{
	SkipperEcho: DefaultSkipperEcho,
}

func EchoMiddleware(config ConfigEcho) echo.MiddlewareFunc {
	if config.SkipperEcho == nil {
		config.SkipperEcho = DefaultSkipperEcho
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			logger := New(WithFormatter(&JSONFormatter{}))
			if config.SkipperEcho(ctx) {
				return next(ctx)
			}

			if config.BeforeFuncEcho != nil {
				config.BeforeFuncEcho(ctx)
			}
			var (
				clientIP  = ctx.RealIP()
				method    = ctx.Request().Method
				userAgent = ctx.Request().UserAgent()
				uri       = ctx.Request().RequestURI
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

			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), Key, logger)))
			ctx.Set(Key, logger)

			err := next(ctx)

			var (
				statusCode = ctx.Response().Status
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
}
