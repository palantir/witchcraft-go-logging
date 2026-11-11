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
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

const (
	TypeValue = "trace.1"

	SpanKey = "span"

	SpanIDKey          = "id"
	SpanNameKey        = "name"
	SpanParentIDKey    = "parentId"
	SpanTimestampKey   = "timestamp"
	SpanDurationKey    = "duration"
	SpanAnnotationsKey = "annotations"
	SpanTagsKey        = "tags"

	AnnotationTimestampKey = "timestamp"
	AnnotationValueKey     = "value"
	AnnotationEndpointKey  = "endpoint"

	EndpointServiceNameKey = "serviceName"
	EndpointIPv4Key        = "ipv4"
	EndpointIPv6Key        = "ipv6"
)

type Param = wloginternal.Param[logtypes.TraceLogV1]

type paramFunc = wloginternal.ParamFunc[logtypes.TraceLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logtypes.TraceLogV1)
	wloginternal.ApplyParam(log, p)
	logentry.TraceLogV1Params(entry, log)
}

// populateLogObject pre-populates the type, time, and span fields of a TraceLogV1 from the span and applies user
// params in-place.
func populateLogObject(log *logtypes.TraceLogV1, span wtracing.SpanModel, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())

	log.Span.TraceId = string(span.TraceID)
	log.Span.Id = string(span.ID)
	log.Span.Name = span.Name
	if span.ParentID != nil {
		log.Span.ParentId = string(*span.ParentID)
	}
	log.Span.Timestamp = timestampMicros(span.Timestamp)
	log.Span.Duration = int64(span.Duration / time.Microsecond)
	log.Span.Annotations = appendSpanAnnotations(log.Span.Annotations, span)
	if len(span.Tags) > 0 {
		log.Span.Tags = span.Tags
	}
	wloginternal.ApplyParams(log, params)
}

func appendSpanAnnotations(dst []logtypes.TraceLogV1SpanAnnotation, span wtracing.SpanModel) []logtypes.TraceLogV1SpanAnnotation {
	var start, end string
	switch span.Kind {
	case wtracing.Server:
		start, end = "sr", "ss"
	case wtracing.Client:
		start, end = "cs", "cr"
	default:
		return dst
	}
	endpoint := spanEndpoint(span)
	return append(dst,
		logtypes.TraceLogV1SpanAnnotation{Timestamp: timestampMicros(span.Timestamp), Value: start, Endpoint: endpoint},
		logtypes.TraceLogV1SpanAnnotation{Timestamp: timestampMicros(span.Timestamp.Add(span.Duration)), Value: end, Endpoint: endpoint},
	)
}

func spanEndpoint(span wtracing.SpanModel) *logtypes.TraceLogV1SpanAnnotationEndpoint {
	if span.LocalEndpoint == nil {
		return nil
	}
	endpoint := &logtypes.TraceLogV1SpanAnnotationEndpoint{ServiceName: span.LocalEndpoint.ServiceName}
	if len(span.LocalEndpoint.IPv4) > 0 {
		endpoint.Ipv4 = span.LocalEndpoint.IPv4.String()
	}
	if len(span.LocalEndpoint.IPv6) > 0 {
		endpoint.Ipv6 = span.LocalEndpoint.IPv6.String()
	}
	return endpoint
}

func timestampMicros(t time.Time) int64 {
	return t.Round(time.Microsecond).UnixNano() / 1e3
}

func UID(uid string) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		log.Uid = uid
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		log.Sid = sid
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		log.TokenId = tokenID
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		log.OrgId = orgID
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logtypes.TraceLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
