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

package metric1log_test

import (
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/stretchr/testify/assert"
)

// paramsRecordingLogger records the org ID that the params of each call resolve to.
type paramsRecordingLogger struct {
	orgIDs []string
}

func (l *paramsRecordingLogger) Metric(name, typ string, params ...metric1log.Param) {
	l.record(params)
}

func (l *paramsRecordingLogger) record(params []metric1log.Param) {
	entry := wlog.NewMapLogEntry()
	for _, p := range params {
		metric1log.ApplyParam(p, entry)
	}
	l.orgIDs = append(l.orgIDs, entry.StringValues()[wlog.OrgIDKey])
}

// Tests that loggers derived from the same parent logger keep their own params: logging with one derived logger must
// not use the params that were added to a sibling.
func TestWithParamsDerivedLoggersDoNotShareParams(t *testing.T) {
	recorder := &paramsRecordingLogger{}

	// a logger that is wrapped repeatedly has more capacity than length in its params
	var parent metric1log.Logger = recorder
	for range 5 {
		parent = metric1log.WithParams(parent, metric1log.OrgID("parent"))
	}

	first := metric1log.WithParams(parent, metric1log.OrgID("first"))
	second := metric1log.WithParams(parent, metric1log.OrgID("second"))

	first.Metric("metric", "gauge")
	second.Metric("metric", "gauge")
	parent.Metric("metric", "gauge")

	assert.Equal(t, []string{"first", "second", "parent"}, recorder.orgIDs)
}
