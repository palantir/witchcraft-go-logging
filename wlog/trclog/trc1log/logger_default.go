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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

type defaultLogger struct {
	logger logentry.Logger
}

func (l *defaultLogger) Log(span wtracing.SpanModel, params ...Param) {
	log := poolTraceLogV1.Get()
	populateLogObject(log, span, params)
	l.logger.LogObject(log)
	poolTraceLogV1.Put(log)
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(span wtracing.SpanModel, inParams []Param) []wlog.Param {
	log := new(logtypes.TraceLogV1)
	populateLogObject(log, span, inParams)
	return []wlog.Param{logentry.Param(log)}
}

func (l *defaultLogger) Send(span wtracing.SpanModel) {
	l.Log(span)
}

// Close is a no-op that satisfies the wtracing.Reporter interface.
func (l *defaultLogger) Close() error {
	return nil
}
