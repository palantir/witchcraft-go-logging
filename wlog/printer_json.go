package wlog

import (
	"encoding/json"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"io"
)

// The size of buffers allocated to consume log data.
// This is a tradeoff between memory usage and the number of allocations.
// Lines longer than this limit will not have their buffers reused.
//const bytesBufferPoolAllocSize = 4096
//
//var bytesBufferPool = bytesbuffers.NewSyncPool(bytesBufferPoolAllocSize)

// jsonPrinter is a LogPrinter that writes JSON-marshaled log objects to an output.
// It minimizes allocations by using a shared buffer pool.
type jsonPrinter[T logging.LogTypes] struct {
	enc *json.Encoder
	out io.Writer
}

func JSONPrinter[T logging.LogTypes](out io.Writer) LogPrinter[T] {
	return &jsonPrinter[T]{enc: newJSONEncoder(out), out: out}
}

func (p *jsonPrinter[T]) Print(log *T) error {
	// write to output
	if err := p.enc.Encode(log); err != nil {
		// In case of error, the json encoder will cache and return the same error for every subsequent call.
		// Reset the encoder to clear the error and continue.
		p.enc = newJSONEncoder(p.out)
		return err
	}
	return nil
}

func newJSONEncoder(out io.Writer) *json.Encoder {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	return enc
}
