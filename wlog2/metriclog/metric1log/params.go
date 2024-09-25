package metric1log

import (
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

type Param = wlog.ConjureLogParam[logging.MetricLogV1]

func Type() Param {
	return func(l *logging.MetricLogV1) {
		l.Type = "metric.1"
	}
}

func Time(time time.Time) Param {
	return func(l *logging.MetricLogV1) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.MetricLogV1) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func MetricName(name string) Param {
	return func(l *logging.MetricLogV1) {
		l.MetricName = name
	}
}

func MetricType(metric string) Param {
	return func(l *logging.MetricLogV1) {
		l.MetricType = metric
	}
}

func Values(params map[string]any) Param {
	return func(l *logging.MetricLogV1) {
		if l.Values == nil {
			l.Values = maps.Clone(params)
		} else {
			for k, v := range params {
				l.Values[k] = v
			}
		}
	}
}

func Value(key string, value any) Param {
	return func(l *logging.MetricLogV1) {
		if l.Values == nil {
			l.Values = map[string]any{key: value}
		} else {
			l.Values[key] = value
		}
	}
}

func Tags(tags map[string]string) Param {
	return func(l *logging.MetricLogV1) {
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
	return func(l *logging.MetricLogV1) {
		if l.Tags == nil {
			l.Tags = map[string]string{key: value}
		} else {
			l.Tags[key] = value
		}
	}
}

func UID(uid string) Param {
	return func(l *logging.MetricLogV1) {
		l.Uid = (*logging.UserId)(&uid)
	}
}

func SID(sid string) Param {
	return func(l *logging.MetricLogV1) {
		l.Sid = (*logging.SessionId)(&sid)
	}
}

func TokenID(tokenId string) Param {
	return func(l *logging.MetricLogV1) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	}
}

func OrgID(orgId string) Param {
	return func(l *logging.MetricLogV1) {
		// TODO: Add OrgID to svc1log
		// l.OrgId = (*logging.OrgId)(&orgId)
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.MetricLogV1) {
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
	return func(l *logging.MetricLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}
