package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Log struct {
	*logrus.Entry
	sync.Mutex
	step int32
}

type TextFormatter = logrus.TextFormatter
type JSONFormatter = logrus.JSONFormatter
type Formatter = logrus.Formatter

const (
	Key                = "Logger"
	ClientIPField      = "client_ip"
	RequestMethodField = "request_method"
	UserAgentField     = "user_agent"
	URIField           = "uri"
	StatusField        = "Status"
	ErrorsField        = "Errors"
	EndField           = "end"
	CodeField          = "code"
	RequestField       = "request"
	ResponseField      = "response"
	StartField         = "start"
)

// New return a new log object with log start time
func New(formatter ...Formatter) *Log {
	if len(formatter) > 0 {
		logrus.SetFormatter(formatter[0])
	}
	logTime := time.Now().Format(time.RFC3339)
	return &Log{
		Entry: logrus.WithField("Time", logTime),
	}
}

// GetLogger get logger from context
func GetLogger(ctx context.Context) *Log {
	logger, ok := ctx.Value(Key).(*Log)
	if !ok {
		logger := New()
		ctx = context.WithValue(
			ctx,
			Key, logger)
		return logger
	}
	return logger
}

// ToJsonString convert an object into json string to beautify log
// return nil if marshalling error
func (l *Log) ToJsonString(input interface{}) string {
	if bytes, err := json.Marshal(input); err == nil {
		return string(bytes)
	}
	return ""
}

func (l *Log) addStep() int32 {
	l.Lock()
	defer l.Unlock()
	l.step += 1
	return l.step
}

// AddLog add a new field to log with step = current step + 1
func (l *Log) AddLog(line string, format ...interface{}) *Log {
	step := fmt.Sprintf("STEP_%d", l.addStep())
	if len(format) > 0 {
		logLine := fmt.Sprintf(line, format...)
		l.Entry = l.Entry.WithField(step, logLine)
		return l
	}
	l.Entry = l.Entry.WithField(step, line)
	return l
}

// WithField add a new key = value to log with key = field, value = value
func (l *Log) WithField(field string, value interface{}) *Log {
	l.Entry = l.Entry.WithField(field, value)
	return l
}

// WithFields add multiple key/value to log: key1 = value1, key2 = value2
func (l *Log) WithFields(fields map[string]interface{}) *Log {
	l.Entry = l.Entry.WithFields(fields)
	return l
}