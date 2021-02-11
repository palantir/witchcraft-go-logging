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

package svc1log

import (
	"github.com/palantir/pkg/refreshable"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

type OriginLevelLoggerConfig struct {
	Level           wlog.Level
	PerOriginLevels map[string]wlog.Level
}

type originLevelLogger struct {
	loggerCreator func(level wlog.LogLevel) wlog.LeveledLogger
	logger        wlog.LeveledLogger
	config        refreshable.Refreshable
	params        []Param
}

func (l *originLevelLogger) Debug(msg string, params ...Param) {
	origin := l.getOriginFromParams(params)
	if !l.shouldLog(origin, wlog.Debug) {
		return
	}
	l.logger.Debug(msg, l.toParams(params)...)
}

func (l *originLevelLogger) Info(msg string, params ...Param) {
	origin := l.getOriginFromParams(params)
	if !l.shouldLog(origin, wlog.Info) {
		return
	}
	l.logger.Info(msg, l.toParams(params)...)
}

func (l *originLevelLogger) Warn(msg string, params ...Param) {
	origin := l.getOriginFromParams(params)
	if !l.shouldLog(origin, wlog.Warn) {
		return
	}
	l.logger.Warn(msg, l.toParams(params)...)
}

func (l *originLevelLogger) Error(msg string, params ...Param) {
	origin := l.getOriginFromParams(params)
	if !l.shouldLog(origin, wlog.Error) {
		return
	}
	l.logger.Error(msg, l.toParams(params)...)
}

func (l *originLevelLogger) SetLevel(level wlog.LogLevel) {
	l.logger.SetLevel(level)
}

func (l *originLevelLogger) toParams(inParams []Param) []wlog.Param {
	if len(inParams) == 0 {
		return defaultTypeParam
	}
	outParams := make([]wlog.Param, len(defaultTypeParam)+len(inParams))
	copy(outParams, defaultTypeParam)
	for idx := range inParams {
		outParams[len(defaultTypeParam)+idx] = wlog.NewParam(inParams[idx].apply)
	}
	return outParams
}

func (l *originLevelLogger) shouldLog(origin string, level wlog.Level) bool {
	config := l.config.Current().(OriginLevelLoggerConfig)

	if originLevel, ok := config.PerOriginLevels[origin]; ok {
		return originLevel <= level
	}
	return config.Level <= level
}

func (l *originLevelLogger) getOriginFromParams(in []Param) string {
	// iterate backwards through the params (most specific to call site first) and grab the first originFunc type
	for i := len(in) - 1; i >= 0; i-- {
		if param, ok := in[i].(originParamFunc); ok {
			e := wlog.NewMapLogEntry()
			param.apply(e)
			return e.StringValues()[OriginKey]
		}
	}
	for i := len(l.params); i >= 0; i-- {
		if param, ok := in[i].(originParamFunc); ok {
			e := wlog.NewMapLogEntry()
			param.apply(e)
			return e.StringValues()[OriginKey]
		}
	}
	return ""
}
