package logging

// LogTypes is a constraint for generic types that combines all the Conjure log types.
// Preferred to LogType when an interface type is not required.
type LogTypes interface {
	AuditLogV2 | DiagnosticLogV1 | EventLogV2 | MetricLogV1 | RequestLogV2 | ServiceLogV1 | TraceLogV1 | WrappedLogV1
}

// LogType is an interface satisfied by all the witchcraft log types.
type LogType interface {
	logType() // marker method for the sealed interface.
}

func (AuditLogV2) logType()      {}
func (DiagnosticLogV1) logType() {}
func (EventLogV2) logType()      {}
func (MetricLogV1) logType()     {}
func (RequestLogV2) logType()    {}
func (ServiceLogV1) logType()    {}
func (TraceLogV1) logType()      {}
func (WrappedLogV1) logType()    {}

var _, _, _, _, _, _, _, _ LogType = AuditLogV2{}, DiagnosticLogV1{}, EventLogV2{}, MetricLogV1{}, RequestLogV2{}, ServiceLogV1{}, TraceLogV1{}, WrappedLogV1{}

// shared string aliases

type OrgId string
type SessionId string
type TokenId string
type TraceId string
type UserId string
