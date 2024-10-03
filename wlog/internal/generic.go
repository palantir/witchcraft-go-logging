package wloginternal

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

type LogTypeMapper[OUT any] interface {
	VisitAuditLogV2(logging.AuditLogV2) OUT
	VisitDiagnosticLogV1(logging.DiagnosticLogV1) OUT
	VisitEventLogV2(logging.EventLogV2) OUT
	VisitMetricLogV1(logging.MetricLogV1) OUT
	VisitRequestLogV2(logging.RequestLogV2) OUT
	VisitServiceLogV1(logging.ServiceLogV1) OUT
	VisitTraceLogV1(logging.TraceLogV1) OUT
	VisitWrappedLogV1(logging.WrappedLogV1) OUT
}

type LogTypeVisitor interface {
	VisitAuditLogV2(logging.AuditLogV2)
	VisitDiagnosticLogV1(logging.DiagnosticLogV1)
	VisitEventLogV2(logging.EventLogV2)
	VisitMetricLogV1(logging.MetricLogV1)
	VisitRequestLogV2(logging.RequestLogV2)
	VisitServiceLogV1(logging.ServiceLogV1)
	VisitTraceLogV1(logging.TraceLogV1)
	VisitWrappedLogV1(logging.WrappedLogV1)
}

type GenericLogTypeVisitor func(log logging.LogType)

func (v GenericLogTypeVisitor) VisitAuditLogV2(log logging.AuditLogV2)           { v(log) }
func (v GenericLogTypeVisitor) VisitDiagnosticLogV1(log logging.DiagnosticLogV1) { v(log) }
func (v GenericLogTypeVisitor) VisitEventLogV2(log logging.EventLogV2)           { v(log) }
func (v GenericLogTypeVisitor) VisitMetricLogV1(log logging.MetricLogV1)         { v(log) }
func (v GenericLogTypeVisitor) VisitRequestLogV2(log logging.RequestLogV2)       { v(log) }
func (v GenericLogTypeVisitor) VisitServiceLogV1(log logging.ServiceLogV1)       { v(log) }
func (v GenericLogTypeVisitor) VisitTraceLogV1(log logging.TraceLogV1)           { v(log) }
func (v GenericLogTypeVisitor) VisitWrappedLogV1(log logging.WrappedLogV1)       { v(log) }

func VisitLogType[T logging.LogTypes](log T, v LogTypeVisitor) {
	switch t := any(log).(type) {
	case logging.AuditLogV2:
		v.VisitAuditLogV2(t)
	case logging.DiagnosticLogV1:
		v.VisitDiagnosticLogV1(t)
	case logging.EventLogV2:
		v.VisitEventLogV2(t)
	case logging.MetricLogV1:
		v.VisitMetricLogV1(t)
	case logging.RequestLogV2:
		v.VisitRequestLogV2(t)
	case logging.ServiceLogV1:
		v.VisitServiceLogV1(t)
	case logging.TraceLogV1:
		v.VisitTraceLogV1(t)
	case logging.WrappedLogV1:
		v.VisitWrappedLogV1(t)
	default:
		panic("unreachable")
	}
}

type GenericLogTypeMapper[T any] func(log logging.LogType) T

func (m GenericLogTypeMapper[T]) VisitAuditLogV2(log logging.AuditLogV2) T           { return m(log) }
func (m GenericLogTypeMapper[T]) VisitDiagnosticLogV1(log logging.DiagnosticLogV1) T { return m(log) }
func (m GenericLogTypeMapper[T]) VisitEventLogV2(log logging.EventLogV2) T           { return m(log) }
func (m GenericLogTypeMapper[T]) VisitMetricLogV1(log logging.MetricLogV1) T         { return m(log) }
func (m GenericLogTypeMapper[T]) VisitRequestLogV2(log logging.RequestLogV2) T       { return m(log) }
func (m GenericLogTypeMapper[T]) VisitServiceLogV1(log logging.ServiceLogV1) T       { return m(log) }
func (m GenericLogTypeMapper[T]) VisitTraceLogV1(log logging.TraceLogV1) T           { return m(log) }
func (m GenericLogTypeMapper[T]) VisitWrappedLogV1(log logging.WrappedLogV1) T       { return m(log) }

func MapLogType[T logging.LogTypes, OUT any](log T, m LogTypeMapper[OUT]) OUT {
	switch t := any(log).(type) {
	case logging.AuditLogV2:
		return m.VisitAuditLogV2(t)
	case logging.DiagnosticLogV1:
		return m.VisitDiagnosticLogV1(t)
	case logging.EventLogV2:
		return m.VisitEventLogV2(t)
	case logging.MetricLogV1:
		return m.VisitMetricLogV1(t)
	case logging.RequestLogV2:
		return m.VisitRequestLogV2(t)
	case logging.ServiceLogV1:
		return m.VisitServiceLogV1(t)
	case logging.TraceLogV1:
		return m.VisitTraceLogV1(t)
	case logging.WrappedLogV1:
		return m.VisitWrappedLogV1(t)
	default:
		panic("unreachable")
	}
}
