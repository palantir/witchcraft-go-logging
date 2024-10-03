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

package wrapped1log

import (
	"io"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

// wrappedPrinter implements Printer for logs included in the wrapped.1 payload field.
// When an underlying log object is constructed and passed to the Print method,
// the delegate WrappedLogV1 logger is called with a new WrappedLogV1Payload object.
type wrappedPrinter[T logging.LogTypes] struct {
	delegate   wlog.Logger[logging.WrappedLogV1]
	newPayload func(payload T) logging.WrappedLogV1Payload
	params     []Param
}

func wrapPrinter[T logging.LogTypes](
	delegate wlog.Logger[logging.WrappedLogV1],
	newPayload func(payload T) logging.WrappedLogV1Payload,
	params []Param,
) wlog.LoggerCreator[T] {
	printer := wrappedPrinter[T]{delegate: delegate, newPayload: newPayload, params: params}
	return func(io.Writer) wlog.Logger[T] { return wlog.NewDefaultLoggerWithPrinter[T](printer) }
}

func (p wrappedPrinter[T]) Print(log logging.LogType) error {
	wloginternal.LogParams(p.delegate.Log, append(append([]Param{}, p.params...), Payload(p.newPayload(log.(T))))...)
	return nil
}
