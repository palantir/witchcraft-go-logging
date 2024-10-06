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
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

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
