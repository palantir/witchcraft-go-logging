package wlog

import (
	"encoding/json"
	"io"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// jsonPrinter is a Printer that writes JSON-marshaled log objects to an output.
// It minimizes allocations by using a shared buffer pool.
type jsonPrinter struct {
	out io.Writer
}

func JSONPrinter(out io.Writer) Printer {
	return jsonPrinter{out: out}
}

func (p jsonPrinter) Print(log logging.LogType) error {
	enc := json.NewEncoder(p.out)
	enc.SetEscapeHTML(false)
	// json.Encoder appends a newline after each JSON object.
	return enc.Encode(log)
}
