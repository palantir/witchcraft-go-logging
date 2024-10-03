package wlog

import (
	"fmt"
	"io"
	"os"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

type LoggerCreator[T logging.LogTypes] func(w io.Writer) Logger[T]

// Logger is a generic logger that can log all Conjure log types.
type Logger[T logging.LogTypes] interface {
	Log(*T)
}

// Param is a function that modifies a Conjure log object.
type Param[T logging.LogTypes] func(*T)

// NewDefaultLogger creates a new logger that writes JSON-marshaled lines to the provided output.
func NewDefaultLogger[T logging.LogTypes](output io.Writer) Logger[T] {
	return NewDefaultLoggerWithPrinter[T](GetDefaultPrinterCreator()(output))
}

// NewDefaultLoggerWithPrinter creates a new logger that writes log objects using the provided printer.
func NewDefaultLoggerWithPrinter[T logging.LogTypes](printer Printer) Logger[T] {
	return &defaultLogger[T]{
		printer: printer,
	}
}

type defaultLogger[T logging.LogTypes] struct {
	printer Printer
}

func (l *defaultLogger[T]) Log(log *T) {
	if err := l.printer.Print(any(*log).(logging.LogType)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		// TODO: something else?
	}
}
