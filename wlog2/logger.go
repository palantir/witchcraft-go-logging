package wlog

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
)

// ConjureLogType is a constraint for generic types that combines all the Conjure log types.
type ConjureLogType interface {
	logging.AuditLogV2 | logging.DiagnosticLogV1 | logging.EventLogV2 | logging.MetricLogV1 | logging.RequestLogV2 | logging.ServiceLogV1 | logging.TraceLogV1 | logging.WrappedLogV1
}

// ConjureLogger is a generic logger that can log all Conjure log types.
type ConjureLogger[T ConjureLogType] interface {
	Log(params ...ConjureLogParam[T])
	WithParams(params ...ConjureLogParam[T]) ConjureLogger[T]
}

// ConjureLogParam is a function that modifies a Conjure log object.
type ConjureLogParam[T ConjureLogType] func(*T)

// NewDefaultLogger creates a new logger that writes JSON-marshaled lines to the provided output.
func NewDefaultLogger[T ConjureLogType](output io.Writer, params ...ConjureLogParam[T]) ConjureLogger[T] {
	return NewDefaultLoggerWithPrinter(jsonPrinter[T]{output}, params...)
}

// NewDefaultLoggerWithPrinter creates a new logger that writes log objects using the provided printer.
func NewDefaultLoggerWithPrinter[T ConjureLogType](printer ConjureLogPrinter[T], params ...ConjureLogParam[T]) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: printer,
		params:  params,
		objPool: &sync.Pool{
			New: func() any { return new(T) },
		},
	}
}

type defaultLogger[T ConjureLogType] struct {
	params  []ConjureLogParam[T]
	printer ConjureLogPrinter[T]
	objPool *sync.Pool
}

func (l *defaultLogger[T]) WithParams(params ...ConjureLogParam[T]) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: l.printer,
		params:  append(append([]ConjureLogParam[T]{}, l.params...), params...),
		objPool: l.objPool,
	}
}

func (l *defaultLogger[T]) Log(params ...ConjureLogParam[T]) {
	log := l.objPool.Get().(*T)
	defer l.resetLogObject(log)
	for _, p := range l.params {
		if p != nil {
			p(log)
		}
	}
	for _, p := range params {
		if p != nil {
			p(log)
		}
	}
	if err := l.printer.Print(*log); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		// TODO: something else?
	}
}

func (l *defaultLogger[T]) resetLogObject(in *T) {
	switch log := any(in).(type) {
	case *logging.AuditLogV2:
		log.Time = datetime.DateTime{}
		log.Uid = nil
		log.Sid = nil
		log.TokenId = nil
		log.TraceId = nil
		log.OtherUids = nil
		log.Origin = nil
		log.Name = ""
		log.Result = logging.AuditResult{}
		clear(log.ResultParams)
		clear(log.RequestParams)

	case *logging.DiagnosticLogV1:
		log.Time = datetime.DateTime{}
		log.Diagnostic = logging.Diagnostic{}
		clear(log.UnsafeParams)

	case *logging.EventLogV2:
		log.Time = datetime.DateTime{}
		log.EventName = ""
		clear(log.Values)
		log.Uid = nil
		log.Sid = nil
		log.TokenId = nil
		log.TraceId = nil
		clear(log.UnsafeParams)
		clear(log.Tags)

	case *logging.MetricLogV1:
		log.Time = datetime.DateTime{}
		log.MetricName = ""
		log.MetricType = ""
		clear(log.Values)
		clear(log.Tags)
		log.Uid = nil
		log.Sid = nil
		log.TokenId = nil
		clear(log.UnsafeParams)

	case *logging.RequestLogV2:
		log.Time = datetime.DateTime{}
		log.Method = nil
		log.Protocol = ""
		log.Path = ""
		clear(log.Params)
		log.Status = 0
		log.RequestSize = 0
		log.ResponseSize = 0
		log.Duration = 0
		log.Uid = nil
		log.Sid = nil
		log.TokenId = nil
		log.TraceId = nil
		clear(log.UnsafeParams)

	case *logging.ServiceLogV1:
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

	case *logging.TraceLogV1:
		log.Time = datetime.DateTime{}
		log.Uid = nil
		log.Sid = nil
		log.TokenId = nil
		clear(log.UnsafeParams)
		log.Span.TraceId = ""
		log.Span.Id = ""
		log.Span.Name = ""
		log.Span.ParentId = nil
		log.Span.Timestamp = 0
		log.Span.Duration = 0
		clear(log.Span.Annotations)
		clear(log.Span.Tags)

	case *logging.WrappedLogV1:
		log.Payload = logging.WrappedLogV1Payload{}
		log.EntityName = ""
		log.EntityVersion = ""

	default:
		panic(fmt.Sprintf("unexpected log type: %T", log))
	}
	l.objPool.Put(in)
}
