package wlog

import (
	"github.com/palantir/pkg/datetime"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	wparams "github.com/palantir/witchcraft-go-params"
	"time"
)

type ConjureLogType interface {
	ConjureLeveledLogType | ConjureUnleveledLogType
}

type ConjureLeveledLogType interface {
	logging.ServiceLogV1
}

type ConjureUnleveledLogType interface {
	logging.AuditLogV2 | logging.DiagnosticLogV1 | logging.EventLogV1 | logging.MetricLogV1 | logging.RequestLogV2 | logging.TraceLogV1
}

type ConjureLogger[T ConjureUnleveledLogType] interface {
	Log(params ...ConjureLogParam[T])
}

type ConjureLeveledLogger[T ConjureLeveledLogType] interface {
	Debug(message string, params ...ConjureLogParam[T])
	Info(message string, params ...ConjureLogParam[T])
	Warn(message string, params ...ConjureLogParam[T])
	Error(message string, params ...ConjureLogParam[T])
	SetLevel(level LogLevel)
}

type ConjureLogParam[T ConjureLogType] func(*T)

type ConjureLogEntry[T ConjureLogType] struct {
	Log *T
}

type ConjureLogEntryBuilder[T ConjureLogType] interface {
	Build() T
}

type ServiceLogV1Param = func(l *logging.ServiceLogV1)

var _ ConjureLogParam[logging.ServiceLogV1] = ServiceLogV1Param(nil)

type ServiceLogV1Param2 interface {
	apply(l *logging.ServiceLogV1)
}

func Type(typ string) ServiceLogV1Param2 {
	return svc1logParamType{typ: typ}
}

type svc1logParamType struct{ typ string }

func (p svc1logParamType) apply(l *logging.ServiceLogV1) { l.Type = p.typ }

func Time(tim datetime.DateTime) ServiceLogV1Param2 {
	return svc1logParamTime{tim: tim}
}

type svc1logParamTime struct{ tim datetime.DateTime }

func (p svc1logParamTime) apply(l *logging.ServiceLogV1) { l.Time = p.tim }

func Level(level LogLevel) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Level = logging.New_LogLevel(logging.LogLevel_Value(level)) }
}

func Origin(origin string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Origin = &origin }
}

func Thread(thread string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Thread = &thread }
}

func Message(message string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Message = message }
}

func SID(sessionID string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Sid = (*logging.SessionId)(&sessionID) }
}

func TokenID(tokenID string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.TokenId = (*logging.TokenId)(&tokenID) }
}

func TraceID(traceID string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.TraceId = (*logging.TraceId)(&traceID) }
}

func UID(userID string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) { l.Uid = (*logging.UserId)(&userID) }
}

func Params(paramStorer wparams.ParamStorer) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		SafeParams(paramStorer.SafeParams())(l)
		UnsafeParams(paramStorer.UnsafeParams())(l)
	}
}

func SafeParam(key string, value any) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = make(map[string]any)
		}
		l.Params[key] = value
	}
}

func SafeParams(safeParams map[string]any) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = make(map[string]any, len(safeParams))
		}
		for k, v := range safeParams {
			l.Params[k] = v
		}
	}
}

func UnsafeParam(key string, value any) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = make(map[string]any)
		}
		l.UnsafeParams[key] = value
	}
}

func UnsafeParams(unsafeParams map[string]any) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = make(map[string]any, len(unsafeParams))
		}
		for k, v := range unsafeParams {
			l.UnsafeParams[k] = v
		}
	}
}

func Tag(key string, value string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.Tags == nil {
			l.Tags = make(map[string]string)
		}
		l.Tags[key] = value
	}
}

func Tags(tags map[string]string) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if l.Tags == nil {
			l.Tags = make(map[string]string, len(tags))
		}
		for k, v := range tags {
			l.Tags[k] = v
		}
	}
}

func Stacktrace(err error) ServiceLogV1Param {
	return func(l *logging.ServiceLogV1) {
		if err == nil {
			return
		}
		errString := werror.GenerateErrorString(err, false)
		l.Stacktrace = &errString
		// add all safe and unsafe parameters stored in error
		safeParams, unsafeParams := werror.ParamsFromError(err)
		SafeParams(safeParams)(l)
		UnsafeParams(unsafeParams)(l)
	}
}

func foo() {
	// Create a service log entry with a message and a timestamp
	logEntry := NewServiceLogV1EntryBuilder().
		Message("Hello, World!").
		Build()
}

func NewServiceLogV1EntryBuilder(params ...ServiceLogV1Param) *ServiceLogV1Builder {
	b := ServiceLogV1Builder{}
	for _, param := range params {
		param(&b.l)
	}
	return &b
}

type ServiceLogV1Builder struct {
	l logging.ServiceLogV1
}

func (b *ServiceLogV1Builder) Type(typ string) *ServiceLogV1Builder {
	b.l.Type = typ
	return b
}

func (b *ServiceLogV1Builder) Time(timestamp datetime.DateTime) *ServiceLogV1Builder {
	b.l.Time = timestamp
	return b
}

func (b *ServiceLogV1Builder) Level(level string) *ServiceLogV1Builder {
	b.l.Level = logging.New_LogLevel(logging.LogLevel_Value(level))
	return b
}

func (b *ServiceLogV1Builder) Origin(origin string) *ServiceLogV1Builder {
	b.l.Origin = &origin
	return b
}

func (b *ServiceLogV1Builder) Thread(thread string) *ServiceLogV1Builder {
	b.l.Thread = &thread
	return b
}

func (b *ServiceLogV1Builder) Message(message stringConst) *ServiceLogV1Builder {
	b.l.Message = message
	return b
}

func (b *ServiceLogV1Builder) SID(sessionID string) *ServiceLogV1Builder {
	b.l.Sid = (*logging.SessionId)(&sessionID)
	return b
}

func (b *ServiceLogV1Builder) TokenID(tokenID string) *ServiceLogV1Builder {
	b.l.TokenId = (*logging.TokenId)(&tokenID)
	return b
}

func (b *ServiceLogV1Builder) TraceID(traceID string) *ServiceLogV1Builder {
	b.l.TraceId = (*logging.TraceId)(&traceID)
	return b
}

func (b *ServiceLogV1Builder) UID(userID string) *ServiceLogV1Builder {
	b.l.Uid = (*logging.UserId)(&userID)
	return b
}

func (b *ServiceLogV1Builder) Params(paramStorer wparams.ParamStorer) *ServiceLogV1Builder {
	b.SafeParams(paramStorer.SafeParams())
	b.UnsafeParams(paramStorer.UnsafeParams())
	return b
}

func (b *ServiceLogV1Builder) SafeParam(key string, value any) *ServiceLogV1Builder {
	if b.l.Params == nil {
		b.l.Params = make(map[string]any)
	}
	b.l.Params[key] = value
	return b
}

func (b *ServiceLogV1Builder) SafeParams(safeParams map[string]any) *ServiceLogV1Builder {
	if b.l.Params == nil {
		b.l.Params = make(map[string]any, len(safeParams))
	}
	for k, v := range safeParams {
		// TODO: Recursively merge maps?
		b.l.Params[k] = v
	}
	return b
}

func (b *ServiceLogV1Builder) UnsafeParam(key string, value any) *ServiceLogV1Builder {
	if b.l.UnsafeParams == nil {
		b.l.UnsafeParams = make(map[string]any)
	}
	b.l.UnsafeParams[key] = value
	return b
}

func (b *ServiceLogV1Builder) UnsafeParams(unsafeParams map[string]any) *ServiceLogV1Builder {
	if b.l.UnsafeParams == nil {
		b.l.UnsafeParams = make(map[string]any, len(unsafeParams))
	}
	for k, v := range unsafeParams {
		// TODO: Recursively merge maps?
		b.l.UnsafeParams[k] = v
	}
	return b
}

func (b *ServiceLogV1Builder) Tag(key string, value string) *ServiceLogV1Builder {
	if b.l.Tags == nil {
		b.l.Tags = make(map[string]string)
	}
	b.l.Tags[key] = value
	return b
}

func (b *ServiceLogV1Builder) Tags(tags map[string]string) *ServiceLogV1Builder {
	if b.l.Tags == nil {
		b.l.Tags = make(map[string]string, len(tags))
	}
	for k, v := range tags {
		b.l.Tags[k] = v
	}
	return b
}

func (b *ServiceLogV1Builder) Stacktrace(stacktrace string) *ServiceLogV1Builder {
	b.l.Stacktrace = &stacktrace
	return b
}

func (b *ServiceLogV1Builder) Build() logging.ServiceLogV1 {
	if b.l.Type == "" {
		b.l.Type = "service.1"
	}
	if time.Time(b.l.Time).IsZero() {
		b.l.Time = datetime.DateTime(time.Now())
	}
	return b.l
}
