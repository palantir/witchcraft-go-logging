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

package trc1log

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

var objectPool = wloginternal.NewPool((*logging.TraceLogV1).Reset)

type defaultLogger struct {
	logger wlog.Logger[logging.TraceLogV1]
}

func (l *defaultLogger) Log(span wtracing.SpanModel, params ...Param) {
	log := objectPool.Get()
	wloginternal.ApplyParams(log, Type(), TimeNow(), Span(span))
	wloginternal.ApplyParams(log, params...)
	l.logger.Log(log)
	objectPool.Put(log)
}

func (l *defaultLogger) Send(span wtracing.SpanModel) {
	l.Log(span)
}

func (l *defaultLogger) Close() error {
	return nil
}
