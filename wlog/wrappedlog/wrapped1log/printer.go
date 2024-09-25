package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog"
)

type wrappedPrinter[T logging.LogTypes] struct {
	delegate   wlog.Logger[logging.WrappedLogV1]
	newPayload func(payload T) logging.WrappedLogV1Payload
}

func (p wrappedPrinter[T]) Print(log T) error {
	p.delegate.Log(Payload(p.newPayload(log)))
	return nil
}

func audit2Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.AuditLogV2] {
	return wrappedPrinter[logging.AuditLogV2]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromAuditLogV2,
	}
}

func diag1Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.DiagnosticLogV1] {
	return wrappedPrinter[logging.DiagnosticLogV1]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromDiagnosticLogV1,
	}
}

func evt2Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.EventLogV2] {
	return wrappedPrinter[logging.EventLogV2]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromEventLogV2,
	}
}

func metric1Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.MetricLogV1] {
	return wrappedPrinter[logging.MetricLogV1]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromMetricLogV1,
	}
}

func req2Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.RequestLogV2] {
	return wrappedPrinter[logging.RequestLogV2]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromRequestLogV2,
	}
}

func svc1Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.ServiceLogV1] {
	return wrappedPrinter[logging.ServiceLogV1]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromServiceLogV1,
	}
}

func trc1Printer(delegate wlog.Logger[logging.WrappedLogV1]) wlog.LogPrinter[logging.TraceLogV1] {
	return wrappedPrinter[logging.TraceLogV1]{
		delegate:   delegate,
		newPayload: logging.NewWrappedLogV1PayloadFromTraceLogV1,
	}
}
