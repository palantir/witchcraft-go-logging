// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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

package wtypes

import (
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
)

// request.2

type RequestLogV2 struct {
	// "request.2"
	Type string `json:"type"`
	// RFC3339Nano timestamp when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// HTTP method of request
	Method *string `json:"method,omitempty"`
	// Protocol, e.g. `HTTP/1.1`, `HTTP/2`
	Protocol string `json:"protocol,omitempty"`
	// Path of request. If templated, the unrendered path, e.g.: `/catalog/dataset/{datasetId}`, `/{rid}/paths/contents/{path:.*}`.
	Path string `json:"path,omitempty"`
	// Known-safe parameters
	Params map[string]any `json:"params,omitempty"`
	// HTTP status code of response
	Status int `json:"status,omitempty"`
	// Size of request (bytes)
	RequestSize safelong.SafeLong `json:"requestSize,omitempty"`
	// Size of response (bytes)
	ResponseSize safelong.SafeLong `json:"responseSize,omitempty"`
	// Amount of time spent handling request (microseconds)
	Duration safelong.SafeLong `json:"duration,omitempty"`
	// User id (if available)
	Uid *UserID `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionID `json:"sid,omitempty"`
	// API token id (if available)
	TokenId *TokenID `json:"tokenId,omitempty"`
	// Organization ID (if available)
	OrgId *OrgID `json:"orgId,omitempty"`
	// Zipkin trace id (if available)
	TraceId *TraceID `json:"traceId,omitempty"`
	// Unredacted parameters such as path, query and header parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
}

func (log *RequestLogV2) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Method = nil
	log.Protocol = ""
	log.Path = ""
	clear(log.Params)
	log.Status = 0
	log.RequestSize = 0
	log.ResponseSize = 0
	log.Duration = 0
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	clear(log.UnsafeParams)
}

func (RequestLogV2) logType() {}
