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
)

const (
	TypeValue = "wrapped.1"

	WrappedEntityNameKey    = "entityName"
	WrappedEntityVersionKey = "entityVersion"

	PayloadKey             = "payload"
	PayloadTypeKey         = "type"
	PayloadServiceLogV1    = "serviceLogV1"
	PayloadRequestLogV2    = "requestLogV2"
	PayloadTraceLogV1      = "traceLogV1"
	PayloadEventLogV2      = "eventLogV2"
	PayloadMetricLogV1     = "metricLogV1"
	PayloadAuditLogV2      = "auditLogV2"
	PayloadAuditLogV3      = "auditLogV3"
	PayloadDiagnosticLogV1 = "diagnosticLogV1"
)

// Param has no implementations; it remains for API compatibility.
type Param interface {
	apply(entry wlog.LogEntry)
}

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	p.apply(entry)
}
