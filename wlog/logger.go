// Copyright (c) 2024 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wlog

import (
	"fmt"
	"io"
	"os"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// LoggerProvider is a function that creates a Printer. The default is JSONPrinter.
type LoggerProvider func(w io.Writer) Printer

// LoggerCreator is a generic constructor for a Logger. The default is NewDefaultLogger.
type LoggerCreator[T logging.LogTypes] func(w io.Writer) Logger[T]

// Logger is a generic logger that can log all Conjure log types.
// It should encapsulate its output write (if applicable) and handle any errors gracefully.
type Logger[T logging.LogTypes] interface {
	Log(*T)
}

// The global variable holding our default LoggerProvider. Settable with SetDefaultLoggerProvider.
var defaultLoggerProvider LoggerProvider = JSONPrinter

// DefaultLoggerProvider returns the default LoggerProvider.
func DefaultLoggerProvider() LoggerProvider {
	return defaultLoggerProvider
}

// SetDefaultLoggerProvider sets the default LoggerProvider.
func SetDefaultLoggerProvider(provider LoggerProvider) {
	defaultLoggerProvider = provider
}

// NewDefaultLogger creates a new logger that writes JSON-marshaled lines to the provided output.
func NewDefaultLogger[T logging.LogTypes](output io.Writer) Logger[T] {
	return NewDefaultLoggerWithLoggerProvider[T](DefaultLoggerProvider())(output)
}

// NewDefaultLoggerWithLoggerProvider creates a new logger that writes log objects using the provided logger provider.
func NewDefaultLoggerWithLoggerProvider[T logging.LogTypes](provider LoggerProvider) LoggerCreator[T] {
	return func(w io.Writer) Logger[T] {
		return NewDefaultLoggerWithPrinter[T](provider(w))
	}
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
	err := l.printer.Print(any(*log).(logging.LogType))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		// TODO: something else?
	}
}
