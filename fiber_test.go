package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFiberMiddleware(t *testing.T) {
	testErr := errors.New("test err")
	tests := []struct {
		name               string
		config             ConfigFiber
		route              func(server *fiber.App)
		requestUri         string
		response           string
		err                error
		responseStatusCode int
		logLevel           string
	}{
		{
			name:   "Test nil config success",
			config: ConfigFiber{},
			route: func(server *fiber.App) {
				server.Get("/hello/:name", func(ctx *fiber.Ctx) error {
					logger := GetLogger(ctx.Context())
					name := ctx.Params("name")
					logger.AddLog("request name %v", name)
					return ctx.Status(200).SendString("hello " + name)
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
			config: DefaultConfigFiber,
			route: func(server *fiber.App) {
				server.Get("/hello/:name", func(ctx *fiber.Ctx) error {
					logger := GetLogger(ctx.Context())
					name := ctx.Params("name")
					logger.AddLog("request name %v", name)
					return ctx.Status(200).SendString("hello " + name)
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
			config: ConfigFiber{
				SkipperFiber: func(context *fiber.Ctx) bool {
					if string(context.Request().RequestURI()) == "/metrics" {
						return true
					}
					return false
				},
			},
			route: func(server *fiber.App) {
				server.Get("/metrics", func(ctx *fiber.Ctx) error {
					logger := GetLogger(ctx.Context())
					logger.AddLog("metrics")
					return ctx.Status(200).SendString("success")
				})
			},
			logLevel:           "info",
			requestUri:         "/metrics",
			err:                nil,
			responseStatusCode: http.StatusOK,
			response:           "success",
		},
		{
			name: "Test error",
			config: ConfigFiber{
				BeforeFuncFiber: func(context *fiber.Ctx) {
					context.Set("X-Metadata", "hello")
				},
			},
			route: func(server *fiber.App) {
				server.Get("/hello/:name", func(ctx *fiber.Ctx) error {
					logger := GetLogger(ctx.Context())
					name := ctx.Params("name")
					logger.AddLog("request name %v", name)
					ctx.Status(500)
					return testErr
				})
			},
			requestUri:         "/hello/world",
			response:           testErr.Error(),
			err:                nil,
			responseStatusCode: http.StatusInternalServerError,
			logLevel:           "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				buf = &bytes.Buffer{}
			)
			New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
			server := fiber.New()
			server.Use(FiberMiddleware(tt.config))
			tt.route(server)
			w, e := performFiberRequest(server, "GET", tt.requestUri)

			out, e := ioutil.ReadAll(w.Body)

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
			assert.Equal(t, tt.responseStatusCode, w.StatusCode)
		})
	}
}

func TestFiberMiddlewareError(t *testing.T) {
	var (
		buf     = &bytes.Buffer{}
		testErr = errors.New("test err")
	)
	New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
	server := fiber.New()
	server.Use(FiberMiddleware(DefaultConfigFiber))

	server.Get("/hello/:name", func(ctx *fiber.Ctx) error {
		logger := GetLogger(ctx.Context())
		name := ctx.Params("name")
		logger.AddLog("request name %v", name)
		ctx.Status(500)
		return testErr
	})

	// RUN
	requestUri := "/hello/world"
	w, e := performFiberRequest(server, "GET", requestUri)

	out, e := ioutil.ReadAll(w.Body)

	t.Logf("Out %v, err %v", string(out), e)
	t.Logf("Log output %v", buf.String())

	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data["STEP_1"]
	uri, uriOk := data[URIField]
	_, errOk := data[ErrorsField]
	level, levelOk := data[FieldKeyLevel]

	// TEST
	assert.True(t, ok, `cannot found expected "STEP_1" field: %v`, data)
	assert.True(t, uriOk, `cannot found expected "%v" field: %v`, URIField, data)
	assert.True(t, errOk, `cannot found expected "%v" field: %v`, ErrorsField, data)
	assert.True(t, levelOk, `cannot found expected "%v" field: %v`, FieldKeyLevel, data)
	assert.Equal(t, level, "error")
	assert.Equal(t, requestUri, uri)
	assert.Equal(t, nil, e)
	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}

func performFiberRequest(app *fiber.App, method, path string) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	return app.Test(req)
}
