package audit2log

import (
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

type Param = wlog.ConjureLogParam[logging.AuditLogV2]

func Type() Param {
	return func(l *logging.AuditLogV2) {
		l.Type = "audit.2"
	}
}

func Time(time time.Time) Param {
	return func(l *logging.AuditLogV2) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.AuditLogV2) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func UID(uid string) Param {
	return func(l *logging.AuditLogV2) {
		l.Uid = (*logging.UserId)(&uid)
	}
}

func SID(sid string) Param {
	return func(l *logging.AuditLogV2) {
		l.Sid = (*logging.SessionId)(&sid)
	}
}

func TokenID(tokenId string) Param {
	return func(l *logging.AuditLogV2) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	}
}

func OrgID(orgId string) Param {
	return func(l *logging.AuditLogV2) {
		// TODO: Add OrgID to svc1log
		// l.OrgId = (*logging.OrgId)(&orgId)
	}
}

func TraceID(traceId string) Param {
	return func(l *logging.AuditLogV2) {
		l.TraceId = (*logging.TraceId)(&traceId)
	}
}

func OtherUIDs(uids ...string) Param {
	return func(l *logging.AuditLogV2) {
		for _, uid := range uids {
			l.OtherUids = append(l.OtherUids, logging.UserId(uid))
		}
	}
}

func Origin(origin string) Param {
	return func(l *logging.AuditLogV2) {
		l.Origin = &origin
	}
}

func Name(name string) Param {
	return func(l *logging.AuditLogV2) {
		l.Name = name
	}
}

func ResultSuccess() Param {
	return func(l *logging.AuditLogV2) {
		l.Result = logging.New_AuditResult(logging.AuditResult_SUCCESS)
	}
}

func ResultUnauthorized() Param {
	return func(l *logging.AuditLogV2) {
		l.Result = logging.New_AuditResult(logging.AuditResult_UNAUTHORIZED)
	}
}

func ResultError() Param {
	return func(l *logging.AuditLogV2) {
		l.Result = logging.New_AuditResult(logging.AuditResult_ERROR)
	}
}

func ResultParams(params map[string]any) Param {
	return func(l *logging.AuditLogV2) {
		if l.ResultParams == nil {
			l.ResultParams = maps.Clone(params)
		} else {
			for k, v := range params {
				l.ResultParams[k] = v
			}
		}
	}
}

func ResultParam(key string, value any) Param {
	return func(l *logging.AuditLogV2) {
		if l.ResultParams == nil {
			l.ResultParams = map[string]any{key: value}
		} else {
			l.ResultParams[key] = value
		}
	}
}

func RequestParams(params map[string]any) Param {
	return func(l *logging.AuditLogV2) {
		if l.RequestParams == nil {
			l.RequestParams = maps.Clone(params)
		} else {
			for k, v := range params {
				l.RequestParams[k] = v
			}
		}
	}
}

func RequestParam(key string, value any) Param {
	return func(l *logging.AuditLogV2) {
		if l.RequestParams == nil {
			l.RequestParams = map[string]any{key: value}
		} else {
			l.RequestParams[key] = value
		}
	}
}
