package wlog

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

// Logger is a generic logger that can log all Conjure log types.
type Logger[T logging.LogTypes] interface {
	Log(params ...Param[T])
}

// Param is a function that modifies a Conjure log object.
type Param[T logging.LogTypes] func(*T)

// NewDefaultLogger creates a new logger that writes JSON-marshaled lines to the provided output.
func NewDefaultLogger[T logging.LogTypes](output io.Writer, params ...Param[T]) Logger[T] {
	return NewDefaultLoggerWithPrinter(JSONPrinter[T](output), params...)
}

// NewDefaultLoggerWithPrinter creates a new logger that writes log objects using the provided printer.
func NewDefaultLoggerWithPrinter[T logging.LogTypes](printer LogPrinter[T], params ...Param[T]) Logger[T] {
	return &wrappedLogger[T]{
		logger: &defaultLogger[T]{
			printer: printer,
			params:  params,
		},
	}
}

type defaultLogger[T logging.LogTypes] struct {
	printer LogPrinter[T]
	params  []Param[T]
}

func (l *defaultLogger[T]) Log(params ...Param[T]) {
	pool := wloginternal.PoolFor[T]()
	obj := pool.Get()
	defer pool.Put(obj)

	for _, p := range l.params {
		if p != nil {
			p(obj)
		}
	}
	for _, p := range params {
		if p != nil {
			p(obj)
		}
	}
	if err := l.printer.Print(obj); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		// TODO: something else?
	}
}

// WithParams returns a new logger that logs with the provided parameters.
func WithParams[T logging.LogTypes](logger Logger[T], params ...Param[T]) Logger[T] {
	// Avoid deeply nested wrappedLogger instances by unwrapping if possible.
	if wl, ok := logger.(*wrappedLogger[T]); ok {
		logger = wl.logger
		params = append(slices.Clone(wl.params), params...)
	}
	return &wrappedLogger[T]{
		logger: logger,
		params: params,
	}
}

// wrappedLogger is a logger that wraps another logger and adds additional parameters to each log call.
type wrappedLogger[T logging.LogTypes] struct {
	params []Param[T]
	logger Logger[T]
}

func (l *wrappedLogger[T]) Log(params ...Param[T]) {
	l.logger.Log(append(l.params, params...)...)
}
