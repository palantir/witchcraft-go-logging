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

package svc1log

import (
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
)

var DebugLevelParam = wlog.NewParam(func(entry wlog.LogEntry) {
	entry.StringValue(LevelKey, LevelDebugValue)
})
var InfoLevelParam = wlog.NewParam(func(entry wlog.LogEntry) {
	entry.StringValue(LevelKey, LevelInfoValue)
})
var WarnLevelParam = wlog.NewParam(func(entry wlog.LogEntry) {
	entry.StringValue(LevelKey, LevelWarnValue)
})
var ErrorLevelParam = wlog.NewParam(func(entry wlog.LogEntry) {
	entry.StringValue(LevelKey, LevelErrorValue)
})

type DefaultLogger struct {
	loggerCreator func(level wlog.LogLevel) wlog.LeveledLogger

	logger wlog.LeveledLogger
	params []Param
}

func (l *DefaultLogger) Debug(msg string, params ...Param) {
	l.logger.Debug("", l.ToParams(msg, DebugLevelParam, params)...)
}

func (l *DefaultLogger) Info(msg string, params ...Param) {
	l.logger.Info("", l.ToParams(msg, InfoLevelParam, params)...)
}

func (l *DefaultLogger) Warn(msg string, params ...Param) {
	l.logger.Warn("", l.ToParams(msg, WarnLevelParam, params)...)
}

func (l *DefaultLogger) Error(msg string, params ...Param) {
	l.logger.Error("", l.ToParams(msg, ErrorLevelParam, params)...)
}

func (l *DefaultLogger) SetLevel(level wlog.LogLevel) {
	l.logger.SetLevel(level)
}

func (l *DefaultLogger) ToParams(msg string, level wlog.Param, inParams []Param) []wlog.Param {
	outParams := make([]wlog.Param, len(defaultParams)+2+len(inParams))
	copy(outParams, defaultParams)
	outParams[len(defaultParams)] = level
	outParams[len(defaultParams)+1] = wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(MessageKey, msg)
	})
	for idx := range inParams {
		outParams[len(defaultParams)+2+idx] = wlog.NewParam(inParams[idx].apply)
	}
	return outParams
}

var defaultParams = []wlog.Param{
	wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, TypeValue)
		entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	}),
}
