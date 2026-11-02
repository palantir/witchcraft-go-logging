// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

type wrappedTrc1Logger struct {
	name    string
	version string

	logger wlog.Logger
}

func (l *wrappedTrc1Logger) Log(span wtracing.SpanModel, params ...trc1log.Param) {
	if direct, ok := l.logger.(conjureLogger); ok {
		emitObject(direct, l.name, l.version, logging.NewWrappedLogV1PayloadFromTraceLogV1(trc1log.ToObject(span, params)))
		return
	}
	l.logger.Log(wrappedParams(l.name, l.version, trc1PayloadParams(span, params))...)
}

func (l *wrappedTrc1Logger) Send(span wtracing.SpanModel) {
	l.Log(span)
}

// Close is a no-op that satisfies the wtracing.Reporter interface.
func (l *wrappedTrc1Logger) Close() error {
	return nil
}
