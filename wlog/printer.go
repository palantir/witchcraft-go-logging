package wlog

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// LogPrinter is a generic interface for printing Conjure log objects.
type LogPrinter[T logging.LogTypes] interface {
	// Print writes the provided log object to an output.
	// Print should not retain the log object after the method returns.
	// Print should return errors sparingly: they are simply logged in plaintext at stderr.
	Print(log *T) error
}
