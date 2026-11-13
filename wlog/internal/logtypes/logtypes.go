// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

// Package logtypes defines the objects that typed loggers build and the built-in JSON provider encodes.
//
// Each type encodes to the witchcraft-logging-api wire format, which the schema test checks against the generated
// types. Fields that no parameter sets are left out. Optional values are omitted when empty, so an empty string
// means unset, and Reset clears an object in place so pooled objects keep their maps.
//
// Optional strings, pointers, structs, and maps use omitzero, which JSON v2 checks before it looks up a marshaler.
// Maps use AnyMap and StringMap so that an empty map Reset kept counts as zero. Slices use omitempty, which takes the
// same fast path because no custom marshaler applies to them, and also omits empty non-nil slices.
package logtypes

import (
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
)

// Log is a pointer to one of the log types. Only this package implements it.
type Log interface {
	isLog()
}

func (*ServiceLogV1) isLog()    {}
func (*EventLogV2) isLog()      {}
func (*MetricLogV1) isLog()     {}
func (*TraceLogV1) isLog()      {}
func (*RequestLogV2) isLog()    {}
func (*DiagnosticLogV1) isLog() {}
func (*AuditLogV2) isLog()      {}
func (*AuditLogV3) isLog()      {}
func (*WrappedLogV1) isLog()    {}

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type ServiceLogV1 struct {
	Type         string            `json:"type"`
	Level        LogLevel          `json:"level"`
	Time         datetime.DateTime `json:"time"`
	Origin       string            `json:"origin,omitzero"`
	Message      string            `json:"message,omitzero"`
	Params       AnyMap            `json:"params,omitzero"`
	Uid          string            `json:"uid,omitzero"`
	Sid          string            `json:"sid,omitzero"`
	TokenId      string            `json:"tokenId,omitzero"`
	OrgId        string            `json:"orgId,omitzero"`
	TraceId      string            `json:"traceId,omitzero"`
	Stacktrace   string            `json:"stacktrace,omitzero"`
	UnsafeParams AnyMap            `json:"unsafeParams,omitzero"`
	Tags         StringMap         `json:"tags,omitzero"`
}

func (l *ServiceLogV1) Reset() {
	*l = ServiceLogV1{Params: cleared(l.Params), UnsafeParams: cleared(l.UnsafeParams), Tags: cleared(l.Tags)}
}

type EventLogV2 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	EventName    string            `json:"eventName,omitzero"`
	Values       AnyMap            `json:"values,omitzero"`
	Uid          string            `json:"uid,omitzero"`
	Sid          string            `json:"sid,omitzero"`
	TokenId      string            `json:"tokenId,omitzero"`
	OrgId        string            `json:"orgId,omitzero"`
	TraceId      string            `json:"traceId,omitzero"`
	UnsafeParams AnyMap            `json:"unsafeParams,omitzero"`
	Tags         StringMap         `json:"tags,omitzero"`
}

func (l *EventLogV2) Reset() {
	*l = EventLogV2{Values: cleared(l.Values), UnsafeParams: cleared(l.UnsafeParams), Tags: cleared(l.Tags)}
}

type MetricLogV1 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	MetricName   string            `json:"metricName"`
	MetricType   string            `json:"metricType"`
	Values       AnyMap            `json:"values,omitzero"`
	Tags         StringMap         `json:"tags,omitzero"`
	Uid          string            `json:"uid,omitzero"`
	Sid          string            `json:"sid,omitzero"`
	TokenId      string            `json:"tokenId,omitzero"`
	OrgId        string            `json:"orgId,omitzero"`
	UnsafeParams AnyMap            `json:"unsafeParams,omitzero"`
}

func (l *MetricLogV1) Reset() {
	*l = MetricLogV1{Values: cleared(l.Values), Tags: cleared(l.Tags), UnsafeParams: cleared(l.UnsafeParams)}
}

type TraceLogV1 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	Uid          string            `json:"uid,omitzero"`
	Sid          string            `json:"sid,omitzero"`
	TokenId      string            `json:"tokenId,omitzero"`
	OrgId        string            `json:"orgId,omitzero"`
	UnsafeParams AnyMap            `json:"unsafeParams,omitzero"`
	Span         TraceLogV1Span    `json:"span"`
}

// Reset keeps UnsafeParams and Annotations but drops TraceLogV1Span.Tags, which borrows the caller's span map.
func (l *TraceLogV1) Reset() {
	*l = TraceLogV1{UnsafeParams: cleared(l.UnsafeParams), Span: TraceLogV1Span{Annotations: truncated(l.Span.Annotations)}}
}

type TraceLogV1Span struct {
	TraceId     string                     `json:"traceId"`
	Id          string                     `json:"id"`
	Name        string                     `json:"name"`
	ParentId    string                     `json:"parentId,omitzero"`
	Timestamp   int64                      `json:"timestamp"`
	Duration    int64                      `json:"duration"`
	Annotations []TraceLogV1SpanAnnotation `json:"annotations,omitempty"`
	Tags        StringMap                  `json:"tags,omitzero"`
}

type TraceLogV1SpanAnnotation struct {
	Timestamp int64                             `json:"timestamp"`
	Value     string                            `json:"value"`
	Endpoint  *TraceLogV1SpanAnnotationEndpoint `json:"endpoint,omitzero"`
}

type TraceLogV1SpanAnnotationEndpoint struct {
	ServiceName string `json:"serviceName"`
	Ipv4        string `json:"ipv4,omitzero"`
	Ipv6        string `json:"ipv6,omitzero"`
}

type RequestLogV2 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	Method       string            `json:"method,omitzero"`
	Protocol     string            `json:"protocol"`
	Path         string            `json:"path"`
	Params       AnyMap            `json:"params,omitzero"`
	Status       int               `json:"status"`
	RequestSize  int64             `json:"requestSize"`
	ResponseSize int64             `json:"responseSize"`
	Duration     int64             `json:"duration"`
	Uid          string            `json:"uid,omitzero"`
	Sid          string            `json:"sid,omitzero"`
	TokenId      string            `json:"tokenId,omitzero"`
	OrgId        string            `json:"orgId,omitzero"`
	TraceId      string            `json:"traceId,omitzero"`
	UnsafeParams AnyMap            `json:"unsafeParams,omitzero"`
}

func (l *RequestLogV2) Reset() {
	*l = RequestLogV2{Params: cleared(l.Params), UnsafeParams: cleared(l.UnsafeParams)}
}

type DiagnosticLogV1 struct {
	Type         string                    `json:"type"`
	Time         datetime.DateTime         `json:"time"`
	Diagnostic   DiagnosticLogV1Diagnostic `json:"diagnostic,omitzero"`
	UnsafeParams AnyMap                    `json:"unsafeParams,omitzero"`
}

func (l *DiagnosticLogV1) Reset() {
	*l = DiagnosticLogV1{UnsafeParams: cleared(l.UnsafeParams)}
}

// DiagnosticLogV1Diagnostic is a union: Type names the one non-nil variant, and an unknown variant is empty. The variants are the
// generated types because the diag1log API accepts them.
type DiagnosticLogV1Diagnostic struct {
	Type       string                     `json:"type"`
	Generic    *logging.GenericDiagnostic `json:"generic,omitzero"`
	ThreadDump *logging.ThreadDumpV1      `json:"threadDump,omitzero"`
}

type AuditResult string

type AuditLogV2 struct {
	Type          string            `json:"type"`
	Time          datetime.DateTime `json:"time"`
	Uid           string            `json:"uid,omitzero"`
	Sid           string            `json:"sid,omitzero"`
	TokenId       string            `json:"tokenId,omitzero"`
	OrgId         string            `json:"orgId,omitzero"`
	TraceId       string            `json:"traceId,omitzero"`
	OtherUids     []string          `json:"otherUids,omitempty"`
	Origin        string            `json:"origin,omitzero"`
	Name          string            `json:"name"`
	Result        AuditResult       `json:"result"`
	RequestParams AnyMap            `json:"requestParams,omitzero"`
	ResultParams  AnyMap            `json:"resultParams,omitzero"`
}

func (l *AuditLogV2) Reset() {
	*l = AuditLogV2{RequestParams: cleared(l.RequestParams), ResultParams: cleared(l.ResultParams)}
}

type AuditLogV3 struct {
	Type           string                         `json:"type"`
	Deployment     string                         `json:"deployment,omitzero"`
	Host           string                         `json:"host,omitzero"`
	Product        string                         `json:"product,omitzero"`
	ProductVersion string                         `json:"productVersion,omitzero"`
	Stack          string                         `json:"stack,omitzero"`
	Service        string                         `json:"service,omitzero"`
	Environment    string                         `json:"environment,omitzero"`
	ProducerType   string                         `json:"producerType,omitzero"`
	Organizations  []AuditLogV3Organization       `json:"organizations,omitempty"`
	EventId        string                         `json:"eventId,omitzero"`
	LogEntryId     string                         `json:"logEntryId"`
	UserAgent      string                         `json:"userAgent,omitzero"`
	Categories     []string                       `json:"categories,omitempty"`
	Entities       []any                          `json:"entities,omitempty"`
	Users          []AuditLogV3ContextualizedUser `json:"users,omitempty"`
	Origins        []string                       `json:"origins,omitempty"`
	SourceOrigin   string                         `json:"sourceOrigin,omitzero"`
	RequestFields  AnyMap                         `json:"requestFields,omitzero"`
	ResultFields   AnyMap                         `json:"resultFields,omitzero"`
	Time           datetime.DateTime              `json:"time"`
	Uid            string                         `json:"uid,omitzero"`
	Sid            string                         `json:"sid,omitzero"`
	TokenId        string                         `json:"tokenId,omitzero"`
	OrgId          string                         `json:"orgId,omitzero"`
	TraceId        string                         `json:"traceId,omitzero"`
	Origin         string                         `json:"origin,omitzero"`
	Name           string                         `json:"name"`
	Result         AuditResult                    `json:"result"`
}

// Reset keeps the maps and the slices, which parameters append to rather than assign.
func (l *AuditLogV3) Reset() {
	*l = AuditLogV3{
		Organizations: truncated(l.Organizations),
		Categories:    truncated(l.Categories),
		Entities:      truncated(l.Entities),
		Users:         truncated(l.Users),
		Origins:       truncated(l.Origins),
		RequestFields: cleared(l.RequestFields),
		ResultFields:  cleared(l.ResultFields),
	}
}

type AuditLogV3Organization struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

type AuditLogV3ContextualizedUser struct {
	Uid       string   `json:"uid"`
	UserName  *string  `json:"userName,omitzero"`
	FirstName *string  `json:"firstName,omitzero"`
	LastName  *string  `json:"lastName,omitzero"`
	Groups    []string `json:"groups,omitempty"`
	Realm     *string  `json:"realm,omitzero"`
}

type WrappedLogV1 struct {
	Type          string              `json:"type"`
	Payload       WrappedLogV1Payload `json:"payload"`
	EntityName    string              `json:"entityName"`
	EntityVersion string              `json:"entityVersion"`
}

// WrappedLogV1Payload is a union: Type names the one non-nil variant.
type WrappedLogV1Payload struct {
	Type            string           `json:"type"`
	ServiceLogV1    *ServiceLogV1    `json:"serviceLogV1,omitzero"`
	RequestLogV2    *RequestLogV2    `json:"requestLogV2,omitzero"`
	TraceLogV1      *TraceLogV1      `json:"traceLogV1,omitzero"`
	EventLogV2      *EventLogV2      `json:"eventLogV2,omitzero"`
	MetricLogV1     *MetricLogV1     `json:"metricLogV1,omitzero"`
	AuditLogV2      *AuditLogV2      `json:"auditLogV2,omitzero"`
	AuditLogV3      *AuditLogV3      `json:"auditLogV3,omitzero"`
	DiagnosticLogV1 *DiagnosticLogV1 `json:"diagnosticLogV1,omitzero"`
}

func cleared[M ~map[string]V, V any](m M) M {
	clear(m)
	return m
}

func truncated[S ~[]E, E any](s S) S {
	clear(s)
	return s[:0]
}
