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
	"slices"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
)

var _ Logger = (*defaultLogger)(nil)

// defaultLogger builds each family's own logger on a sink that wraps its log objects.
type defaultLogger struct {
	name    string
	version string
	writer  io.Writer
	creator wlog.LoggerCreator

	sink        *sink
	leveledSink *sink
}

func (l *defaultLogger) logger(io.Writer) wlog.Logger {
	return l.sink
}

// leveledLogger ignores the level so that all service loggers share the sink's level.
func (l *defaultLogger) leveledLogger(io.Writer, wlog.LogLevel) wlog.LeveledLogger {
	return l.leveledSink
}

func (l *defaultLogger) Audit() audit2log.Logger {
	return audit2log.NewFromCreator(l.writer, l.logger)
}

func (l *defaultLogger) AuditV3() audit3log.Logger {
	return audit3log.NewFromCreator(l.writer, l.logger)
}

func (l *defaultLogger) Diagnostic() diag1log.Logger {
	return diag1log.NewFromCreator(l.writer, l.logger)
}

func (l *defaultLogger) Event() evt2log.Logger {
	return evt2log.NewFromCreator(l.writer, l.logger)
}

func (l *defaultLogger) Metric() metric1log.Logger {
	return metric1log.NewFromCreator(l.writer, l.logger)
}

// Request writes wrapped logs to the logger of a Creator param, if one is given.
func (l *defaultLogger) Request(params ...req2log.LoggerCreatorParam) req2log.Logger {
	var b creatorBuilder
	for _, p := range params {
		p.Apply(&b)
	}
	creator := l.logger
	if b.creator != nil {
		creator = func(w io.Writer) wlog.Logger { return newSink(l.name, l.version, b.creator(w)) }
	}
	return req2log.NewFromCreator(l.writer, creator, append(slices.Clip(params), req2log.Creator(creator))...)
}

func (l *defaultLogger) Service(params ...svc1log.Param) svc1log.Logger {
	return svc1log.NewFromCreator(l.writer, wlog.InfoLevel, l.leveledLogger, params...)
}

func (l *defaultLogger) Trace() trc1log.Logger {
	return trc1log.NewFromCreator(l.writer, l.logger)
}

// creatorBuilder records the last Creator param and ignores all other params.
type creatorBuilder struct {
	creator wlog.LoggerCreator
}

func (b *creatorBuilder) LoggerCreator(creator wlog.LoggerCreator) { b.creator = creator }
func (*creatorBuilder) IdsExtractor(extractor.IDsFromRequest)      {}
func (*creatorBuilder) SafePathParams([]string)                    {}
func (*creatorBuilder) ForbiddenPathParams([]string)               {}
func (*creatorBuilder) SafeQueryParams([]string)                   {}
func (*creatorBuilder) ForbiddenQueryParams([]string)              {}
func (*creatorBuilder) SafeHeaderParams([]string)                  {}
func (*creatorBuilder) ForbiddenHeaderParams([]string)             {}
