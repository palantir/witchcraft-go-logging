package wlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/palantir/pkg/bytesbuffers"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
)

// ConjureLogType is a constraint for generic types that combines all the Conjure log types.
type ConjureLogType interface {
	logging.AuditLogV2 | logging.DiagnosticLogV1 | logging.EventLogV2 | logging.MetricLogV1 | logging.RequestLogV2 | logging.ServiceLogV1 | logging.TraceLogV1 | logging.WrappedLogV1
}

// ConjureLogger is a generic logger that can log all Conjure log types.
type ConjureLogger[T ConjureLogType] interface {
	Log(params ...ConjureLogParam[T])
	WithParams(params ...ConjureLogParam[T]) ConjureLogger[T]
}

// ConjureLogParam is a function that modifies a Conjure log object.
type ConjureLogParam[T ConjureLogType] func(*T)

// ConjureLogPrinter is a generic interface for printing Conjure log objects.
type ConjureLogPrinter[T ConjureLogType] interface {
	Print(log T) error
}

// NewDefaultLogger creates a new logger that writes JSON-marshaled lines to the provided output.
func NewDefaultLogger[T ConjureLogType](output io.Writer) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: jsonPrinter[T]{output: output},
	}
}

// NewDefaultLoggerWithPrinter creates a new logger that writes log objects using the provided printer.
func NewDefaultLoggerWithPrinter[T ConjureLogType](printer ConjureLogPrinter[T]) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: printer,
	}
}

type defaultLogger[T ConjureLogType] struct {
	params  []ConjureLogParam[T]
	printer ConjureLogPrinter[T]
}

func (l *defaultLogger[T]) WithParams(params ...ConjureLogParam[T]) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: l.printer,
		params:  append(append([]ConjureLogParam[T]{}, l.params...), params...)}
}

func (l *defaultLogger[T]) Log(params ...ConjureLogParam[T]) {
	// construct object
	var log T
	for _, p := range l.params {
		p(&log)
	}
	for _, p := range params {
		p(&log)
	}
	if err := l.printer.Print(log); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		// TODO: something else?
	}
}

// The size of buffers allocated to consume log data.
// This is a tradeoff between memory usage and the number of allocations.
// Lines longer than this limit will not have their buffers reused.
const bytesBufferPoolAllocSize = 4096

var bytesBufferPool = bytesbuffers.NewSyncPool(bytesBufferPoolAllocSize)

// jsonPrinter is a ConjureLogPrinter that writes JSON-marshaled log objects to an output.
// It minimizes allocations by using a shared buffer pool.
type jsonPrinter[T ConjureLogType] struct {
	output io.Writer
}

func (p jsonPrinter[T]) Print(log T) error {
	buf := bytesBufferPool.Get()
	defer bytesBufferPool.Put(buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(log); err != nil {
		return err
	}
	out := buf.Bytes()
	if out[len(out)-1] != '\n' {
		_ = buf.WriteByte('\n')
		out = buf.Bytes()
	}
	// write to output
	if _, err := p.output.Write(out); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		return err
	}
	return nil
}
