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
	//keys map[string]struct{}

	// a map from key to the operation that should be applied to the *zerolog.Event for that key.
	// This is needed because *zerolog.Event itself is append-only/set-once, but the wlog.LogEntry interface defines
	// overwrite behavior.
	entryOps map[string]zeroLogEntryOp

	// stores values for StringMapValue, AnyMapValue, and
	wlog.MutableValueEntries
}

func (e *zeroLogEntry) StringValue(key, value string) {
	e.DeleteKey(key)
	e.entryOps[key] = func() {
		e.evt = e.evt.Str(key, value)
	}
}

func (e *zeroLogEntry) OptionalStringValue(key, value string) {
	e.DeleteKey(key)
	if value == "" {
		delete(e.entryOps, key)
	} else {
		e.StringValue(key, value)
	}
}

func (e *zeroLogEntry) StringListValue(key string, value []string) {
	e.DeleteKey(key)
	e.MutableValueEntries.StringListValue(key, value)
}

func (e *zeroLogEntry) StringListValueAppend(k string, v []string) {
	e.StringListValue(k, append(e.MutableValueEntries.StringListValues()[k], v...))
}

func (e *zeroLogEntry) SafeLongValue(key string, value int64) {
	e.DeleteKey(key)
	e.entryOps[key] = func() {
		e.evt = e.evt.Int64(key, value)
	}
}

func (e *zeroLogEntry) IntValue(key string, value int32) {
	e.DeleteKey(key)
	e.entryOps[key] = func() {
		e.evt = e.evt.Int32(key, value)
	}
}

func (e *zeroLogEntry) ObjectValue(key string, value interface{}, marshalerType reflect.Type) {
	e.DeleteKey(key)
	e.entryOps[key] = func() {
		e.evt.Interface(key, value)
	}
}

// StringMapValue adds or merges the strings in values
// Since wlog overrides duplicates with a preference for the last parameter
// The parameters should not replace an existing key because parameters are passed to zerolog in reverse
// This differs from the default wlog StringMapValue since parameters are not reversed
func (e *zeroLogEntry) StringMapValue(key string, values map[string]string) {
	//mapValueHelper(&e.stringMapValues, key, values)
	delete(e.entryOps, key)
	e.MutableValueEntries.StringMapValue(key, values)
}

// AnyMapValue adds or merges the values in values
// Since wlog overrides duplicates with a preference for the last parameter
// The parameters should not replace an existing key because parameters are passed to zerolog in reverse
// This differs from the default wlog AnyMapValue since parameters are not reversed
func (e *zeroLogEntry) AnyMapValue(key string, values map[string]interface{}) {
	//mapValueHelper(&e.anyMapValues, key, values)
	delete(e.entryOps, key)
	e.MutableValueEntries.AnyMapValue(key, values)
}

func (e *zeroLogEntry) DeleteKey(key string) {
	e.MutableValueEntries.DeleteKey(key)
	delete(e.entryOps, key)
}

func (e *zeroLogEntry) Evt() *zerolog.Event {
	evt := e.evt

	// apply all operations
	for _, opFn := range e.entryOps {
		opFn()
	}

	for k, v := range e.MutableValueEntries.StringListValues() {
		evt = e.evt.Strs(k, v)
	}

	evt = addMapToEvt(evt, e.StringMapValues(), func(dictEvt *zerolog.Event, k string, v string) *zerolog.Event {
		return dictEvt.Str(k, v)
	})

	evt = addMapToEvt(evt, e.AnyMapValues(), func(dictEvt *zerolog.Event, k string, v any) *zerolog.Event {
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
