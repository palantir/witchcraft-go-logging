package logging

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
