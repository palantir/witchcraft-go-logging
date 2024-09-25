package logging

import (
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
)

// LogTypes is a constraint for generic types that combines all the Conjure log types.
type LogTypes interface {
	AuditLogV2 | DiagnosticLogV1 | EventLogV2 | MetricLogV1 | RequestLogV2 | ServiceLogV1 | TraceLogV1 | WrappedLogV1
}

// shared string aliases

type SessionId string
type TokenId string
type TraceId string
type UserId string

// audit.2

type AuditLogV2 struct {
	Type          string            `json:"type"`
	Time          datetime.DateTime `json:"time"`
	Name          string            `json:"name"`
	Result        AuditResult       `json:"result"`
	Uid           *UserId           `json:"uid,omitempty"`
	Sid           *SessionId        `json:"sid,omitempty"`
	TokenId       *TokenId          `json:"tokenId,omitempty"`
	TraceId       *TraceId          `json:"traceId,omitempty"`
	OtherUids     []UserId          `json:"otherUids,omitempty"`
	Origin        *string           `json:"origin,omitempty"`
	RequestParams map[string]any    `json:"requestParams,omitempty"`
	ResultParams  map[string]any    `json:"resultParams,omitempty"`
}

type AuditResult string

const (
	AuditResultSUCCESS      AuditResult = "SUCCESS"
	AuditResultUNAUTHORIZED AuditResult = "UNAUTHORIZED"
	AuditResultERROR        AuditResult = "ERROR"
)

// diagnostic.1

type DiagnosticLogV1 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	Diagnostic   Diagnostic        `json:"diagnostic"`
	UnsafeParams map[string]any    `json:"unsafeParams,omitempty"`
}

type Diagnostic struct {
	Type       string             `json:"type"`
	Generic    *GenericDiagnostic `json:"generic,omitempty"`
	ThreadDump *ThreadDumpV1      `json:"threadDump,omitempty"`
}

func NewDiagnosticFromGeneric(v GenericDiagnostic) Diagnostic {
	return Diagnostic{Type: "generic", Generic: &v}
}

func NewDiagnosticFromThreadDump(v ThreadDumpV1) Diagnostic {
	return Diagnostic{Type: "threadDump", ThreadDump: &v}
}

type GenericDiagnostic struct {
	DiagnosticType string `json:"diagnosticType"`
	Value          any    `json:"value"`
}

type ThreadDumpV1 struct {
	// Information about each of the threads in the thread dump. "Thread" may refer to a userland thread such as a goroutine, or an OS-level thread.
	Threads []ThreadInfoV1 `json:"threads,omitempty"`
}

type ThreadInfoV1 struct {
	// The ID of the thread.
	Id *safelong.SafeLong `json:"id,omitempty"`
	// The name of the thread. Note that thread names may include unsafe information such as the path of the HTTP request being processed. It must be safely redacted.
	Name *string `json:"name,omitempty"`
	// A list of stack frames for the thread, ordered with the current frame first.
	StackTrace []StackFrameV1 `json:"stackTrace,omitempty"`
	// Other thread-level information.
	Params map[string]any `json:"params,omitempty"`
}

type StackFrameV1 struct {
	// The address of the execution point of this stack frame. This is a string because a safelong can't represent the full 64 bit address space.
	Address *string `json:"address,omitempty"`
	// The identifier of the procedure containing the execution point of this stack frame. This is a fully qualified method name in Java and a demangled symbol name in native code, for example. Note that procedure names may include unsafe information if a service is, for exmaple, running user-defined code. It must be safely redacted.
	Procedure *string `json:"procedure,omitempty"`
	// The name of the file containing the source location of the execution point of this stack frame. Note that file names may include unsafe information if a service is, for example, running user-defined code. It must be safely redacted.
	File *string `json:"file,omitempty"`
	// The line number of the source location of the execution point of this stack frame.
	Line *int `json:"line,omitempty"`
	// Other frame-level information.
	Params map[string]any `json:"params,omitempty"`
}

// event.2

type EventLogV2 struct {
	Type string            `json:"type"`
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
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// Unsafe metadata describing the event
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
	// Additional dimensions that describe the instance of the log event
	Tags map[string]string `json:"tags,omitempty"`
}

// metric.1

type MetricLogV1 struct {
	Type string            `json:"type"`
	Time datetime.DateTime `json:"time"`
	// Dot-delimited name of metric, e.g. `com.foundry.compass.api.Compass.http.ping.failures`
	MetricName string `json:"metricName"`
	// Type of metric being represented, e.g. `gauge`, `histogram`, `counter`
	MetricType string `json:"metricType"`
	// Observations, measurements and context associated with the metric
	Values map[string]any `json:"values,omitempty"`
	// Additional dimensions that describe the instance of the metric
	Tags map[string]string `json:"tags,omitempty"`
	// User id (if available)
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// API token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Unsafe metadata describing the event
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
}

// request.2

type RequestLogV2 struct {
	Type string            `json:"type"`
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
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// API token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// Unredacted parameters such as path, query and header parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
}

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
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// Language-specific stack trace. Content is knowably safe. Renderers should substitute named placeholders ({name}, for name as a key) with keyed value from unsafeParams and leave non-matching keys as the original placeholder text.
	Stacktrace *string `json:"stacktrace,omitempty"`
	// Unredacted parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
	// Additional dimensions that describe the instance of the log event
	Tags map[string]string `json:"tags,omitempty"`
}

type LogLevel string

const (
	LogLevelDEBUG LogLevel = "DEBUG"
	LogLevelINFO  LogLevel = "INFO"
	LogLevelWARN  LogLevel = "WARN"
	LogLevelERROR LogLevel = "ERROR"
	LogLevelFATAL LogLevel = "FATAL"
)

// trace.1

type TraceLogV1 struct {
	Type         string            `json:"type"`
	Time         datetime.DateTime `json:"time"`
	Uid          *UserId           `json:"uid,omitempty"`
	Sid          *SessionId        `json:"sid,omitempty"`
	TokenId      *TokenId          `json:"tokenId,omitempty"`
	UnsafeParams map[string]any    `json:"unsafeParams,omitempty"`
	Span         Span              `json:"span"`
}

type Span struct {
	// 16-digit hex trace identifier
	TraceId string `json:"traceId"`
	// 16-digit hex span identifier
	Id string `json:"id"`
	// Name of the span (typically the operation/RPC/method name for corresponding to this span)
	Name string `json:"name"`
	// 16-digit hex identifer of the parent span
	ParentId *string `json:"parentId,omitempty"`
	// Timestamp of the start of this span (epoch microsecond value)
	Timestamp safelong.SafeLong `json:"timestamp"`
	// Duration of this span (microseconds)
	Duration    safelong.SafeLong `json:"duration"`
	Annotations []Annotation      `json:"annotations,omitempty"`
	// Additional dimensions that describe the instance of the trace span
	Tags map[string]string `json:"tags,omitempty"`
}

type Annotation struct {
	// Time annotation was created (epoch microsecond value)
	Timestamp safelong.SafeLong `json:"timestamp"`
	// Value encapsulated by this annotation
	Value    string   `json:"value"`
	Endpoint Endpoint `json:"endpoint"`
}

type Endpoint struct {
	// Name of the service that generated the annotation
	ServiceName string `json:"serviceName"`
	// IPv4 address of the machine that generated this annotation (`xxx.xxx.xxx.xxx`)
	Ipv4 *string `json:"ipv4,omitempty"`
	// IPv6 address of the machine that generated this annotation (standard hextet form)
	Ipv6 *string `json:"ipv6,omitempty"`
}

// wrapped.1

type WrappedLogV1 struct {
	// "wrapped.1"
	Type    string              `json:"type"`
	Payload WrappedLogV1Payload `json:"payload"`
	// Artifact part of entity's maven coordinate
	EntityName    string `json:"entityName"`
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

// Clear methods

func (log *AuditLogV2) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	log.OtherUids = nil
	log.Origin = nil
	log.Name = ""
	log.Result = ""
	clear(log.ResultParams)
	clear(log.RequestParams)
}

func (log *DiagnosticLogV1) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Diagnostic = Diagnostic{}
	clear(log.UnsafeParams)
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

func (log *MetricLogV1) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.MetricName = ""
	log.MetricType = ""
	clear(log.Values)
	clear(log.Tags)
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	clear(log.UnsafeParams)
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

func (log *TraceLogV1) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	clear(log.UnsafeParams)
	log.Span.TraceId = ""
	log.Span.Id = ""
	log.Span.Name = ""
	log.Span.ParentId = nil
	log.Span.Timestamp = 0
	log.Span.Duration = 0
	clear(log.Span.Annotations)
	clear(log.Span.Tags)
}

func (log *WrappedLogV1) Reset() {
	log.Type = ""
	log.Payload = WrappedLogV1Payload{}
	log.EntityName = ""
	log.EntityVersion = ""
}
