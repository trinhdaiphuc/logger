package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGinMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		config             ConfigGin
		route              func(server *gin.Engine)
		requestUri         string
		response           string
		err                error
		responseStatusCode int
		logLevel           string
	}{
		{
			name:   "Test nil config success",
			config: ConfigGin{},
			route: func(server *gin.Engine) {
				server.GET("/hello/:name", func(ctx *gin.Context) {
					logger := GetLogger(ctx)
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					ctx.String(200, "hello "+name)
				})
			},
			requestUri:         "/hello/world",
			response:           "hello world",
			err:                nil,
			responseStatusCode: http.StatusOK,
			logLevel:           "info",
		},
		{
			name:   "Test default config success",
			config: DefaultConfigGin,
			route: func(server *gin.Engine) {
				server.GET("/hello/:name", func(ctx *gin.Context) {
					logger := GetLogger(ctx)
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					ctx.String(200, "hello "+name)
				})
			},
			requestUri:         "/hello/world",
			response:           "hello world",
			err:                nil,
			responseStatusCode: http.StatusOK,
			logLevel:           "info",
		},
		{
			name: "Test skipp config success",
			config: ConfigGin{
				SkipperGin: func(context *gin.Context) bool {
					if context.Request.RequestURI == "/metrics" {
						return true
					}
					return false
				},
			},
			route: func(server *gin.Engine) {
				server.GET("/metrics", func(ctx *gin.Context) {
					logger := GetLogger(ctx)
					logger.AddLog("metrics")
					ctx.String(200, "success")
				})
			},
			requestUri:         "/metrics",
			err:                nil,
			responseStatusCode: http.StatusOK,
			response:           "success",
		},
		{
			name: "Test error",
			config: ConfigGin{
				BeforeFuncGin: func(context *gin.Context) {
					context.Header("X-Metadata", "hello")
				},
			},
			route: func(server *gin.Engine) {
				server.GET("/hello/:name", func(ctx *gin.Context) {
					logger := GetLogger(ctx)
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					ctx.Error(errors.New("test err"))
					ctx.Status(500)
				})
			},
			requestUri:         "/hello/world",
			response:           "",
			err:                nil,
			responseStatusCode: http.StatusInternalServerError,
			logLevel:           "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
			server := gin.New()
			server.Use(GinMiddleware(tt.config))
			tt.route(server)

			w := performRequest(server, "GET", tt.requestUri)

			out, e := ioutil.ReadAll(w.Result().Body)

			t.Logf("Out %v, err %v", string(out), e)
			t.Logf("Log output %v", buf.String())

			var data map[string]interface{}
			if len(buf.String()) > 0 {
				if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
					t.Error("unexpected error", err)
				}
				_, ok := data["STEP_1"]
				uri, uriOk := data[URIField]
				level, levelOk := data[FieldKeyLevel]
				assert.True(t, ok, `cannot found expected "STEP_1" field: %v`, data)
				assert.True(t, uriOk, `cannot found expected "%v" field: %v`, URIField, data)
				assert.Equal(t, tt.requestUri, uri)
				assert.True(t, levelOk, `cannot found expected "%v" field: %v`, FieldKeyLevel, data)
				assert.Equal(t, tt.logLevel, level)
			}

			// TEST
			assert.Equal(t, tt.err, e)
			assert.Equal(t, tt.response, string(out))
			assert.Equal(t, tt.responseStatusCode, w.Code)
		})
	}
}
