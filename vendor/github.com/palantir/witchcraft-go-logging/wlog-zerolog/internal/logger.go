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
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog-zerolog/internal/marshalers"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/rs/zerolog"
)

type zeroLogEntry struct {
	evt *zerolog.Event
	wlog.MapValueEntries
}

func (e *zeroLogEntry) StringValue(key, value string) {
	e.evt = e.evt.Str(key, value)
}

func (e *zeroLogEntry) OptionalStringValue(key, value string) {
	if value != "" {
		e.StringValue(key, value)
	}
}

func (e *zeroLogEntry) StringListValue(k string, v []string) {
	if len(v) > 0 {
		e.evt.Strs(k, v)
	}
}

func (e *zeroLogEntry) SafeLongValue(key string, value int64) {
	e.evt = e.evt.Int64(key, value)
}

func (e *zeroLogEntry) IntValue(key string, value int32) {
	e.evt = e.evt.Int32(key, value)
}

func (e *zeroLogEntry) ObjectValue(k string, v interface{}, marshalerType reflect.Type) {
	ok := marshalers.EncodeType(e.evt, marshalerType, k, v)
	if !ok {
		e.evt.Interface(k, v)
	}
}

func (e *zeroLogEntry) Evt() *zerolog.Event {
	evt := e.evt
	for key, values := range e.StringMapValues() {
		key := key
		values := values
		dictEvt := zerolog.Dict()
		for k, v := range values {
			dictEvt = dictEvt.Str(k, v)
		}
		evt = evt.Dict(key, dictEvt)
	}
	for key, values := range e.AnyMapValues() {
		key := key
		values := values
		dictEvt := zerolog.Dict()
		for k, v := range values {
			dictEvt = dictEvt.Interface(k, v)
		}
		evt = evt.Dict(key, dictEvt)
	}
	return evt
}

type zeroLogger struct {
	logger zerolog.Logger
}

func (l *zeroLogger) Log(params ...wlog.Param) {
	logOutput(l.logger.Log, "", "", params)
}

func (l *zeroLogger) Debug(msg string, params ...wlog.Param) {
	logOutput(l.logger.Debug, msg, svc1log.LevelDebugValue, params)
}

func (l *zeroLogger) Info(msg string, params ...wlog.Param) {
	logOutput(l.logger.Info, msg, svc1log.LevelInfoValue, params)
}

func (l *zeroLogger) Warn(msg string, params ...wlog.Param) {
	logOutput(l.logger.Warn, msg, svc1log.LevelWarnValue, params)
}

func (l *zeroLogger) Error(msg string, params ...wlog.Param) {
	logOutput(l.logger.Error, msg, svc1log.LevelErrorValue, params)
}

func (l *zeroLogger) SetLevel(level wlog.LogLevel) {
	l.logger = l.logger.Level(toZeroLevel(level))
}

func logOutput(newEvt func() *zerolog.Event, msg, levelVal string, params []wlog.Param) {
	entry := &zeroLogEntry{
		evt: newEvt(),
	}
	if !entry.evt.Enabled() {
		return
	}
	entry.evt = entry.evt.Str(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	if levelVal != "" {
		entry.evt = entry.evt.Str(svc1log.LevelKey, levelVal)
	}
	wlog.ApplyParams(entry, params)
	entry.Evt().Msg(msg)
}
