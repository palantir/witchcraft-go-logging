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

package metric1log

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

// conjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type conjureLogger interface {
	LogMetricLogV1(*logging.MetricLogV1)
}

type defaultLogger struct {
	logger conjureLogger
}

func (l *defaultLogger) Metric(name, typ string, params ...Param) {
	log := poolMetricLogV1.Get()
	populateLogObject(log, name, typ, params)
	l.logger.LogMetricLogV1(log)
	poolMetricLogV1.Put(log)
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(metricName, metricType string, inParams []Param) []wlog.Param {
	log := new(logging.MetricLogV1)
	populateLogObject(log, metricName, metricType, inParams)
	return []wlog.Param{objectToParam(log)}
}

func ToObject(metricName, metricType string, inParams []Param) logging.MetricLogV1 {
	var log logging.MetricLogV1
	populateLogObject(&log, metricName, metricType, inParams)
	return log
}

type entryLogger struct {
	logger wlog.Logger
}

func (l entryLogger) LogMetricLogV1(log *logging.MetricLogV1) {
	l.logger.Log(objectToParam(log))
}
