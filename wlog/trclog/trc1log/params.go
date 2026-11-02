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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
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

type Param = wloginternal.Param[logging.TraceLogV1]

type paramFunc = wloginternal.ParamFunc[logging.TraceLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logging.TraceLogV1)
	wloginternal.ApplyParam(log, p)
	applyObjectParams(log, entry)
}

// populateLogObject pre-populates the type, time, and span fields of a TraceLogV1 from the span and applies user
// params in-place.
func populateLogObject(log *logging.TraceLogV1, span wtracing.SpanModel, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())

	log.Span.TraceId = string(span.TraceID)
	log.Span.Id = string(span.ID)
	log.Span.Name = span.Name
	log.Span.ParentId = (*string)(span.ParentID)
	log.Span.Timestamp = timestampMicros(span.Timestamp)
	log.Span.Duration = safelong.SafeLong(span.Duration / time.Microsecond)
	log.Span.Annotations = spanAnnotations(span)
	if len(span.Tags) > 0 {
		log.Span.Tags = span.Tags
	}
	wloginternal.ApplyParams(log, params)
	wloginternal.OmitEmpty(&log.Uid)
	wloginternal.OmitEmpty(&log.Sid)
	wloginternal.OmitEmpty(&log.TokenId)
	wloginternal.OmitEmpty(&log.OrgId)
}

func spanAnnotations(span wtracing.SpanModel) []logging.Annotation {
	if span.Kind != wtracing.Server && span.Kind != wtracing.Client {
		return nil
	}
	endpoint := spanEndpoint(span)
	startTime := timestampMicros(span.Timestamp)
	endTime := timestampMicros(span.Timestamp.Add(span.Duration))
	switch span.Kind {
	case wtracing.Server:
		return []logging.Annotation{
			{Timestamp: startTime, Value: "sr", Endpoint: endpoint},
			{Timestamp: endTime, Value: "ss", Endpoint: endpoint},
		}
	case wtracing.Client:
		return []logging.Annotation{
			{Timestamp: startTime, Value: "cs", Endpoint: endpoint},
			{Timestamp: endTime, Value: "cr", Endpoint: endpoint},
		}
	default:
		return nil
	}
}

func spanEndpoint(span wtracing.SpanModel) logging.Endpoint {
	endpoint := logging.Endpoint{}
	if span.LocalEndpoint != nil {
		endpoint.ServiceName = span.LocalEndpoint.ServiceName
		if len(span.LocalEndpoint.IPv4) > 0 {
			endpoint.Ipv4 = new(span.LocalEndpoint.IPv4.String())
		}
		if len(span.LocalEndpoint.IPv6) > 0 {
			endpoint.Ipv6 = new(span.LocalEndpoint.IPv6.String())
		}
	}
	return endpoint
}

func timestampMicros(t time.Time) safelong.SafeLong {
	return safelong.SafeLong(t.Round(time.Microsecond).UnixNano() / 1e3)
}

func objectToParam(log *logging.TraceLogV1, hasLocalEndpoint bool) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())

		// Build nested values directly to avoid allocating MapLogEntry bookkeeping for each object.
		spanValues := map[string]any{
			wlog.TraceIDKey:  log.Span.TraceId,
			SpanIDKey:        log.Span.Id,
			SpanNameKey:      log.Span.Name,
			SpanTimestampKey: int64(log.Span.Timestamp),
			SpanDurationKey:  int64(log.Span.Duration),
		}

		if parentID := log.Span.ParentId; parentID != nil {
			spanValues[SpanParentIDKey] = *parentID
		}
		if log.Span.Annotations != nil {
			annotations := make([]any, len(log.Span.Annotations))
			for i, annotation := range log.Span.Annotations {
				fields := map[string]any{
					AnnotationTimestampKey: int64(annotation.Timestamp),
					AnnotationValueKey:     annotation.Value,
				}
				// The required object field cannot distinguish a missing endpoint from an explicitly empty one.
				if hasLocalEndpoint {
					endpoint := map[string]any{
						EndpointServiceNameKey: annotation.Endpoint.ServiceName,
					}
					if annotation.Endpoint.Ipv4 != nil {
						endpoint[EndpointIPv4Key] = *annotation.Endpoint.Ipv4
					}
					if annotation.Endpoint.Ipv6 != nil {
						endpoint[EndpointIPv6Key] = *annotation.Endpoint.Ipv6
					}
					fields[AnnotationEndpointKey] = endpoint
				}
				annotations[i] = fields
			}
			spanValues[SpanAnnotationsKey] = annotations
		}
		if log.Span.Tags != nil {
			spanValues[SpanTagsKey] = log.Span.Tags
		}

		entry.AnyMapValue(SpanKey, spanValues)

		applyObjectParams(log, entry)
	})
}

func applyObjectParams(log *logging.TraceLogV1, entry wlog.LogEntry) {
	if log.Uid != nil {
		entry.OptionalStringValue(wlog.UIDKey, string(*log.Uid))
	}
	if log.Sid != nil {
		entry.OptionalStringValue(wlog.SIDKey, string(*log.Sid))
	}
	if log.TokenId != nil {
		entry.OptionalStringValue(wlog.TokenIDKey, string(*log.TokenId))
	}
	if log.OrgId != nil {
		entry.OptionalStringValue(wlog.OrgIDKey, string(*log.OrgId))
	}
	if len(log.UnsafeParams) > 0 {
		entry.AnyMapValue(wlog.UnsafeParamsKey, log.UnsafeParams)
	}
}
func UID(uid string) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		log.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		log.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		log.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		log.OrgId = (*logging.OrganizationId)(&orgID)
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logging.TraceLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
