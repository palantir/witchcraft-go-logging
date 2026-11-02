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
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
)

type wrappedEvt2Logger struct {
	name    string
	version string

	logger wlog.Logger
}

func (l *wrappedEvt2Logger) Event(name string, params ...evt2log.Param) {
	if direct, ok := l.logger.(conjureLogger); ok {
		emitObject(direct, l.name, l.version, logging.NewWrappedLogV1PayloadFromEventLogV2(evt2log.ToObject(name, params)))
		return
	}
	l.logger.Log(wrappedParams(l.name, l.version, evt2PayloadParams(name, params))...)
}
