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
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type wrappedSvc1Logger struct {
	name    string
	version string
	params  []svc1log.Param

	logger wlog.LeveledLogger
	level  wlog.LevelChecker
}

func (l *wrappedSvc1Logger) Debug(msg string, params ...svc1log.Param) {
	if l.Enabled(wlog.DebugLevel) {
		l.log(svc1log.LevelDebug, msg, params)
	}
}

func (l *wrappedSvc1Logger) Info(msg string, params ...svc1log.Param) {
	if l.Enabled(wlog.InfoLevel) {
		l.log(svc1log.LevelInfo, msg, params)
	}
}

func (l *wrappedSvc1Logger) Warn(msg string, params ...svc1log.Param) {
	if l.Enabled(wlog.WarnLevel) {
		l.log(svc1log.LevelWarn, msg, params)
	}
}

func (l *wrappedSvc1Logger) Error(msg string, params ...svc1log.Param) {
	if l.Enabled(wlog.ErrorLevel) {
		l.log(svc1log.LevelError, msg, params)
	}
}

func (l *wrappedSvc1Logger) SetLevel(level wlog.LogLevel) {
	l.logger.SetLevel(level)
}

func (l *wrappedSvc1Logger) Enabled(level wlog.LogLevel) bool {
	return l.level == nil || l.level.Enabled(level)
}

func (l *wrappedSvc1Logger) log(level logging.LogLevel, message string, params []svc1log.Param) {
	// Evaluate params before provider dispatch so caller attribution is independent of the backend.
	obj := svc1log.ToObject(level, message, wloginternal.AppendParams(l.params, params))
	if direct, ok := l.logger.(conjureLogger); ok {
		emitObject(direct, l.name, l.version, logging.NewWrappedLogV1PayloadFromServiceLogV1(obj))
		return
	}

	entryParams := wrappedParams(l.name, l.version, svc1PayloadParams(obj))
	switch level.Value() {
	case logging.LogLevel_DEBUG:
		l.logger.Debug("", entryParams...)
	case logging.LogLevel_INFO:
		l.logger.Info("", entryParams...)
	case logging.LogLevel_WARN:
		l.logger.Warn("", entryParams...)
	case logging.LogLevel_ERROR:
		l.logger.Error("", entryParams...)
	}
}
