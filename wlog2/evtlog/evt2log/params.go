package evt2log

import (
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
	wparams "github.com/palantir/witchcraft-go-params"
)

type Param = wlog.ConjureLogParam[logging.EventLogV2]

func Type() Param {
	return func(l *logging.EventLogV2) {
		l.Type = "event.2"
	}
}

func Time(time time.Time) Param {
	return func(l *logging.EventLogV2) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.EventLogV2) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func EventName(event string) Param {
	return func(l *logging.EventLogV2) {
		l.EventName = event
	}
}

func Params(params wparams.ParamStorer) Param {
	return func(l *logging.EventLogV2) {
		if params != nil {
			SafeParams(params.SafeParams())(l)
			UnsafeParams(params.UnsafeParams())(l)
		}
	}
}

func SafeParams(params map[string]any) Param {
	return func(l *logging.EventLogV2) {
		if l.Values == nil {
			l.Values = maps.Clone(params)
		} else {
			for k, v := range params {
				l.Values[k] = v
			}
		}
	}
}

func SafeParam(key string, value any) Param {
	return func(l *logging.EventLogV2) {
		if l.Values == nil {
			l.Values = map[string]any{key: value}
		} else {
			l.Values[key] = value
		}
	}
}

func UID(uid string) Param {
	return func(l *logging.EventLogV2) {
		l.Uid = (*logging.UserId)(&uid)
	}
}

func SID(sid string) Param {
	return func(l *logging.EventLogV2) {
		l.Sid = (*logging.SessionId)(&sid)
	}
}

func TokenID(tokenId string) Param {
	return func(l *logging.EventLogV2) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	}
}

func OrgID(orgId string) Param {
	return func(l *logging.EventLogV2) {
		// TODO: Add OrgID to svc1log
		// l.OrgId = (*logging.OrgId)(&orgId)
	}
}

func TraceID(traceId string) Param {
	return func(l *logging.EventLogV2) {
		l.TraceId = (*logging.TraceId)(&traceId)
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.EventLogV2) {
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
	return func(l *logging.EventLogV2) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}

func Tags(tags map[string]string) Param {
	return func(l *logging.EventLogV2) {
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
	return func(l *logging.EventLogV2) {
		if l.Tags == nil {
			l.Tags = map[string]string{key: value}
		} else {
			l.Tags[key] = value
		}
	}
}
