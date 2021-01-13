package gologrwlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"

	// Use zap as logger implementation
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
)

type logEntry struct {
	Level      wlog.LogLevel          `json:"level"`
	Message    string                 `json:"message"`
	Type       string                 `json:"type"`
	Origin     string                 `json:"origin"`
	Params     map[string]interface{} `json:"params"`
	Stacktrace string                 `json:"stacktrace"`
}

func TestSvc1LogrWrapper(t *testing.T) {
	byteWriter := new(bytes.Buffer)
	logger := svc1log.New(byteWriter, wlog.DebugLevel)

	logr1 := NewGoLogrWrapper(logger, "foo")
	logr2 := logr1.WithValues("key2", "val2", "key3", "val3")
	logr3 := logr1.WithName("bar")

	logr1.Info("logr 1", "key1", "val1")
	assertLogLine(t, byteWriter.Bytes(), logEntry{
		Level:   wlog.InfoLevel,
		Message: "logr 1",
		Type:    svc1log.TypeValue,
		Origin:  "foo",
		Params: map[string]interface{}{
			"key1": "val1",
		},
	})
	byteWriter.Reset()

	logr2.Info("logr 2")
	assertLogLine(t, byteWriter.Bytes(), logEntry{
		Level:   wlog.InfoLevel,
		Message: "logr 2",
		Type:    svc1log.TypeValue,
		Origin:  "foo",
		Params: map[string]interface{}{
			"key2": "val2",
			"key3": "val3",
		},
	})
	byteWriter.Reset()

	logr3.Error(fmt.Errorf("test error"), "logr 3")
	assertLogLine(t, byteWriter.Bytes(), logEntry{
		Level:      wlog.ErrorLevel,
		Message:    "logr 3",
		Type:       svc1log.TypeValue,
		Origin:     "foo/bar",
		Stacktrace: "test error",
	})
	byteWriter.Reset()
}

func assertLogLine(t *testing.T, logLine []byte, expectedLogEntry logEntry) {
	fmt.Println(string(logLine))
	logEntry := new(logEntry)
	err := json.Unmarshal(logLine, &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, expectedLogEntry, *logEntry)
}

