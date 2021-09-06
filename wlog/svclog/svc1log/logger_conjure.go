package svc1log

import (
	"io"
	"sync"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

var (
	svclogObjPool = &sync.Pool{New: func() interface{} {
		return new(logging.ServiceLogV1)
	}}
	svclogBufPool = &sync.Pool{New: func() interface{} {
		b := make([]byte, 0, 128)
		return &b
	}}
)

func NewConjureLogger(w io.Writer, level wlog.LogLevel, params ...Param) Logger {

}

type conjureLogger struct {
	output io.Writer
	level  wlog.LogLevel
}

func (l *conjureLogger) Debug(msg string, params ...Param) {
	panic("implement me")
}

func (l *conjureLogger) Info(msg string, params ...Param) {
	panic("implement me")
}

func (l *conjureLogger) Warn(msg string, params ...Param) {
	panic("implement me")
}

func (l *conjureLogger) Error(msg string, params ...Param) {
	panic("implement me")
}

func (l *conjureLogger) SetLevel(level wlog.LogLevel) {
	panic("implement me")
}
