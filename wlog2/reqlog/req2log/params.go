package req2log

import (
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

type Param = wlog.ConjureLogParam[logging.RequestLogV2]

func Type() Param {
	return func(l *logging.RequestLogV2) {
		l.Type = "request.2"
	}
}

func Time(time time.Time) Param {
	return func(l *logging.RequestLogV2) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.RequestLogV2) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func Method(method string) Param {
	return func(l *logging.RequestLogV2) {
		l.Method = &method
	}
}

func Protocol(protocol string) Param {
	return func(l *logging.RequestLogV2) {
		l.Protocol = protocol
	}
}

func Path(path string) Param {
	return func(l *logging.RequestLogV2) {
		l.Path = path
	}
}

func SafeParams(params map[string]any) Param {
	return func(l *logging.RequestLogV2) {
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
	return func(l *logging.RequestLogV2) {
		if l.Params == nil {
			l.Params = map[string]any{key: value}
		} else {
			l.Params[key] = value
		}
	}
}

func Status(status int) Param {
	return func(l *logging.RequestLogV2) {
		l.Status = status
	}
}

func RequestSize(requestSize int64) Param {
	return func(l *logging.RequestLogV2) {
		l.RequestSize = safelong.SafeLong(requestSize)
	}
}

func ResponseSize(responseSize int64) Param {
	return func(l *logging.RequestLogV2) {
		l.ResponseSize = safelong.SafeLong(responseSize)
	}
}

func Duration(duration time.Duration) Param {
	return func(l *logging.RequestLogV2) {
		l.Duration = safelong.SafeLong(duration.Microseconds())
	}
}

func UID(uid string) Param {
	return func(l *logging.RequestLogV2) {
		l.Uid = (*logging.UserId)(&uid)
	}
}

func SID(sid string) Param {
	return func(l *logging.RequestLogV2) {
		l.Sid = (*logging.SessionId)(&sid)
	}
}

func TokenID(tokenId string) Param {
	return func(l *logging.RequestLogV2) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	}
}

func TraceID(traceId string) Param {
	return func(l *logging.RequestLogV2) {
		l.TraceId = (*logging.TraceId)(&traceId)
	}
}

func OrgID(orgId string) Param {
	return func(l *logging.RequestLogV2) {
		// TODO: Add OrgID to svc1log
		// l.OrgId = (*logging.OrgId)(&orgId)
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.RequestLogV2) {
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
	return func(l *logging.RequestLogV2) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}
