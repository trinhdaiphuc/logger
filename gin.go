package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var DefaultConfigGin = ConfigGin{
	SkipperGin: DefaultSkipperGin,
}

func GinMiddleware(config ConfigGin) gin.HandlerFunc {
	if config.SkipperGin == nil {
		config.SkipperGin = DefaultSkipperGin
	}
	logger := New(WithFormatter(&JSONFormatter{}))

	return func(ctx *gin.Context) {
		if config.SkipperGin(ctx) {
			ctx.Next()
			return
		}

		if config.BeforeFuncGin != nil {
			config.BeforeFuncGin(ctx)
		}
		var (
			clientIP = ctx.ClientIP()
			method   = ctx.Request.Method

			userAgent = ctx.Request.UserAgent()
			uri       = ctx.Request.RequestURI
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
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), Key, logger))
		ctx.Set(Key, logger)
		ctx.Next()
		var (
			statusCode = ctx.Writer.Status()
		)
		logger.WithField(StatusField, statusCode)
		if ctx.Errors != nil {
			isError = true
			bs, err := ctx.Errors.MarshalJSON()
			if err == nil {
				errs = string(bs)
			} else {
				errs = ctx.Errors.String()
			}
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
	}
}
