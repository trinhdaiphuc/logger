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
	var (
		l   = New(&JSONFormatter{})
		buf = &bytes.Buffer{}
	)
	server := fiber.New()
	server.Use(FiberMiddleware())
	l.Logger.SetOutput(buf)
	server.Get("/hello/:name", func(ctx *fiber.Ctx) error {
		logger := GetLogger(ctx.Context())
		name := ctx.Params("name")
		logger.AddLog("request name %v", name)
		return ctx.Status(200).SendString( "hello "+name)
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
	level, levelOk := data[FieldKeyLevel]

	// TEST
	assert.True(t, ok, `cannot found expected "STEP_1" field: %v`, data)
	assert.True(t, uriOk, `cannot found expected "%v" field: %v`, URIField, data)
	assert.Equal(t, requestUri, uri)
	assert.True(t, levelOk, `cannot found expected "%v" field: %v`, FieldKeyLevel, data)
	assert.Equal(t, level, "info")
	assert.Equal(t, nil, e)
	assert.Equal(t, "hello world", string(out))
	assert.Equal(t, http.StatusOK, w.StatusCode)
}

func TestFiberMiddlewareError(t *testing.T) {
	var (
		l       = New(&JSONFormatter{})
		buf     = &bytes.Buffer{}
		testErr = errors.New("test err")
	)
	server := fiber.New()
	server.Use(FiberMiddleware())
	l.Logger.SetOutput(buf)
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
