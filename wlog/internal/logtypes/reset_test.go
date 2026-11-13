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

package logtypes_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetKeepsCapacityAndClearsValues(t *testing.T) {
	audit := logtypes.AuditLogV3{
		Organizations: []logtypes.AuditLogV3Organization{{Id: "org"}},
		Categories:    []string{"category"},
		Entities:      []any{map[string]any{"k": "v"}},
		Users:         []logtypes.AuditLogV3ContextualizedUser{{Uid: "uid", Groups: []string{"group"}}},
		Origins:       []string{"origin"},
		RequestFields: logtypes.AnyMap{"k": "v"},
	}
	entities := audit.Entities
	audit.Reset()
	assert.Equal(t, 1, cap(audit.Entities))
	assert.Nil(t, entities[:1][0], "Reset must release referenced values")
	assert.Empty(t, audit.Users)
	assertEncodes(t, `{"type":"","logEntryId":"","time":"0001-01-01T00:00:00Z","name":"","result":""}`, &audit)

	trace := logtypes.TraceLogV1{Span: logtypes.TraceLogV1Span{Annotations: []logtypes.TraceLogV1SpanAnnotation{{Value: "sr"}}}}
	trace.Reset()
	assert.Equal(t, 1, cap(trace.Span.Annotations))
	assertEncodes(t, `{"type":"","time":"0001-01-01T00:00:00Z","span":{"traceId":"","id":"","name":"","timestamp":0,"duration":0}}`, &trace)
}

func assertEncodes(t *testing.T, want string, v any) {
	t.Helper()
	got, err := json.Marshal(v)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(got))
}
