// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trc1log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

const (
	TypeValue = "trace.1"
)

type Param = wloginternal.Param[logging.TraceLogV1]

type paramFunc = wloginternal.ParamFunc[logging.TraceLogV1]

func Type() Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Type = TypeValue
	})
}

func Time(time time.Time) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Time = datetime.DateTime(time)
	})
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Time = datetime.DateTime(time.Now())
	})
}

func Span(span wtracing.SpanModel) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Span = logging.Span{
			TraceId:     string(span.TraceID),
			Id:          string(span.ID),
			Name:        span.Name,
			ParentId:    (*string)(span.ParentID),
			Timestamp:   safelong.SafeLong(span.Timestamp.Round(time.Microsecond).UnixNano() / 1e3),
			Duration:    safelong.SafeLong(span.Duration / time.Microsecond),
			Annotations: spanAnnotationsParam(span),
			Tags:        span.Tags,
		}
	})
}

func spanAnnotationsParam(span wtracing.SpanModel) []logging.Annotation {
	var startVal, endVal string
	switch span.Kind {
	case wtracing.Server:
		startVal, endVal = "sr", "ss"
	case wtracing.Client:
		startVal, endVal = "cs", "cr"
	default:
		return nil
	}
	return []logging.Annotation{
		{
			Timestamp: timestampMicros(span.Timestamp),
			Value:     startVal,
			Endpoint:  spanEndpoint(span.LocalEndpoint),
		},
		{
			Timestamp: timestampMicros(span.Timestamp.Add(span.Duration)),
			Value:     endVal,
			Endpoint:  spanEndpoint(span.LocalEndpoint),
		},
	}
}

func spanEndpoint(endpoint *wtracing.Endpoint) logging.Endpoint {
	e := logging.Endpoint{}
	if endpoint != nil {
		e.ServiceName = endpoint.ServiceName
		if endpoint.IPv4 != nil {
			s := endpoint.IPv4.String()
			e.Ipv4 = &s
		}
		if endpoint.IPv6 != nil {
			s := endpoint.IPv6.String()
			e.Ipv6 = &s
		}
	}
	return e
}

func timestampMicros(t time.Time) safelong.SafeLong {
	return safelong.SafeLong(t.Round(time.Microsecond).UnixNano() / 1e3)
}

func UID(uid string) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		l.OrgId = (*logging.OrgId)(&orgID)
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		wloginternal.SetMapParams(&l.UnsafeParams, unsafe)
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(l *logging.TraceLogV1) {
		wloginternal.SetMapParam(&l.UnsafeParams, key, value)
	})
}
