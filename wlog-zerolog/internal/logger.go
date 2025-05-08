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

package zeroimpl

import (
	"reflect"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/rs/zerolog"
)

var _ wlog.LogEntry = (*zeroLogEntry)(nil)

// zeroLogEntryOp represents a single operation that operates on a *zeroLogEntry.
// The function does not take a parameter because it should be captured.
type zeroLogEntryOp func()

type zeroLogEntry struct {
	evt *zerolog.Event

	// a map from key to the operation that should be applied to the *zerolog.Event for that key.
	// This is needed because *zerolog.Event itself is append-only/set-once, but the wlog.LogEntry interface defines
	// overwrite behavior.
	entryOps map[string]zeroLogEntryOp

	// stores values for StringMapValue, AnyMapValue, and StringListValue
	mutableValueEntries wlog.MutableValueEntries
}

func zeroLogEntrySetValue[ValT any](entry *zeroLogEntry, fn func(k string, v ValT) *zerolog.Event, k string, v ValT) {
	entry.delete(k)
	entry.entryOps[k] = func() {
		entry.evt = fn(k, v)
	}
}

func (e *zeroLogEntry) delete(k string) {
	delete(e.entryOps, k)
	e.mutableValueEntries.Delete(k)
}

func (e *zeroLogEntry) StringValue(k, v string) {
	zeroLogEntrySetValue(e, e.evt.Str, k, v)
}

func (e *zeroLogEntry) OptionalStringValue(key, value string) {
	if value == "" {
		e.delete(key)
	} else {
		e.StringValue(key, value)
	}
}

func (e *zeroLogEntry) StringListValue(k string, v []string) {
	e.mutableValueEntries.StringListValue(k, v)
}

func (e *zeroLogEntry) StringListAppendValue(k string, v []string) {
	e.StringListValue(k, append(e.mutableValueEntries.StringListValues()[k], v...))
}

func (e *zeroLogEntry) SafeLongValue(k string, v int64) {
	zeroLogEntrySetValue(e, e.evt.Int64, k, v)
}

func (e *zeroLogEntry) IntValue(k string, v int32) {
	zeroLogEntrySetValue(e, e.evt.Int32, k, v)
}

func (e *zeroLogEntry) ObjectValue(k string, v interface{}, marshalerType reflect.Type) {
	zeroLogEntrySetValue(e, e.evt.Interface, k, v)
}

func (e *zeroLogEntry) StringMapValue(key string, values map[string]string) {
	delete(e.entryOps, key)
	e.mutableValueEntries.StringMapValue(key, values)
}

func (e *zeroLogEntry) AnyMapValue(key string, values map[string]interface{}) {
	delete(e.entryOps, key)
	e.mutableValueEntries.AnyMapValue(key, values)
}

func (e *zeroLogEntry) ObjectListValue(k string, v []any) {
	e.mutableValueEntries.ObjectListValue(k, v)
}

func (e *zeroLogEntry) ObjectListAppendValue(k string, v []any) {
	e.ObjectListValue(k, append(e.mutableValueEntries.ObjectListValues()[k], v...))
}

func (e *zeroLogEntry) Evt() *zerolog.Event {
	evt := e.evt

	for _, opFn := range e.entryOps {
		opFn()
	}

	for k, v := range e.mutableValueEntries.StringListValues() {
		evt = e.evt.Strs(k, v)
	}

	for k, v := range e.mutableValueEntries.ObjectListValues() {
		evt = e.evt.Any(k, v)
	}

	evt = addMapToEvt(evt, e.mutableValueEntries.StringMapValues(), func(dictEvt *zerolog.Event, k string, v string) *zerolog.Event {
		return dictEvt.Str(k, v)
	})

	evt = addMapToEvt(evt, e.mutableValueEntries.AnyMapValues(), func(dictEvt *zerolog.Event, k string, v any) *zerolog.Event {
		return dictEvt.Interface(k, v)
	})

	return evt
}

func addMapToEvt[ValT any](evt *zerolog.Event, inputMap map[string]map[string]ValT, evtFn func(dictEvt *zerolog.Event, k string, v ValT) *zerolog.Event) *zerolog.Event {
	for key, values := range inputMap {
		dictEvt := zerolog.Dict()
		for k, v := range values {
			dictEvt = evtFn(dictEvt, k, v)
		}
		evt = evt.Dict(key, dictEvt)
	}
	return evt
}

type zeroLogger struct {
	logger zerolog.Logger
	*wlog.AtomicLogLevel
}

func (l *zeroLogger) Log(params ...wlog.Param) {
	logOutput(l.logger.Log, "", params)
}

func (l *zeroLogger) Debug(msg string, params ...wlog.Param) {
	if l.Enabled(wlog.DebugLevel) {
		logOutput(l.logger.Log, msg, params)
	}
}

func (l *zeroLogger) Info(msg string, params ...wlog.Param) {
	if l.Enabled(wlog.InfoLevel) {
		logOutput(l.logger.Log, msg, params)
	}
}

func (l *zeroLogger) Warn(msg string, params ...wlog.Param) {
	if l.Enabled(wlog.WarnLevel) {
		logOutput(l.logger.Log, msg, params)
	}
}

func (l *zeroLogger) Error(msg string, params ...wlog.Param) {
	if l.Enabled(wlog.ErrorLevel) {
		logOutput(l.logger.Log, msg, params)
	}
}

func logOutput(newEvt func() *zerolog.Event, msg string, params []wlog.Param) {
	entry := &zeroLogEntry{
		evt:      newEvt(),
		entryOps: make(map[string]zeroLogEntryOp),
	}
	if !entry.evt.Enabled() {
		return
	}
	wlog.ApplyParams(entry, params)
	entry.Evt().Msg(msg)
}
