package svc1log_literal

import (
	"io"
	golog "log"
	"sync"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

const (
	TypeValue = "service.1"
)

var (
	svclogObjPool = &sync.Pool{New: func() interface{} {
		return &logging.ServiceLogV1{Type: TypeValue}
	}}
	//svclogBufPool = &sync.Pool{New: func() interface{} {
	//	b := make([]byte, 0, 1024)
	//	return &b
	//}}

	levelDebug = logging.New_LogLevel(logging.LogLevel_DEBUG)
	levelInfo  = logging.New_LogLevel(logging.LogLevel_INFO)
	levelWarn  = logging.New_LogLevel(logging.LogLevel_WARN)
	levelError = logging.New_LogLevel(logging.LogLevel_ERROR)
)

type Logger interface {
	Debug(msg string, params ...Param)
	Info(msg string, params ...Param)
	Warn(msg string, params ...Param)
	Error(msg string, params ...Param)
	SetLevel(level wlog.LogLevel)
}

func NewConjureLogger(w io.Writer, level wlog.LogLevel, params ...Param) Logger {
	return &conjureLogger{
		output:         w,
		params:         params,
		nl:             []byte("\n"),
		AtomicLogLevel: wlog.NewAtomicLogLevel(level),
	}
}

type conjureLogger struct {
	output io.Writer
	params []Param
	nl     []byte
	*wlog.AtomicLogLevel
}

func (l *conjureLogger) Debug(msg string, params ...Param) {
	l.doLog(wlog.DebugLevel, levelDebug, msg, params...)
}

func (l *conjureLogger) Info(msg string, params ...Param) {
	l.doLog(wlog.InfoLevel, levelInfo, msg, params...)
}

func (l *conjureLogger) Warn(msg string, params ...Param) {
	l.doLog(wlog.WarnLevel, levelWarn, msg, params...)
}

func (l *conjureLogger) Error(msg string, params ...Param) {
	l.doLog(wlog.ErrorLevel, levelError, msg, params...)
}

func (l *conjureLogger) doLog(wLevel wlog.LogLevel, cLevel logging.LogLevel, msg string, params ...Param) {
	if l.Enabled(wLevel) {
		log := svclogObjPool.Get().(*logging.ServiceLogV1)
		defer resetSvc1Log(log)
		log.Type = TypeValue
		log.Time = datetime.DateTime(time.Now())
		log.Level = cLevel
		log.Message = msg
		for _, p := range l.params {
			p.apply(log)
		}
		for _, p := range params {
			p.apply(log)
		}

		out, err := log.MarshalJSON() // TODO: reuse memory with a different json writer interface
		if err != nil {
			golog.Printf("failed to marshal service.1 log: %v", err)
			return
		}
		if _, err := l.output.Write(out); err != nil {
			golog.Printf("failed to write service.1 log: %v", err)
			golog.Println(string(out))
		}
	}
}

func resetSvc1Log(log *logging.ServiceLogV1) {
	log.Type = TypeValue
	log.Level = logging.LogLevel{}
	log.Time = datetime.DateTime{}
	log.Origin = nil
	log.Thread = nil
	log.Message = ""
	clear(log.Params)
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	log.Stacktrace = nil
	clear(log.UnsafeParams)
	clear(log.Tags)
	// return log to pool
	svclogObjPool.Put(log)
}
