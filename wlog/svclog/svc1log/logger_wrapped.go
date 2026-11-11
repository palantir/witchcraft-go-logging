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
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

type wrappedLogger struct {
	logger Logger
	params []Param
}

func (w *wrappedLogger) Debug(msg string, params ...Param) {
	if w.skip(wlog.DebugLevel, params) {
		return
	}
	w.logger.Debug(msg, wloginternal.AppendParams(w.params, params)...)
}

func (w *wrappedLogger) Info(msg string, params ...Param) {
	if w.skip(wlog.InfoLevel, params) {
		return
	}
	w.logger.Info(msg, wloginternal.AppendParams(w.params, params)...)
}

func (w *wrappedLogger) Warn(msg string, params ...Param) {
	if w.skip(wlog.WarnLevel, params) {
		return
	}
	w.logger.Warn(msg, wloginternal.AppendParams(w.params, params)...)
}

func (w *wrappedLogger) Error(msg string, params ...Param) {
	if w.skip(wlog.ErrorLevel, params) {
		return
	}
	w.logger.Error(msg, wloginternal.AppendParams(w.params, params)...)
}

// skip reports whether a disabled level makes appending the parameters wasted work.
func (w *wrappedLogger) skip(level wlog.LogLevel, params []Param) bool {
	return len(w.params) > 0 && len(params) > 0 && !w.Enabled(level)
}

func (w *wrappedLogger) SetLevel(level wlog.LogLevel) {
	w.logger.SetLevel(level)
}

func (w *wrappedLogger) Enabled(level wlog.LogLevel) bool {
	if l, ok := w.logger.(wlog.LevelChecker); ok {
		return l.Enabled(level)
	}
	return true
}
