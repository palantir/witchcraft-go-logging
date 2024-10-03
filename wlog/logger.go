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

type LoggerCreator[T logging.LogTypes] func(w io.Writer) Logger[T]

// Logger is a generic logger that can log all Conjure log types.
type Logger[T logging.LogTypes] interface {
	Log(*T)
}

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
