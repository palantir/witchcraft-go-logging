package wloginternal

import (
	"sync"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

func GetPooled[T logging.LogTypes]() (log *T, release func(*T)) {
	pool := poolFor[T]()
	return pool.Get(), pool.Put
}

type syncPool[T logging.LogTypes] struct {
	pool  sync.Pool
	reset func(*T)
}

func newPool[T logging.LogTypes](reset func(*T)) *syncPool[T] {
	return &syncPool[T]{pool: sync.Pool{New: func() any { return new(T) }}, reset: reset}
}

func (p *syncPool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *syncPool[T]) Put(obj *T) {
	p.reset(obj)
	p.pool.Put(obj)
}

var (
	poolAuditLogV2   = newPool((*logging.AuditLogV2).Reset)
	poolDiagLogV1    = newPool((*logging.DiagnosticLogV1).Reset)
	poolEventLogV2   = newPool((*logging.EventLogV2).Reset)
	poolMetricLogV1  = newPool((*logging.MetricLogV1).Reset)
	poolReqLogV2     = newPool((*logging.RequestLogV2).Reset)
	poolSvcLogV1     = newPool((*logging.ServiceLogV1).Reset)
	poolTraceLogV1   = newPool((*logging.TraceLogV1).Reset)
	poolWrappedLogV1 = newPool((*logging.WrappedLogV1).Reset)
)

// poolFor exposes a generic entrypoint to the object pools.
func poolFor[T logging.LogTypes]() *syncPool[T] {
	switch any(*new(T)).(type) {
	case logging.AuditLogV2:
		return any(poolAuditLogV2).(*syncPool[T])
	case logging.DiagnosticLogV1:
		return any(poolDiagLogV1).(*syncPool[T])
	case logging.EventLogV2:
		return any(poolEventLogV2).(*syncPool[T])
	case logging.MetricLogV1:
		return any(poolMetricLogV1).(*syncPool[T])
	case logging.RequestLogV2:
		return any(poolReqLogV2).(*syncPool[T])
	case logging.ServiceLogV1:
		return any(poolSvcLogV1).(*syncPool[T])
	case logging.TraceLogV1:
		return any(poolTraceLogV1).(*syncPool[T])
	case logging.WrappedLogV1:
		return any(poolWrappedLogV1).(*syncPool[T])
	default:
		panic("unhandled log type")
	}
}
