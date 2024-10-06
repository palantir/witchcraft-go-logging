package wloginternal

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

// LogObject logs a log object with the given parameters to the provided logger.
// The heap object is retrieved and released using the objectPool.
// The defaultParam is applied first, followed by the params.
func LogObject[T logging.LogTypes](
	logger wlog.Logger[T],
	objectPool *SyncPool[T],
	defaultParam Param[T],
	params ...Param[T],
) {
	log := objectPool.Get()
	ApplyParams(log, defaultParam)
	ApplyParams(log, params...)
	logger.Log(log)
	objectPool.Put(log)
}
