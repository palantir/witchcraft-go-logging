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

package logging

import (
	"github.com/palantir/pkg/datetime"
)

type LogLevel string

const (
	LogLevelDEBUG LogLevel = "DEBUG"
	LogLevelINFO  LogLevel = "INFO"
	LogLevelWARN  LogLevel = "WARN"
	LogLevelERROR LogLevel = "ERROR"
	LogLevelFATAL LogLevel = "FATAL"
)

// service.1

type ServiceLogV1 struct {
	// "service.1"
	Type string `json:"type"`
	// The logger output level. One of {FATAL,ERROR,WARN,INFO,DEBUG,TRACE}.
	Level LogLevel `json:"level"`
	// RFC3339Nano UTC datetime string when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// Class or file name. May include line number.
	Origin *string `json:"origin,omitempty"`
	// Thread name
	Thread *string `json:"thread,omitempty"`
	// Log message. Palantir Java services using slf4j should not use slf4j placeholders ({}). Logs obtained from 3rd party libraries or services that use slf4j and contain slf4j placeholders will always produce `unsafeParams` with numeric indexes corresponding to the zero-indexed order of placeholders. Renderers should substitute numeric parameters from `unsafeParams` and may leave placeholders that do not match indexes as the original placeholder text.
	Message string `json:"message,omitempty"`
	// Known-safe parameters (redaction may be used to make params knowably safe, but is not required).
	Params map[string]any `json:"params,omitempty"`
	// User id (if available).
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// API token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Organization ID (if available)
	OrgId *OrgId `json:"orgId,omitempty"`
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// Language-specific stack trace. Content is knowably safe. Renderers should substitute named placeholders ({name}, for name as a key) with keyed value from unsafeParams and leave non-matching keys as the original placeholder text.
	Stacktrace *string `json:"stacktrace,omitempty"`
	// Unredacted parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
	// Additional dimensions that describe the instance of the log event
	Tags map[string]string `json:"tags,omitempty"`
}

func (log *ServiceLogV1) Reset() {
	log.Type = ""
	log.Level = ""
	log.Time = datetime.DateTime{}
	log.Origin = nil
	log.Thread = nil
	log.Message = ""
	clear(log.Params)
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	log.Stacktrace = nil
	clear(log.UnsafeParams)
	clear(log.Tags)
}
