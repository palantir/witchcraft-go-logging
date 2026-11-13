// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

// Package origintest wraps service loggers from outside the wlog packages, like a library that instruments them.
package origintest

import (
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

const Package = "github.com/palantir/witchcraft-go-logging/internal/origintest"

type Logger struct {
	Delegate svc1log.Logger
}

func (l Logger) Debug(msg string, params ...svc1log.Param) { l.Delegate.Debug(msg, params...) }
func (l Logger) Info(msg string, params ...svc1log.Param)  { l.Delegate.Info(msg, params...) }
func (l Logger) Warn(msg string, params ...svc1log.Param)  { l.Delegate.Warn(msg, params...) }
func (l Logger) Error(msg string, params ...svc1log.Param) { l.Delegate.Error(msg, params...) }
func (l Logger) SetLevel(level wlog.LogLevel)              { l.Delegate.SetLevel(level) }
