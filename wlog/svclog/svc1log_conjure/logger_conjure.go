package svc1log

import (
	"io"
	golog "log"
	"sync"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
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
		//marshaler:      (*logging.ServiceLogV1).WriteJSON,
		output:         w,
		params:         params,
		nl:             []byte("\n"),
		AtomicLogLevel: wlog.NewAtomicLogLevel(level),
	}
}

//type encoderFunc func(log *logging.ServiceLogV1, writer io.Writer) (int, error)

type conjureLogger struct {
	//marshaler encoderFunc
	output io.Writer
	params []Param
	nl     []byte
	*wlog.AtomicLogLevel
}

func (l *conjureLogger) Debug(msg string, params ...Param) {
	if l.Enabled(wlog.DebugLevel) {
		log := svclogObjPool.Get().(*logging.ServiceLogV1)
		defer resetSvc1Log(log)
		log.Time = datetime.DateTime(time.Now())
		log.Level = levelDebug
		log.Message = msg
		l.applyParams(log, params...)
		l.write(log)
	}
}

func (l *conjureLogger) Info(msg string, params ...Param) {
	if l.Enabled(wlog.InfoLevel) {
		log := svclogObjPool.Get().(*logging.ServiceLogV1)
		defer resetSvc1Log(log)
		log.Type = TypeValue
		log.Time = datetime.DateTime(time.Now())
		log.Level = levelInfo
		log.Message = msg
		l.applyParams(log, params...)
		l.write(log)
	}
}

func (l *conjureLogger) Warn(msg string, params ...Param) {
	if l.Enabled(wlog.WarnLevel) {
		log := svclogObjPool.Get().(*logging.ServiceLogV1)
		defer resetSvc1Log(log)
		log.Time = datetime.DateTime(time.Now())
		log.Level = levelWarn
		log.Message = msg
		l.applyParams(log, params...)
		l.write(log)
	}
}

func (l *conjureLogger) Error(msg string, params ...Param) {
	if l.Enabled(wlog.ErrorLevel) {
		log := svclogObjPool.Get().(*logging.ServiceLogV1)
		defer resetSvc1Log(log)
		log.Time = datetime.DateTime(time.Now())
		log.Level = levelError
		log.Message = msg
		l.applyParams(log, params...)
		l.write(log)
	}
}

func (l *conjureLogger) applyParams(log *logging.ServiceLogV1, params ...Param) {
	for _, p := range l.params {
		p.apply(log)
	}
	for _, p := range params {
		p.apply(log)
	}
}

func (l *conjureLogger) write(log *logging.ServiceLogV1) {
	//buf := svclogBufPool.Get().(*[]byte)
	//defer svclogBufPool.Put(buf)
	//*buf = (*buf)[:0]
	//_, err := l.marshaler(log, dj.NewAppender(buf))
	out, err := log.MarshalJSON()
	if err != nil {
		golog.Printf("failed to marshal service.1 log: %v", err)
		return
	}
	if _, err := l.output.Write(out); err != nil {
		golog.Printf("failed to write service.1 log: %v", err)
		golog.Println(string(out))
	}
}

func resetSvc1Log(log *logging.ServiceLogV1) {
	log.Type = TypeValue
	log.Level = logging.LogLevel{}
	log.Time = datetime.DateTime{}
	log.Origin = nil
	log.Thread = nil
	log.Message = ""
	for k := range log.Params {
		delete(log.Params, k)
	}
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	log.Stacktrace = nil
	for k := range log.UnsafeParams {
		delete(log.UnsafeParams, k)
	}
	for k := range log.Tags {
		delete(log.Tags, k)
	}
	// return log to pool
	svclogObjPool.Put(log)
}
