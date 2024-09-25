package wloginternal

import (
	"sync"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

type resettable interface {
	// Reset resets the object to its zero state so that it can be reused.
	Reset()
}

var (
	poolAuditLogV2   SyncPool[logging.AuditLogV2]      = new(syncPool[logging.AuditLogV2])
	poolDiagLogV1    SyncPool[logging.DiagnosticLogV1] = new(syncPool[logging.DiagnosticLogV1])
	poolEventLogV2   SyncPool[logging.EventLogV2]      = new(syncPool[logging.EventLogV2])
	poolMetricLogV1  SyncPool[logging.MetricLogV1]     = new(syncPool[logging.MetricLogV1])
	poolReqLogV2     SyncPool[logging.RequestLogV2]    = new(syncPool[logging.RequestLogV2])
	poolSvcLogV1     SyncPool[logging.ServiceLogV1]    = new(syncPool[logging.ServiceLogV1])
	poolTraceLogV1   SyncPool[logging.TraceLogV1]      = new(syncPool[logging.TraceLogV1])
	poolWrappedLogV1 SyncPool[logging.WrappedLogV1]    = new(syncPool[logging.WrappedLogV1])
)

type SyncPool[T logging.LogTypes | []byte] interface {
	Get() *T
	Put(*T)
}

// PoolFor returns a SyncPool for the provided log type.
func PoolFor[T logging.LogTypes]() SyncPool[T] {
	switch any(*new(T)).(type) {
	case logging.AuditLogV2:
		return any(poolAuditLogV2).(SyncPool[T])
	case logging.DiagnosticLogV1:
		return any(poolDiagLogV1).(SyncPool[T])
	case logging.EventLogV2:
		return any(poolEventLogV2).(SyncPool[T])
	case logging.MetricLogV1:
		return any(poolMetricLogV1).(SyncPool[T])
	case logging.RequestLogV2:
		return any(poolReqLogV2).(SyncPool[T])
	case logging.ServiceLogV1:
		return any(poolSvcLogV1).(SyncPool[T])
	case logging.TraceLogV1:
		return any(poolTraceLogV1).(SyncPool[T])
	case logging.WrappedLogV1:
		return any(poolWrappedLogV1).(SyncPool[T])
	default:
		// This should never happen, but just construct a new pool for short term use.
		return new(syncPool[T])
	}
}

type syncPool[T logging.LogTypes | []byte] struct {
	pool sync.Pool
}

func (p *syncPool[T]) Get() *T {
	if v := p.pool.Get(); v != nil {
		return v.(*T)
	}
	return new(T)
}

func (p *syncPool[T]) Put(obj *T) {
	if r, ok := any(obj).(resettable); ok {
		r.Reset()
		p.pool.Put(obj)
	}
}
