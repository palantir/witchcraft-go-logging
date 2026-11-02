// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package trc1log

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

// conjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type conjureLogger interface {
	LogTraceLogV1(*logging.TraceLogV1)
}

type defaultLogger struct {
	logger conjureLogger
	entry  wlog.Logger
}

func (l *defaultLogger) Log(span wtracing.SpanModel, params ...Param) {
	log := poolTraceLogV1.Get()
	populateLogObject(log, span, params)
	if l.logger != nil {
		l.logger.LogTraceLogV1(log)
	} else {
		l.entry.Log(objectToParam(log, span.LocalEndpoint != nil))
	}
	poolTraceLogV1.Put(log)
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(span wtracing.SpanModel, inParams []Param) []wlog.Param {
	log := new(logging.TraceLogV1)
	populateLogObject(log, span, inParams)
	return []wlog.Param{objectToParam(log, span.LocalEndpoint != nil)}
}

func ToObject(span wtracing.SpanModel, inParams []Param) logging.TraceLogV1 {
	var log logging.TraceLogV1
	populateLogObject(&log, span, inParams)
	return log
}

func (l *defaultLogger) Send(span wtracing.SpanModel) {
	l.Log(span)
}

// Close is a no-op that satisfies the wtracing.Reporter interface.
func (l *defaultLogger) Close() error {
	return nil
}
