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
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wtypes"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

var objectPool = wloginternal.NewPool((*wtypes.TraceLogV1).Reset)

type defaultLogger struct {
	logger wlog.Logger[wtypes.TraceLogV1]
}

func (l *defaultLogger) Log(span wtracing.SpanModel, params ...Param) {
	wloginternal.LogObject(l.logger.Log, objectPool, defaultParam(span), params...)
}

func (l *defaultLogger) Send(span wtracing.SpanModel) {
	l.Log(span)
}

func (l *defaultLogger) Close() error {
	return nil
}

func defaultParam(span wtracing.SpanModel) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Type = TypeValue
		l.Time = datetime.DateTime(time.Now())

		l.Span.TraceId = string(span.TraceID)
		l.Span.Id = string(span.ID)
		l.Span.Name = span.Name
		l.Span.ParentId = (*string)(span.ParentID)
		l.Span.Timestamp = safelong.SafeLong(span.Timestamp.Round(time.Microsecond).UnixNano() / 1e3)
		l.Span.Duration = safelong.SafeLong(span.Duration / time.Microsecond)
		l.Span.Annotations = spanAnnotationsParam(span)
		l.Span.Tags = span.Tags
	})
}

func spanAnnotationsParam(span wtracing.SpanModel) []wtypes.Annotation {
	var startVal, endVal string
	switch span.Kind {
	case wtracing.Server:
		startVal, endVal = "sr", "ss"
	case wtracing.Client:
		startVal, endVal = "cs", "cr"
	default:
		return nil
	}
	return []wtypes.Annotation{
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

func spanEndpoint(endpoint *wtracing.Endpoint) wtypes.Endpoint {
	e := wtypes.Endpoint{}
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
