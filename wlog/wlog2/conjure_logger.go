package wlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/palantir/pkg/bytesbuffers"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
)

var (
	bytesBufferPool = bytesbuffers.NewSyncPool(4096)
)

type ConjureLogType interface {
	logging.AuditLogV2 | logging.DiagnosticLogV1 | logging.EventLogV2 | logging.MetricLogV1 | logging.RequestLogV2 | logging.ServiceLogV1 | logging.TraceLogV1 | logging.WrappedLogV1
}

type ConjureLogger[T ConjureLogType] interface {
	Log(params ...ConjureLogParam[T])
	WithParams(params ...ConjureLogParam[T]) ConjureLogger[T]
}

type ConjureLogParam[T ConjureLogType] func(*T)

type defaultLogger[T ConjureLogType] struct {
	params  []ConjureLogParam[T]
	printer ConjureLogPrinter[T]
}

func NewDefaultLogger[T ConjureLogType](output io.Writer) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: jsonPrinter[T]{output: output},
	}
}

func NewDefaultLoggerWithPrinter[T ConjureLogType](printer ConjureLogPrinter[T]) ConjureLogger[T] {
	return &defaultLogger[T]{
		printer: printer,
	}
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
		return // TODO: something else?
	}
}

type ConjureLogPrinter[T ConjureLogType] interface {
	Print(log T) error
}

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
