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

package audit2log_test

import (
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/stretchr/testify/assert"
)

// paramsRecordingLogger records the org ID that the params of each call resolve to.
type paramsRecordingLogger struct {
	orgIDs []string
}

func (l *paramsRecordingLogger) Audit(name string, result audit2log.AuditResultType, params ...audit2log.Param) {
	l.record(params)
}

func (l *paramsRecordingLogger) record(params []audit2log.Param) {
	entry := wlog.NewMapLogEntry()
	for _, p := range params {
		audit2log.ApplyParam(p, entry)
	}
	l.orgIDs = append(l.orgIDs, entry.StringValues()[wlog.OrgIDKey])
}

// Tests that loggers derived from the same parent logger keep their own params: logging with one derived logger must
// not use the params that were added to a sibling.
func TestWithParamsDerivedLoggersDoNotShareParams(t *testing.T) {
	recorder := &paramsRecordingLogger{}

	// a logger that is wrapped repeatedly has more capacity than length in its params
	var parent audit2log.Logger = recorder
	for range 5 {
		parent = audit2log.WithParams(parent, audit2log.OrgID("parent"))
	}

	first := audit2log.WithParams(parent, audit2log.OrgID("first"))
	second := audit2log.WithParams(parent, audit2log.OrgID("second"))

	first.Audit("audit", audit2log.AuditResultSuccess)
	second.Audit("audit", audit2log.AuditResultSuccess)
	parent.Audit("audit", audit2log.AuditResultSuccess)

	assert.Equal(t, []string{"first", "second", "parent"}, recorder.orgIDs)
}
