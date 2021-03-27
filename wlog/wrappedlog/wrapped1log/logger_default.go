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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
)

type defaultLogger struct {
	name        string
	version     string
	logger      wlog.Logger
	levellogger wlog.LeveledLogger
}

func (l *defaultLogger) Request() req2log.Logger {
	return nil
}

func (l *defaultLogger) Service(params ...svc1log.Param) svc1log.Logger {
	return &wrappedSvc1Logger{
		params:  params,
		name:    l.name,
		version: l.version,
		logger:  l.levellogger,
	}
}

func (l *defaultLogger) Trace() trc1log.Logger {
	return &wrappedTrc1Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

var defaultTypeParam = []wlog.Param{
	wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, TypeValue)
	}),
}
