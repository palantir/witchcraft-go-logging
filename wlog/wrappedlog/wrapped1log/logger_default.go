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
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
)

type DefaultLogger struct {
	name    string
	version string

	// creator and writer are used only by the request logger, to allow consumers to override the default creator with a Param
	creator wlog.LoggerCreator
	writer  io.Writer

	logger wlog.Logger
	// levellogger is only used by the service logger which supports logging at different log levels
	levellogger wlog.LeveledLogger
}

func (l *DefaultLogger) Audit() audit2log.Logger {
	return &wrappedAudit2Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

func (l *DefaultLogger) Diagnostic() diag1log.Logger {
	return &wrappedDiag1Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

func (l *DefaultLogger) Event() evt2log.Logger {
	return &wrappedEvt2Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

func (l *DefaultLogger) Metric() metric1log.Logger {
	return &wrappedMetric1Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

func (l *DefaultLogger) Request(params ...req2log.LoggerCreatorParam) req2log.Logger {
	loggerBuilder := &req2LoggerBuilder{
		name:          l.name,
		version:       l.version,
		loggerCreator: l.creator,
		idsExtractor:  extractor.NewDefaultIDsExtractor(),
	}
	for _, p := range params {
		p.Apply(loggerBuilder)
	}
	return loggerBuilder.build(l.writer)
}

func (l *DefaultLogger) Service(params ...svc1log.Param) svc1log.Logger {
	return &wrappedSvc1Logger{
		params:  params,
		name:    l.name,
		version: l.version,
		logger:  l.levellogger,
	}
}

func (l *DefaultLogger) Trace() trc1log.Logger {
	return &wrappedTrc1Logger{
		name:    l.name,
		version: l.version,
		logger:  l.logger,
	}
}

func (l *DefaultLogger) WithName(name string) {
	l.name = name
}

func (l *DefaultLogger) WithVersion(version string) {
	l.version = version
}

var defaultTypeParam = []wlog.Param{
	wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, TypeValue)
	}),
}
