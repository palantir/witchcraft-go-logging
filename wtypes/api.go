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

// LogTypes is a constraint for generic types that combines all the Conjure log types.
// Preferred to LogType when an interface type is not required.
type LogTypes interface {
	AuditLogV2 | DiagnosticLogV1 | EventLogV2 | MetricLogV1 | RequestLogV2 | ServiceLogV1 | TraceLogV1 | WrappedLogV1
}

// LogType is an interface satisfied by all the witchcraft log types.
type LogType interface {
	logType() // marker method for the sealed interface.
}

var _, _, _, _, _, _, _, _ LogType = AuditLogV2{}, DiagnosticLogV1{}, EventLogV2{}, MetricLogV1{}, RequestLogV2{}, ServiceLogV1{}, TraceLogV1{}, WrappedLogV1{}

// shared string aliases

type OrgID string
type SessionID string
type TokenID string
type TraceID string
type UserID string
