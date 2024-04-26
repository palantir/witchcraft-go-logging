package svc1log

import (
	"maps"
	"strconv"
	"time"

	"github.com/palantir/pkg/datetime"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wparams "github.com/palantir/witchcraft-go-params"
)

type Param = wlog.ConjureLogParam[logging.ServiceLogV1]

func Type() Param {
	return func(l *logging.ServiceLogV1) { l.Type = "service.1" }
}

func Level(level wlog.LogLevel) Param {
	return func(l *logging.ServiceLogV1) { l.Level = logging.New_LogLevel(logging.LogLevel_Value(level)) }
}

func Time(time time.Time) Param {
	return func(l *logging.ServiceLogV1) { l.Time = datetime.DateTime(time) }
}

func TimeNow() Param {
	return func(l *logging.ServiceLogV1) { l.Time = datetime.DateTime(time.Now()) }
}

func Origin(origin string) Param {
	return func(l *logging.ServiceLogV1) { l.Origin = &origin }
}

// OriginFromInitLine sets the "origin" field to be the filename and line of the location at which this function is
// called.
func OriginFromInitLine() Param {
	origin := ""
	if file, line, ok := initLineCaller(1); ok {
		origin = file + ":" + strconv.Itoa(line)
	}
	return Origin(origin)
}

// OriginFromInitPkg sets the "origin" field to be the package path of the location at which this function is called.
// The skipPkg parameter determines the level of the parent package that should be used. For example, if this function
// is called in a file with the package path "github.com/palantir/witchcraft-go-logging/wlog", then with skipPkg=0 the
// origin would be "github.com/palantir/witchcraft-go-logging/wlog", while with skipPkg=1 the origin would be
// "github.com/palantir/witchcraft-go-logging".
func OriginFromInitPkg(skipPkg int) Param {
	return Origin(CallerPkg(1, skipPkg))
}

// OriginFromCallLine sets the "origin" field to be the filename and line of the location at which the logger invocation
// is performed.
//
// Note that, when this parameter is used, every log invocation will perform a "runtime.Caller" call, which may not be
// suitable for performance-critical scenarios.
//
// Note that this parameter is tied to the implementation details of the logger implementations defined in the svc1log
// package (it hard-codes assumptions relating to the number of call stacks that must be skipped to reach the log site).
// Using this parameter with an svc1log.Logger implementation not defined in the svc1log package may result in incorrect
// output. If wrapping the default implementation of svc1log.Logger, OriginFromCallLineWithSkip allows for trimming
// additional stack frames.
func OriginFromCallLine() Param {
	return OriginFromCallLineWithSkip(0)
}

const defaultOriginFromCallLineStackSkip = 8

// OriginFromCallLineWithSkip is like OriginFromCallLine but allows for configuring additional skipped stack frames.
// This allows for libraries wrapping loggers to hide their implementation frames from the caller.
func OriginFromCallLineWithSkip(skipFrames int) Param {
	return func(l *logging.ServiceLogV1) {
		origin := ""
		if file, line, ok := initLineCaller(defaultOriginFromCallLineStackSkip + skipFrames); ok {
			origin = file + ":" + strconv.Itoa(line)
		}
		l.Origin = &origin
	}
}

func Thread(thread string) Param {
	return func(l *logging.ServiceLogV1) { l.Thread = &thread }
}

func Message(message string) Param {
	return func(l *logging.ServiceLogV1) { l.Message = message }
}

func Params(params wparams.ParamStorer) Param {
	return func(l *logging.ServiceLogV1) {
		if params != nil {
			SafeParams(params.SafeParams())(l)
			UnsafeParams(params.UnsafeParams())(l)
		}
	}
}

func SafeParams(params map[string]any) Param {
	return func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = maps.Clone(params)
		} else {
			for k, v := range params {
				l.Params[k] = v
			}
		}
	}
}

func SafeParam(key string, value any) Param {
	return func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = map[string]any{key: value}
		} else {
			l.Params[key] = value
		}
	}
}

func UID(uid string) Param {
	return func(l *logging.ServiceLogV1) {
		l.Uid = (*logging.UserId)(&uid)
	}
}

func SID(sid string) Param {
	return func(l *logging.ServiceLogV1) {
		l.Sid = (*logging.SessionId)(&sid)
	}
}

func TokenID(tokenId string) Param {
	return func(l *logging.ServiceLogV1) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	}
}

func TraceID(traceId string) Param {
	return func(l *logging.ServiceLogV1) {
		l.TraceId = (*logging.TraceId)(&traceId)
	}
}

func Stacktrace(err error) Param {
	return func(l *logging.ServiceLogV1) {
		if err != nil {
			stacktrace := werror.GenerateErrorString(err, false)
			l.Stacktrace = &stacktrace

			// add all safe and unsafe parameters stored in error
			safeParams, unsafeParams := werror.ParamsFromError(err)
			SafeParams(safeParams)(l)
			UnsafeParams(unsafeParams)(l)
		}
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.ServiceLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = maps.Clone(unsafeParams)
		} else {
			for k, v := range unsafeParams {
				l.UnsafeParams[k] = v
			}
		}
	}
}

func UnsafeParam(key string, value any) Param {
	return func(l *logging.ServiceLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}

func Tags(tags map[string]string) Param {
	return func(l *logging.ServiceLogV1) {
		if l.Tags == nil {
			l.Tags = maps.Clone(tags)
		} else {
			for k, v := range tags {
				l.Tags[k] = v
			}
		}
	}
}

func Tag(key, value string) Param {
	return func(l *logging.ServiceLogV1) {
		if l.Tags == nil {
			l.Tags = map[string]string{key: value}
		} else {
			l.Tags[key] = value
		}
	}
}
