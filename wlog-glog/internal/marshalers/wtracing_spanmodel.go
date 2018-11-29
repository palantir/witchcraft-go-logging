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
	"fmt"
	"strings"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

func marshalWTracingSpanModel(key string, val interface{}) string {
	builder := &strings.Builder{}
	needSeparator := false

	span := val.(wtracing.SpanModel)
	_, _ = builder.WriteString(fmt.Sprintf("%s: %s", wlog.TraceIDKey, span.TraceID))
	_, _ = builder.WriteString(separator)

	_, _ = builder.WriteString(fmt.Sprintf("%s: %s", trc1log.SpanIDKey, span.ID))
	_, _ = builder.WriteString(separator)

	_, _ = builder.WriteString(trc1log.SpanNameKey + ": " + span.Name)
	_, _ = builder.WriteString(" {")

	if parentID := span.ParentID; parentID != nil {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString(fmt.Sprintf("%s: %s", trc1log.SpanParentIDKey, *parentID))
		needSeparator = true
	}

	if needSeparator {
		_, _ = builder.WriteString(separator)
	}
	_, _ = builder.WriteString(fmt.Sprintf("%s: %d", trc1log.SpanTimestampKey, span.Timestamp.Round(time.Microsecond).UnixNano()/1e3))
	needSeparator = true

	if needSeparator {
		_, _ = builder.WriteString(separator)
	}
	_, _ = builder.WriteString(fmt.Sprintf("%s: %d", trc1log.SpanDurationKey, int64(span.Duration.Round(time.Microsecond))))
	needSeparator = true

	if kind := span.Kind; kind != "" {
		// if kind is non-empty, manually create v1-style annotations
		switch kind {
		case wtracing.Server:
			encodeSpanModelAnnotations(builder, &needSeparator, "sr", "ss", span)
		case wtracing.Client:
			encodeSpanModelAnnotations(builder, &needSeparator, "cs", "cr", span)
		}
	}

	_, _ = builder.WriteString("}")
	return builder.String()
}

func encodeSpanModelAnnotations(builder *strings.Builder, needSeparator *bool, startVal, endVal string, span wtracing.SpanModel) {
	if *needSeparator {
		_, _ = builder.WriteString(separator)
	}
	_, _ = builder.WriteString("[")

	spanModelAnnotationEncoder(builder, startVal, span.Timestamp, span.LocalEndpoint)
	_, _ = builder.WriteString(separator)
	spanModelAnnotationEncoder(builder, endVal, span.Timestamp.Add(span.Duration), span.LocalEndpoint)

	_, _ = builder.WriteString("]")
	*needSeparator = true
}

func spanModelAnnotationEncoder(builder *strings.Builder, value string, timeStamp time.Time, endpoint *wtracing.Endpoint) {
	_, _ = builder.WriteString(trc1log.AnnotationValueKey + ": " + value)
	_, _ = builder.WriteString(separator)

	_, _ = builder.WriteString(fmt.Sprintf("%s: %d", trc1log.AnnotationTimestampKey, timeStamp.Round(time.Microsecond).UnixNano()/1e3))
	_, _ = builder.WriteString(separator)

	_, _ = builder.WriteString(trc1log.AnnotationEndpointKey + ": {")
	_, _ = builder.WriteString(trc1log.EndpointServiceNameKey + ": " + endpoint.ServiceName)
	needSeparator := true

	if len(endpoint.IPv4) > 0 {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString(trc1log.EndpointIPv4Key + ": " + endpoint.IPv4.String())
		needSeparator = true
	}
	if len(endpoint.IPv6) > 0 {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString(trc1log.EndpointIPv6Key + ": " + endpoint.IPv6.String())
	}
	_, _ = builder.WriteString("}")
}
