package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoMiddleware(t *testing.T) {
	testErr := errors.New("test err")
	tests := []struct {
		name               string
		config             ConfigEcho
		route              func(server *echo.Echo)
		requestUri         string
		response           string
		err                error
		responseStatusCode int
		logLevel           string
	}{
		{
			name:   "Test nil config success",
			config: ConfigEcho{},
			route: func(server *echo.Echo) {
				server.GET("/hello/:name", func(ctx echo.Context) error {
					logger := GetLogger(ctx.Request().Context())
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					return ctx.String(200, "hello "+name)
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
			config: DefaultConfigEcho,
			route: func(server *echo.Echo) {
				server.GET("/hello/:name", func(ctx echo.Context) error {
					logger := GetLogger(ctx.Request().Context())
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					return ctx.String(200, "hello "+name)
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
			config: ConfigEcho{
				SkipperEcho: func(context echo.Context) bool {
					if context.Request().RequestURI == "/metrics" {
						return true
					}
					return false
				},
			},
			route: func(server *echo.Echo) {
				server.GET("/metrics", func(ctx echo.Context) error {
					logger := GetLogger(ctx.Request().Context())
					logger.AddLog("metrics")
					return ctx.String(200, "success")
				})
			},
			requestUri:         "/metrics",
			err:                nil,
			responseStatusCode: http.StatusOK,
			logLevel:           "info",
			response:           "success",
		},
		{
			name: "Test error",
			config: ConfigEcho{BeforeFuncEcho: func(ctx echo.Context) {

			}},
			route: func(server *echo.Echo) {
				server.GET("/hello/:name", func(ctx echo.Context) error {
					logger := GetLogger(ctx.Request().Context())
					name := ctx.Param("name")
					logger.AddLog("request name %v", name)
					ctx.Error(testErr)
					return testErr
				})
			},
			requestUri:         "/hello/world",
			response:           "{\"message\":\"Internal Server Error\"}\n",
			err:                nil,
			responseStatusCode: http.StatusInternalServerError,
			logLevel:           "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
			server := echo.New()
			server.Use(EchoMiddleware(tt.config))
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

				// TEST
				assert.True(t, ok, `cannot found expected "STEP_1" field: %v`, data)
				assert.True(t, uriOk, `cannot found expected "%v" field: %v`, URIField, data)
				assert.Equal(t, tt.requestUri, uri)
				assert.True(t, levelOk, `cannot found expected "%v" field: %v`, FieldKeyLevel, data)
				assert.Equal(t, tt.logLevel, level)
			}
			assert.Equal(t, tt.err, e)
			assert.Equal(t, tt.response, string(out))
			assert.Equal(t, tt.responseStatusCode, w.Code)
		})
	}
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
