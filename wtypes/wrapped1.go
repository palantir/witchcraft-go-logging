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

// wrapped.1

type WrappedLogV1 struct {
	// "wrapped.1"
	Type string `json:"type"`
	// The log event
	Payload WrappedLogV1Payload `json:"payload"`
	// Artifact part of entity's maven coordinate
	EntityName string `json:"entityName"`
	// Version part of entity's maven coordinate
	EntityVersion string `json:"entityVersion"`
}

type WrappedLogV1Payload struct {
	Type            string           `json:"type"`
	ServiceLogV1    *ServiceLogV1    `json:"serviceLogV1,omitempty"`
	RequestLogV2    *RequestLogV2    `json:"requestLogV2,omitempty"`
	TraceLogV1      *TraceLogV1      `json:"traceLogV1,omitempty"`
	EventLogV2      *EventLogV2      `json:"eventLogV2,omitempty"`
	MetricLogV1     *MetricLogV1     `json:"metricLogV1,omitempty"`
	AuditLogV2      *AuditLogV2      `json:"auditLogV2,omitempty"`
	DiagnosticLogV1 *DiagnosticLogV1 `json:"diagnosticLogV1,omitempty"`
}

func NewWrappedLogV1PayloadFromServiceLogV1(v ServiceLogV1) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "serviceLogV1", ServiceLogV1: &v}
}

func NewWrappedLogV1PayloadFromRequestLogV2(v RequestLogV2) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "requestLogV2", RequestLogV2: &v}
}

func NewWrappedLogV1PayloadFromTraceLogV1(v TraceLogV1) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "traceLogV1", TraceLogV1: &v}
}

func NewWrappedLogV1PayloadFromEventLogV2(v EventLogV2) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "eventLogV2", EventLogV2: &v}
}

func NewWrappedLogV1PayloadFromMetricLogV1(v MetricLogV1) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "metricLogV1", MetricLogV1: &v}
}

func NewWrappedLogV1PayloadFromAuditLogV2(v AuditLogV2) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "auditLogV2", AuditLogV2: &v}
}

func NewWrappedLogV1PayloadFromDiagnosticLogV1(v DiagnosticLogV1) WrappedLogV1Payload {
	return WrappedLogV1Payload{Type: "diagnosticLogV1", DiagnosticLogV1: &v}
}

func (log *WrappedLogV1) Reset() {
	log.Type = ""
	log.Payload = WrappedLogV1Payload{}
	log.EntityName = ""
	log.EntityVersion = ""
}
