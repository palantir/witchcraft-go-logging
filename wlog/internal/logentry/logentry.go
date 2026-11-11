// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

// Package logentry writes log objects to a wlog.LogEntry for providers that do not accept them directly. Each XParams
// function writes only the fields that parameters set, so per-parameter callers can use it. The parity tests in each
// log package keep this output equal to the JSON encoding of the same object.
package logentry

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

// Logger consumes borrowed log objects synchronously. Implementations must copy any object or nested data they retain
// after the method returns.
type Logger interface {
	LogObject(logtypes.Log)
}

// NewLogger returns logger itself if it accepts log objects, and otherwise writes objects to its entries.
func NewLogger(logger wlog.Logger) Logger {
	if direct, ok := logger.(Logger); ok {
		return direct
	}
	return EntryLogger{Logger: logger}
}

// NewLeveledLogger is like NewLogger for loggers that log service logs, and wrapped service logs, at their level.
func NewLeveledLogger(logger wlog.LeveledLogger) Logger {
	if direct, ok := logger.(Logger); ok {
		return direct
	}
	return LeveledEntryLogger{Logger: logger}
}

type EntryLogger struct {
	Logger wlog.Logger
}

func (l EntryLogger) LogObject(log logtypes.Log) {
	l.Logger.Log(Param(log))
}

// LeveledEntryLogger drops logs without a known level, which only service logs and wrapped service logs have.
type LeveledEntryLogger struct {
	Logger wlog.LeveledLogger
}

func (l LeveledEntryLogger) LogObject(log logtypes.Log) {
	var level logtypes.LogLevel
	var message string
	switch log := log.(type) {
	case *logtypes.ServiceLogV1:
		level, message = log.Level, log.Message
	case *logtypes.WrappedLogV1:
		if svc := log.Payload.ServiceLogV1; svc != nil {
			level = svc.Level
		}
	}
	param := Param(log)
	switch level {
	case logtypes.LogLevelDebug:
		l.Logger.Debug(message, param)
	case logtypes.LogLevelInfo:
		l.Logger.Info(message, param)
	case logtypes.LogLevelWarn:
		l.Logger.Warn(message, param)
	case logtypes.LogLevelError:
		l.Logger.Error(message, param)
	}
}

// Param returns a parameter that writes log to an entry.
func Param(log logtypes.Log) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		Write(entry, log)
	})
}

// Write writes log to entry. A service log's message is left to the leveled logger's message argument.
func Write(entry wlog.LogEntry, log logtypes.Log) {
	switch log := log.(type) {
	case *logtypes.ServiceLogV1:
		ServiceLogV1(entry, log)
	case *logtypes.EventLogV2:
		EventLogV2(entry, log)
	case *logtypes.MetricLogV1:
		MetricLogV1(entry, log)
	case *logtypes.TraceLogV1:
		TraceLogV1(entry, log)
	case *logtypes.RequestLogV2:
		RequestLogV2(entry, log)
	case *logtypes.DiagnosticLogV1:
		DiagnosticLogV1(entry, log)
	case *logtypes.AuditLogV2:
		AuditLogV2(entry, log)
	case *logtypes.AuditLogV3:
		AuditLogV3(entry, log)
	case *logtypes.WrappedLogV1:
		WrappedLogV1(entry, log)
	}
}

// ServiceLogV1 writes every field except "message", which leveled providers take as the message argument.
func ServiceLogV1(entry wlog.LogEntry, log *logtypes.ServiceLogV1) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.StringValue("level", string(log.Level))
	ServiceLogV1Params(entry, log)
}

func ServiceLogV1Params(entry wlog.LogEntry, log *logtypes.ServiceLogV1) {
	optionalString(entry, "origin", log.Origin)
	anyMap(entry, "params", log.Params)
	optionalString(entry, "stacktrace", log.Stacktrace)
	optionalString(entry, wlog.UIDKey, log.Uid)
	optionalString(entry, wlog.SIDKey, log.Sid)
	optionalString(entry, wlog.TokenIDKey, log.TokenId)
	optionalString(entry, wlog.OrgIDKey, log.OrgId)
	optionalString(entry, wlog.TraceIDKey, log.TraceId)
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
	stringMap(entry, "tags", log.Tags)
}

func EventLogV2(entry wlog.LogEntry, log *logtypes.EventLogV2) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.OptionalStringValue("eventName", log.EventName)
	EventLogV2Params(entry, log)
}

func EventLogV2Params(entry wlog.LogEntry, log *logtypes.EventLogV2) {
	anyMap(entry, "values", log.Values)
	optionalString(entry, wlog.UIDKey, log.Uid)
	optionalString(entry, wlog.SIDKey, log.Sid)
	optionalString(entry, wlog.TokenIDKey, log.TokenId)
	optionalString(entry, wlog.OrgIDKey, log.OrgId)
	optionalString(entry, wlog.TraceIDKey, log.TraceId)
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
	stringMap(entry, "tags", log.Tags)
}

func MetricLogV1(entry wlog.LogEntry, log *logtypes.MetricLogV1) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.StringValue("metricName", log.MetricName)
	entry.StringValue("metricType", log.MetricType)
	MetricLogV1Params(entry, log)
}

// MetricLogV1Params does not write "time"; metric1log.ApplyParam writes it only for its time parameter.
func MetricLogV1Params(entry wlog.LogEntry, log *logtypes.MetricLogV1) {
	anyMap(entry, "values", log.Values)
	stringMap(entry, "tags", log.Tags)
	optionalString(entry, wlog.UIDKey, log.Uid)
	optionalString(entry, wlog.SIDKey, log.Sid)
	optionalString(entry, wlog.TokenIDKey, log.TokenId)
	optionalString(entry, wlog.OrgIDKey, log.OrgId)
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
}

func TraceLogV1(entry wlog.LogEntry, log *logtypes.TraceLogV1) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.AnyMapValue("span", spanValue(log.Span))
	TraceLogV1Params(entry, log)
}

func TraceLogV1Params(entry wlog.LogEntry, log *logtypes.TraceLogV1) {
	optionalString(entry, wlog.UIDKey, log.Uid)
	optionalString(entry, wlog.SIDKey, log.Sid)
	optionalString(entry, wlog.TokenIDKey, log.TokenId)
	optionalString(entry, wlog.OrgIDKey, log.OrgId)
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
}

// spanValue builds nested values directly to avoid allocating MapLogEntry bookkeeping for each object.
func spanValue(span logtypes.TraceLogV1Span) map[string]any {
	out := map[string]any{
		wlog.TraceIDKey: span.TraceId,
		"id":            span.Id,
		"name":          span.Name,
		"timestamp":     span.Timestamp,
		"duration":      span.Duration,
	}
	if span.ParentId != "" {
		out["parentId"] = span.ParentId
	}
	if len(span.Annotations) > 0 {
		annotations := make([]any, len(span.Annotations))
		for i, annotation := range span.Annotations {
			fields := map[string]any{
				"timestamp": annotation.Timestamp,
				"value":     annotation.Value,
			}
			if e := annotation.Endpoint; e != nil {
				endpoint := map[string]any{"serviceName": e.ServiceName}
				if e.Ipv4 != "" {
					endpoint["ipv4"] = e.Ipv4
				}
				if e.Ipv6 != "" {
					endpoint["ipv6"] = e.Ipv6
				}
				fields["endpoint"] = endpoint
			}
			annotations[i] = fields
		}
		out["annotations"] = annotations
	}
	if len(span.Tags) > 0 {
		out["tags"] = span.Tags
	}
	return out
}

func RequestLogV2(entry wlog.LogEntry, log *logtypes.RequestLogV2) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.OptionalStringValue("method", log.Method)
	entry.StringValue("protocol", log.Protocol)
	entry.StringValue("path", log.Path)
	anyMap(entry, "params", log.Params)
	entry.IntValue("status", int32(log.Status))
	entry.SafeLongValue("requestSize", log.RequestSize)
	entry.SafeLongValue("responseSize", log.ResponseSize)
	entry.SafeLongValue("duration", log.Duration)
	entry.OptionalStringValue(wlog.UIDKey, log.Uid)
	entry.OptionalStringValue(wlog.SIDKey, log.Sid)
	entry.OptionalStringValue(wlog.TokenIDKey, log.TokenId)
	entry.OptionalStringValue(wlog.OrgIDKey, log.OrgId)
	entry.OptionalStringValue(wlog.TraceIDKey, log.TraceId)
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
}

func DiagnosticLogV1(entry wlog.LogEntry, log *logtypes.DiagnosticLogV1) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	if diagnostic := diagnosticValue(log.Diagnostic); diagnostic != nil {
		entry.AnyMapValue("diagnostic", diagnostic)
	}
	DiagnosticLogV1Params(entry, log)
}

func DiagnosticLogV1Params(entry wlog.LogEntry, log *logtypes.DiagnosticLogV1) {
	anyMap(entry, wlog.UnsafeParamsKey, log.UnsafeParams)
}

func diagnosticValue(d logtypes.DiagnosticLogV1Diagnostic) map[string]any {
	switch {
	case d.Generic != nil:
		return map[string]any{
			"type": "generic",
			"generic": map[string]any{
				"diagnosticType": d.Generic.DiagnosticType,
				"value":          d.Generic.Value,
			},
		}
	case d.ThreadDump != nil:
		return map[string]any{
			"type": "threadDump",
			"threadDump": map[string]any{
				"threads": threadsField(d.ThreadDump.Threads),
			},
		}
	}
	return nil
}

func threadsField(threads []logging.ThreadInfoV1) []map[string]any {
	encodedThreads := make([]map[string]any, len(threads))
	for idx, threadInfo := range threads {
		encodedThreads[idx] = encodeThreadInfo(threadInfo)
	}
	return encodedThreads
}

func encodeThreadInfo(threadInfo logging.ThreadInfoV1) map[string]any {
	fields := make(map[string]any)
	if threadInfo.Id != nil {
		fields["id"] = threadInfo.Id
	}
	encodeNonEmptyString(fields, "name", threadInfo.Name)
	if len(threadInfo.StackTrace) > 0 {
		fields["stackTrace"] = stackTraceField(threadInfo.StackTrace)
	}
	if len(threadInfo.Params) > 0 {
		fields["params"] = threadInfo.Params
	}
	return fields
}

func stackTraceField(stackFrames []logging.StackFrameV1) []map[string]any {
	encodedStackTrace := make([]map[string]any, len(stackFrames))
	for idx, stackFrame := range stackFrames {
		encodedStackTrace[idx] = encodeStackFrame(stackFrame)
	}
	return encodedStackTrace
}

func encodeStackFrame(stackFrame logging.StackFrameV1) map[string]any {
	fields := make(map[string]any)
	encodeNonEmptyString(fields, "address", stackFrame.Address)
	encodeNonEmptyString(fields, "procedure", stackFrame.Procedure)
	encodeNonEmptyString(fields, "file", stackFrame.File)
	if stackFrame.Line != nil {
		fields["line"] = stackFrame.Line
	}
	if len(stackFrame.Params) > 0 {
		fields["params"] = stackFrame.Params
	}
	return fields
}

func encodeNonEmptyString(fields map[string]any, key string, val *string) {
	if val == nil || len(*val) == 0 {
		return
	}
	fields[key] = val
}

func AuditLogV2(entry wlog.LogEntry, log *logtypes.AuditLogV2) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.StringValue("name", log.Name)
	entry.StringValue("result", string(log.Result))
	AuditLogV2Params(entry, log)
}

func AuditLogV2Params(entry wlog.LogEntry, log *logtypes.AuditLogV2) {
	optionalString(entry, wlog.UIDKey, log.Uid)
	optionalString(entry, wlog.SIDKey, log.Sid)
	optionalString(entry, wlog.TokenIDKey, log.TokenId)
	optionalString(entry, wlog.OrgIDKey, log.OrgId)
	optionalString(entry, wlog.TraceIDKey, log.TraceId)
	if len(log.OtherUids) > 0 {
		entry.StringListValue("otherUids", log.OtherUids)
	}
	optionalString(entry, "origin", log.Origin)
	anyMap(entry, "requestParams", log.RequestParams)
	anyMap(entry, "resultParams", log.ResultParams)
}

func AuditLogV3(entry wlog.LogEntry, log *logtypes.AuditLogV3) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue(wlog.TimeKey, log.Time.String())
	entry.OptionalStringValue("deployment", log.Deployment)
	entry.OptionalStringValue("host", log.Host)
	entry.OptionalStringValue("product", log.Product)
	entry.OptionalStringValue("productVersion", log.ProductVersion)
	entry.OptionalStringValue("stack", log.Stack)
	entry.OptionalStringValue("service", log.Service)
	entry.OptionalStringValue("environment", log.Environment)
	entry.OptionalStringValue("producerType", log.ProducerType)
	if len(log.Organizations) > 0 {
		entry.ObjectListAppendValue("organizations", toAnySlice(log.Organizations))
	}
	entry.OptionalStringValue("eventId", log.EventId)
	entry.StringValue("logEntryId", log.LogEntryId)
	entry.OptionalStringValue("userAgent", log.UserAgent)
	if len(log.Categories) > 0 {
		entry.StringListAppendValue("categories", log.Categories)
	}
	if len(log.Entities) > 0 {
		entry.ObjectListAppendValue("entities", log.Entities)
	}
	if len(log.Users) > 0 {
		entry.ObjectListAppendValue("users", toAnySlice(log.Users))
	}
	if len(log.Origins) > 0 {
		entry.StringListAppendValue("origins", log.Origins)
	}
	entry.OptionalStringValue("sourceOrigin", log.SourceOrigin)
	anyMap(entry, "requestFields", log.RequestFields)
	anyMap(entry, "resultFields", log.ResultFields)
	entry.OptionalStringValue(wlog.UIDKey, log.Uid)
	entry.OptionalStringValue(wlog.SIDKey, log.Sid)
	entry.OptionalStringValue(wlog.TokenIDKey, log.TokenId)
	entry.OptionalStringValue(wlog.OrgIDKey, log.OrgId)
	entry.OptionalStringValue(wlog.TraceIDKey, log.TraceId)
	entry.OptionalStringValue("origin", log.Origin)
	entry.StringValue("name", log.Name)
	entry.StringValue("result", string(log.Result))
}

func toAnySlice[T any](values []T) []any {
	out := make([]any, len(values))
	for i, v := range values {
		out[i] = v
	}
	return out
}

func WrappedLogV1(entry wlog.LogEntry, log *logtypes.WrappedLogV1) {
	entry.StringValue(wlog.TypeKey, log.Type)
	entry.StringValue("entityName", log.EntityName)
	entry.StringValue("entityVersion", log.EntityVersion)
	p := log.Payload
	inner := wlog.NewMapLogEntry()
	switch {
	case p.ServiceLogV1 != nil:
		ServiceLogV1(inner, p.ServiceLogV1)
		inner.OptionalStringValue("message", p.ServiceLogV1.Message)
	case p.RequestLogV2 != nil:
		RequestLogV2(inner, p.RequestLogV2)
	case p.TraceLogV1 != nil:
		TraceLogV1(inner, p.TraceLogV1)
	case p.EventLogV2 != nil:
		EventLogV2(inner, p.EventLogV2)
	case p.MetricLogV1 != nil:
		MetricLogV1(inner, p.MetricLogV1)
	case p.AuditLogV2 != nil:
		AuditLogV2(inner, p.AuditLogV2)
	case p.AuditLogV3 != nil:
		AuditLogV3(inner, p.AuditLogV3)
	case p.DiagnosticLogV1 != nil:
		DiagnosticLogV1(inner, p.DiagnosticLogV1)
	}
	entry.AnyMapValue("payload", map[string]any{
		"type": p.Type,
		p.Type: inner.AllValues(),
	})
}

// optionalString skips empty values, so that applying one parameter to an entry leaves the fields it did not set.
func optionalString(entry wlog.LogEntry, key, value string) {
	if value != "" {
		entry.StringValue(key, value)
	}
}

func anyMap(entry wlog.LogEntry, key string, value logtypes.AnyMap) {
	if len(value) > 0 {
		entry.AnyMapValue(key, value)
	}
}

func stringMap(entry wlog.LogEntry, key string, value logtypes.StringMap) {
	if len(value) > 0 {
		entry.StringMapValue(key, value)
	}
}
