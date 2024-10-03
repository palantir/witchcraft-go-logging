package audit2log

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

type defaultLogger struct {
	logger wlog.Logger[logging.AuditLogV2]
}

func (l defaultLogger) Audit(name string, result AuditResultType, params ...Param) {
	wloginternal.LogParams(l.logger.Log, append([]Param{Type(), TimeNow(), Name(name), Result(result)}, params...)...)
}
