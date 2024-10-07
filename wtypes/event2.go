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
)

// event.2

type EventLogV2 struct {
	// "event.2"
	Type string `json:"type"`
	// RFC3339Nano timestamp when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// Dot-delimited name of event, e.g. `com.foundry.compass.api.Compass.http.ping.failures`
	EventName string `json:"eventName"`
	// Observations, measurements and context associated with the event
	Values map[string]any `json:"values,omitempty"`
	// User id (if available)
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// API token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Organization ID (if available)
	OrgId *OrgId `json:"orgId,omitempty"`
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// Unsafe metadata describing the event
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
	// Additional dimensions that describe the instance of the log event
	Tags map[string]string `json:"tags,omitempty"`
}

func (log *EventLogV2) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.EventName = ""
	clear(log.Values)
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	clear(log.UnsafeParams)
	clear(log.Tags)
}

func (EventLogV2) logType() {}
