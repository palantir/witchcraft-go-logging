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

package wtracing

import (
	"net/http"
	"strings"

	"github.com/palantir/witchcraft-go-error"
)

const (
	b3TraceID      = "X-B3-TraceId"
	b3SpanID       = "X-B3-SpanId"
	b3ParentSpanID = "X-B3-ParentSpanId"
	b3Sampled      = "X-B3-Sampled"
	b3Flags        = "X-B3-Flags"
)

// InjectB3HeaderVals takes the provided SpanContext and sets the appropriate B3 values in the provided header.
func InjectB3HeaderVals(req *http.Request, sc SpanContext) {
	if len(sc.TraceID) > 0 && len(sc.ID) > 0 {
		req.Header.Set(b3TraceID, string(sc.TraceID))
		req.Header.Set(b3SpanID, string(sc.ID))
		if parentID := sc.ParentID; parentID != nil {
			req.Header.Set(b3ParentSpanID, string(*sc.ParentID))
		}
	}

	if sc.Debug {
		req.Header.Set(b3Flags, "1")
	} else if sampled := sc.Sampled; sampled != nil {
		sampledVal := "0"
		if *sampled {
			sampledVal = "1"
		}
		req.Header.Set(b3Sampled, sampledVal)
	}
}

// ExtractB3HeaderVals returns a SpanContext created by extracting the B3 values from the provided header. Returns an
// error if the B3 values in the provided header are corrupt (for example, if only a TraceID or SpanID is specified or
// if a ParentID is specified when a TraceID or SpanID is not present). However, if both TraceID and SpanID are missing,
// other header are still extracted and set on the returned SpanContext. This means that it is possible for the returned
// SpanContext to have a blank TraceID and SpanID.
func ExtractB3HeaderVals(req *http.Request) (*SpanContext, error) {
	traceID := strings.ToLower(req.Header.Get(b3TraceID))
	spanID := strings.ToLower(req.Header.Get(b3SpanID))
	if (traceID == "") != (spanID == "") {
		// either both traceID and spanID must be present or neither must be present
		return nil, werror.Error("TraceID and SpanID must both be present or both be absent",
			werror.SafeParam("traceId", traceID),
			werror.SafeParam("spanId", spanID),
		)
	}

	var parentIDVal *SpanID
	if parentID := strings.ToLower(req.Header.Get(b3ParentSpanID)); parentID != "" {
		if traceID == "" || spanID == "" {
			return nil, werror.Error("ParentID was present but TraceID or SpanID was not",
				werror.SafeParam("parentId", parentID),
				werror.SafeParam("traceId", traceID),
				werror.SafeParam("spanId", spanID),
			)
		}
		parentIDVal = (*SpanID)(&parentID)
	}

	var sampledVal *bool
	switch sampledHeader := strings.ToLower(req.Header.Get(b3Sampled)); sampledHeader {
	case "0", "false":
		boolVal := false
		sampledVal = &boolVal
	case "1", "true":
		boolVal := true
		sampledVal = &boolVal
	case "":
		// keep nil
	default:
		return nil, werror.Error("Sampled value was invalid", werror.SafeParam("sampledHeaderVal", sampledHeader))
	}

	debug := req.Header.Get(b3Flags) == "1"
	if debug {
		sampledVal = nil
	}

	return &SpanContext{
		TraceID:  TraceID(traceID),
		ID:       SpanID(spanID),
		ParentID: parentIDVal,
		Debug:    debug,
		Sampled:  sampledVal,
	}, nil
}
