package logging

import "sync"

var (
	poolAuditLogV2   = newSyncPool[*AuditLogV2]()
	poolDiagLogV1    = newSyncPool[*DiagnosticLogV1]()
	poolEventLogV2   = newSyncPool[*EventLogV2]()
	poolMetricLogV1  = newSyncPool[*MetricLogV1]()
	poolReqLogV2     = newSyncPool[*RequestLogV2]()
	poolSvcLogV1     = newSyncPool[*ServiceLogV1]()
	poolTraceLogV1   = newSyncPool[*TraceLogV1]()
	poolWrappedLogV1 = newSyncPool[*WrappedLogV1]()
)

func BorrowObject[T LogTypes](p SyncPool[T]) (obj *T, release func(*T)) {
	switch any(*new(T)).(type) {
	case AuditLogV2:
		return obj, release
	}
	return p.pool.Get().(*T)
}

type resetLog interface {
	Reset()
}

type SyncPool[T resetLog] struct {
	pool *sync.Pool
}

func newSyncPool[T resetLog]() SyncPool[T] {
	return SyncPool[T]{
		pool: &sync.Pool{
			New: func() any { return new(T) },
		},
	}
}

func (p SyncPool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p SyncPool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
