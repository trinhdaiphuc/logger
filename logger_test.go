package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type Resp struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
}

// Default key names for the default fields
const (
	FieldKeyLevel       = "level"
	FieldKeyLogrusError = "logrus_error"
)

func TestFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
	l.WithField("func", func() {}).Info("test")
	t.Log("buffer", buf.String())
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data[FieldKeyLogrusError]
	assert.True(t, ok, `cannot found expected "logrus_error" field: %v`, data)
}

func TestContextWithNoLogger(t *testing.T) {
	var (
		buf    = &bytes.Buffer{}
		ctx    = context.Background()
		logger = GetLogger(ctx)
	)
	logger.Logger.SetFormatter(&TextFormatter{})
	logger.Logger.SetOutput(buf)
	logger.AddLog("hello")
	logger.Info("end")
	output := buf.String()
	t.Log("buffer", output)
	ok := strings.Contains(output, "STEP_1=hello")
	assert.True(t, ok, `cannot found expected "STEP_1=hello" field: %v`, output)
	ok = strings.Contains(output, "level=info")
	assert.True(t, ok, `cannot found expected "level=info" field: %v`, output)
}

func TestContextWithLogger(t *testing.T) {
	var (
		buf = &bytes.Buffer{}
		l   = New(WithFormatter(&TextFormatter{}), WithOutput(buf))
		ctx = context.WithValue(context.Background(), Key, l)
	)
	logger := GetLogger(ctx)
	logger.AddLog("hello %v", "world")
	logger.WithFields(map[string]interface{}{
		"K1": "V1",
		"K2": "V2",
	})
	logger.Info("end")
	output := buf.String()
	fmt.Println("buffer", output)
	ok := strings.Contains(output, `STEP_1="hello world"`)
	assert.True(t, ok, `cannot found expected "STEP_1=hello world" field: %v`, output)
	ok = strings.Contains(output, "level=info")
	assert.True(t, ok, `cannot found expected "level=info" field: %v`, output)
	ok = strings.Contains(output, "K1=V1")
	assert.True(t, ok, `cannot found expected "K1=V1" field: %v`, output)
	ok = strings.Contains(output, "K2=V2")
	assert.True(t, ok, `cannot found expected "K2=V2" field: %v`, output)
}

func TestLoggerJson(t *testing.T) {
	var (
		data = Resp{
			Message: "OK",
			Code:    200,
		}
		errData = func() {}
	)
	loggerStr := ToJsonString(data)
	buf, err := json.Marshal(data)
	expect := string(buf)
	t.Logf("Logger json to string %v, data expect %v", loggerStr, expect)
	if err != nil {
		assert.Errorf(t, err, "Marshal object failed")
	}

	assert.Equalf(t, loggerStr, expect, "Expected string data %v, output logger %v", expect, loggerStr)
	emptyStr := ToJsonString(errData)
	t.Logf("Logger empty to string %v", emptyStr)
	assert.Equalf(t, emptyStr, "", "Expected empty string data, output logger %v", loggerStr)
}
