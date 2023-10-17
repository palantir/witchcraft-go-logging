package svc1log

import (
	"bytes"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleLoggers(t *testing.T) {
	wlog.SetDefaultLoggerProvider(wlog.NewJSONMarshalLoggerProvider())
	out1 := &bytes.Buffer{}
	log1 := New(out1, wlog.WarnLevel)
	out2 := &bytes.Buffer{}
	log2 := New(out2, wlog.InfoLevel)

	multi := NewMultiLogger(log1, log2)

	assert.True(t, multi.(wlog.LevelChecker).Enabled(wlog.WarnLevel))
	multi.Warn("warn message")
	assert.Contains(t, out1.String(), `"message":"warn message"`)
	assert.Contains(t, out2.String(), `"message":"warn message"`)

	out1.Reset()
	out2.Reset()
	assert.True(t, multi.(wlog.LevelChecker).Enabled(wlog.InfoLevel))
	multi.Info("info message")
	assert.Empty(t, out1.String())
	assert.Contains(t, out2.String(), `"message":"info message"`)

	out1.Reset()
	out2.Reset()
	assert.False(t, multi.(wlog.LevelChecker).Enabled(wlog.DebugLevel))
	multi.Debug("debug message")
	assert.Empty(t, out1.String())
	assert.Empty(t, out2.String())
}
