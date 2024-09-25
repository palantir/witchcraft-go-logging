package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

// wrappedPrinter implements Printer for logs included in the wrapped.1 payload field.
// When an underlying log object is constructed and passed to the Print method,
// the delegate WrappedLogV1 logger is called with a new WrappedLogV1Payload object.
type wrappedPrinter[T logging.LogTypes] struct {
	delegate   wlog.Logger[logging.WrappedLogV1]
	newPayload func(payload T) logging.WrappedLogV1Payload
}

func wrapPrinter[T logging.LogTypes](
	delegate wlog.Logger[logging.WrappedLogV1],
	newPayload func(payload T) logging.WrappedLogV1Payload,
) wlog.LogPrinter[T] {
	return wrappedPrinter[T]{
		delegate:   delegate,
		newPayload: newPayload,
	}
}

func (p wrappedPrinter[T]) Print(log *T) error {
	p.delegate.Log(Payload(p.newPayload(*log)))
	return nil
}
