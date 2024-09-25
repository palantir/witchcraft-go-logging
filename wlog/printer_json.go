package wlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/palantir/pkg/bytesbuffers"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// The size of buffers allocated to consume log data.
// This is a tradeoff between memory usage and the number of allocations.
// Lines longer than this limit will not have their buffers reused.
const bytesBufferPoolAllocSize = 4096

var bytesBufferPool = bytesbuffers.NewSyncPool(bytesBufferPoolAllocSize)

// jsonPrinter is a LogPrinter that writes JSON-marshaled log objects to an output.
// It minimizes allocations by using a shared buffer pool.
type jsonPrinter[T logging.LogTypes] struct {
	output io.Writer
}

func (p jsonPrinter[T]) Print(log *T) error {
	buf := bytesBufferPool.Get()
	defer bytesBufferPool.Put(buf)

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(log); err != nil {
		return err
	}
	if buf.Len() == 0 {
		return nil
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
