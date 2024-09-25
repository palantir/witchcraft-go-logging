package diag1log

import (
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

type Param = wlog.ConjureLogParam[logging.DiagnosticLogV1]

func Type() Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Type = "diagnostic.1"
	}
}

func Time(time time.Time) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.DiagnosticLogV1) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func GenericDiagnostic(genericDiagnostic logging.GenericDiagnostic) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Diagnostic = logging.NewDiagnosticFromGeneric(genericDiagnostic)
	}
}

func ThreadDump(threadDumpV1 logging.ThreadDumpV1) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Diagnostic = logging.NewDiagnosticFromThreadDump(threadDumpV1)
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.DiagnosticLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = maps.Clone(unsafeParams)
		} else {
			for k, v := range unsafeParams {
				l.UnsafeParams[k] = v
			}
		}
	}
}

func UnsafeParam(key string, value any) Param {
	return func(l *logging.DiagnosticLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}
