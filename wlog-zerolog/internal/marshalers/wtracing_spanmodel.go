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

package marshalers

import (
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/rs/zerolog"
)

func marshalWTracingSpanModel(evt *zerolog.Event, key string, val interface{}) *zerolog.Event {
	return evt.Object(key, logObjectMarshalerFn(func(e *zerolog.Event) {
		span := val.(wtracing.SpanModel)
		e.Str(wlog.TraceIDKey, string(span.TraceID))
		e.Str(trc1log.SpanIDKey, string(span.ID))
		e.Str(trc1log.SpanNameKey, span.Name)

		if parentID := span.ParentID; parentID != nil {
			e.Str(trc1log.SpanParentIDKey, string(*parentID))
		}
		e.Int64(trc1log.SpanTimestampKey, span.Timestamp.Round(time.Microsecond).UnixNano()/1e3)
		e.Int64(trc1log.SpanDurationKey, int64(span.Duration/time.Microsecond))
		if kind := span.Kind; kind != "" {
			// if kind is non-empty, manually create v1-style annotations
			switch kind {
			case wtracing.Server:
				encodeSpanModelAnnotations(evt, "sr", "ss", span)
			case wtracing.Client:
				encodeSpanModelAnnotations(evt, "cs", "cr", span)
			}
		}
		if tags := span.Tags; tags != nil && len(tags) > 0 {
			e.Interface(trc1log.SpanTagsKey, tags)
		}
	}))
}

func encodeSpanModelAnnotations(evt *zerolog.Event, startVal, endVal string, span wtracing.SpanModel) {
	evt.Array(trc1log.SpanAnnotationsKey, logArrayMarshalerFn(func(a *zerolog.Array) {
		// add "sr" annotation
		a.Object(spanModelAnnotationEncoder(startVal, span.Timestamp, span.LocalEndpoint))
		// add "ss" annotation
		a.Object(spanModelAnnotationEncoder(endVal, span.Timestamp.Add(span.Duration), span.LocalEndpoint))
	}))
}

func spanModelAnnotationEncoder(value string, timeStamp time.Time, endpoint *wtracing.Endpoint) logObjectMarshalerFn {
	return logObjectMarshalerFn(func(e *zerolog.Event) {
		e.Str(trc1log.AnnotationValueKey, value)
		e.Int64(trc1log.AnnotationTimestampKey, timeStamp.Round(time.Microsecond).UnixNano()/1e3)
		e.Object(trc1log.AnnotationEndpointKey, logObjectMarshalerFn(func(e *zerolog.Event) {
			e.Str(trc1log.EndpointServiceNameKey, endpoint.ServiceName)
			if len(endpoint.IPv4) > 0 {
				e.Str(trc1log.EndpointIPv4Key, endpoint.IPv4.String())
			}
			if len(endpoint.IPv6) > 0 {
				e.Str(trc1log.EndpointIPv6Key, endpoint.IPv6.String())
			}
		}))
	})
}
