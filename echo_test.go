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
	var (
		l   = New(&JSONFormatter{})
		buf = &bytes.Buffer{}
	)
	server := echo.New()
	server.Use(EchoMiddleware)
	l.Logger.SetOutput(buf)
	server.GET("/hello/:name", func(ctx echo.Context) error {
		logger := GetLogger(ctx.Request().Context())
		name := ctx.Param("name")
		logger.AddLog("request name %v", name)
		return ctx.String(200, "hello "+name)
	})

	// RUN
	requestUri := "/hello/world"
	w := performRequest(server, "GET", requestUri)

	out, e := ioutil.ReadAll(w.Result().Body)

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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEchoMiddlewareError(t *testing.T) {
	var (
		l       = New(&JSONFormatter{})
		buf     = &bytes.Buffer{}
		testErr = errors.New("test err")
	)
	server := echo.New()
	server.Use(EchoMiddleware)
	l.Logger.SetOutput(buf)
	server.GET("/hello/:name", func(ctx echo.Context) error {
		logger := GetLogger(ctx.Request().Context())
		name := ctx.Param("name")
		logger.AddLog("request name %v", name)
		ctx.Error(testErr)
		return testErr
	})

	// RUN
	requestUri := "/hello/world"
	w := performRequest(server, "GET", requestUri)

	out, e := ioutil.ReadAll(w.Result().Body)

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
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
