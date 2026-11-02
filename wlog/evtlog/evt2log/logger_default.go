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

package evt2log

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

// conjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type conjureLogger interface {
	LogEventLogV2(*logging.EventLogV2)
}

type defaultLogger struct {
	logger conjureLogger
}

func (l *defaultLogger) Event(name string, params ...Param) {
	log := poolEventLogV2.Get()
	populateLogObject(log, name, params)
	l.logger.LogEventLogV2(log)
	poolEventLogV2.Put(log)
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(evtName string, inParams []Param) []wlog.Param {
	log := new(logging.EventLogV2)
	populateLogObject(log, evtName, inParams)
	return []wlog.Param{objectToParam(log)}
}

func ToObject(evtName string, inParams []Param) logging.EventLogV2 {
	var log logging.EventLogV2
	populateLogObject(&log, evtName, inParams)
	return log
}

type entryLogger struct {
	logger wlog.Logger
}

func (l entryLogger) LogEventLogV2(log *logging.EventLogV2) {
	l.logger.Log(objectToParam(log))
}
