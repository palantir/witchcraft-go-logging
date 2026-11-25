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

package audit3log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// paramsRecordingLogger records the last param of each call.
type paramsRecordingLogger struct {
	lastParams []Param
}

func (l *paramsRecordingLogger) Audit(name string, result AuditResultType, params ...Param) {
	l.lastParams = append(l.lastParams, params[len(params)-1])
}

// Tests that loggers derived from the same parent logger keep their own params: logging with one derived logger must
// not use the params that were added to a sibling.
func TestWithParamsDerivedLoggersDoNotShareParams(t *testing.T) {
	recorder := &paramsRecordingLogger{}

	// a logger that is wrapped repeatedly has more capacity than length in its params
	var parent Logger = recorder
	for range 5 {
		parent = WithParams(parent, OrgID("parent"))
	}

	firstParam := OrgID("first")
	secondParam := OrgID("second")
	first := WithParams(parent, firstParam)
	second := WithParams(parent, secondParam)

	first.Audit("audit", AuditResultSuccess)
	second.Audit("audit", AuditResultSuccess)

	if assert.Len(t, recorder.lastParams, 2) {
		assert.Same(t, firstParam, recorder.lastParams[0], "the first logger logged with a param of its sibling")
		assert.Same(t, secondParam, recorder.lastParams[1])
	}
}
