package wlog

import (
	"io"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

var defaultPrinterCreator PrinterCreator = JSONPrinter

// GetDefaultPrinterCreator returns the default PrinterCreator.
func GetDefaultPrinterCreator() PrinterCreator {
	return defaultPrinterCreator
}

// SetDefaultPrinterCreator sets the default PrinterCreator.
func SetDefaultPrinterCreator(creator PrinterCreator) {
	defaultPrinterCreator = creator
}

// PrinterCreator is a function that creates a Printer. The default is JSONPrinter.
type PrinterCreator func(w io.Writer) Printer

// Printer is a generic interface for printing Conjure log objects.
type Printer interface {
	// Print writes the provided log object to an output.
	// Print should return errors sparingly: they are simply logged in plaintext at stderr.
	// Print should not retain the log object after the method returns.
	Print(log logging.LogType) error
}

// NoopPrinter returns a Printer that does nothing.
func NoopPrinter() Printer {
	return funcPrinter(func(logging.LogType) error { return nil })
}

type funcPrinter func(log logging.LogType) error

func (p funcPrinter) Print(log logging.LogType) error { return p(log) }
