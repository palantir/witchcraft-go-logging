package logging

import (
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
)

// trace.1

type TraceLogV1 struct {
	// "trace.1"
	Type string `json:"type"`
	// RFC3339Nano timestamp when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// User id (if available)
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// Token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Organization ID (if available)
	OrgId *OrgId `json:"orgId,omitempty"`
	// Unredacted parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
	// Span information
	Span Span `json:"span"`
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
	Duration safelong.SafeLong `json:"duration"`
	// Annotations that describe the instance of the trace span
	Annotations []Annotation `json:"annotations,omitempty"`
	// Additional dimensions that describe the instance of the trace span
	Tags map[string]string `json:"tags,omitempty"`
}

type Annotation struct {
	// Time annotation was created (epoch microsecond value)
	Timestamp safelong.SafeLong `json:"timestamp"`
	// Value encapsulated by this annotation
	Value string `json:"value"`
	// Endpoint that generated this annotation
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
