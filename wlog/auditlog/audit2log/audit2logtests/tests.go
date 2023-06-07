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

package audit2logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name          string
	UID           string
	SID           string
	TokenID       string
	OrgID         string
	TraceID       string
	OtherUIDs     []string
	Origin        string
	AuditName     string
	AuditResult   audit2log.AuditResultType
	RequestParams map[string]interface{}
	ResultParams  map[string]interface{}
	JSONMatcher   objmatcher.MapMatcher
}

func (tc TestCase) Params() []audit2log.Param {
	return []audit2log.Param{
		audit2log.UID(tc.UID),
		audit2log.SID(tc.SID),
		audit2log.TokenID(tc.TokenID),
		audit2log.OrgID(tc.OrgID),
		audit2log.TraceID(tc.TraceID),
		audit2log.OtherUIDs(tc.OtherUIDs...),
		audit2log.Origin(tc.Origin),
		audit2log.RequestParams(tc.RequestParams),
		audit2log.ResultParams(tc.ResultParams),
	}
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name:          "basic audit log entry",
			UID:           "user-1",
			SID:           "session-1",
			TokenID:       "X-Y-Z",
			OrgID:         "org-1",
			TraceID:       "trace-id-1",
			OtherUIDs:     []string{"user-2", "user-3"},
			Origin:        "0.0.0.0",
			AuditName:     "AUDITED_ACTION_NAME",
			AuditResult:   audit2log.AuditResultSuccess,
			RequestParams: map[string]interface{}{"requestKey": "requestValue"},
			ResultParams:  map[string]interface{}{"resultKey": "resultValue"},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":      objmatcher.NewRegExpMatcher(".+"),
				"uid":       objmatcher.NewEqualsMatcher("user-1"),
				"sid":       objmatcher.NewEqualsMatcher("session-1"),
				"tokenId":   objmatcher.NewEqualsMatcher("X-Y-Z"),
				"orgId":     objmatcher.NewEqualsMatcher("org-1"),
				"traceId":   objmatcher.NewEqualsMatcher("trace-id-1"),
				"otherUids": objmatcher.NewEqualsMatcher([]interface{}{"user-2", "user-3"}),
				"origin":    objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":      objmatcher.NewEqualsMatcher("AUDITED_ACTION_NAME"),
				"result":    objmatcher.NewEqualsMatcher("SUCCESS"),
				"requestParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"requestKey": objmatcher.NewEqualsMatcher("requestValue"),
				}),
				"resultParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"resultKey": objmatcher.NewEqualsMatcher("resultValue"),
				}),
				"type": objmatcher.NewEqualsMatcher("audit.2"),
			}),
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) audit2log.Logger) {
	jsonOutputTests(t, loggerProvider)
	rParamIsntOverwrittenByRParamsTest(t, loggerProvider)
	extraRParamsDoNotAppear(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer) audit2log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Audit(
				tc.AuditName,
				tc.AuditResult,
				audit2log.UID(tc.UID),
				audit2log.SID(tc.SID),
				audit2log.TokenID(tc.TokenID),
				audit2log.OrgID(tc.OrgID),
				audit2log.TraceID(tc.TraceID),
				audit2log.OtherUIDs(tc.OtherUIDs...),
				audit2log.Origin(tc.Origin),
				audit2log.RequestParams(tc.RequestParams),
				audit2log.ResultParams(tc.ResultParams))

			gotAuditLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotAuditLog)
			require.NoError(t, err, "Case %d: %s\nAudit log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotAuditLog), "Case %d: %s", i, tc.Name)
		})
	}
}

// Verifies that if different parameters are specified using ResultParam/RequestParam and ResultParams/RequestParams,
// all of the values are present in the final output (that is, these parameters should be additive).
func rParamIsntOverwrittenByRParamsTest(t *testing.T, loggerProvider func(w io.Writer) audit2log.Logger) {
	mapFieldMatcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"key1": objmatcher.NewEqualsMatcher("val1"),
		"key2": objmatcher.NewEqualsMatcher("val2"),
	})
	for i, tc := range []struct {
		name   string
		params []audit2log.Param
		want   objmatcher.MapMatcher
	}{
		{
			name: "ResultParam params are additive",
			params: []audit2log.Param{
				audit2log.ResultParam("key1", "val1"),
				audit2log.ResultParams(map[string]interface{}{"key2": "val2"}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":         objmatcher.NewRegExpMatcher(".+"),
				"name":         objmatcher.NewEqualsMatcher("audited action name"),
				"type":         objmatcher.NewEqualsMatcher("audit.2"),
				"result":       objmatcher.NewEqualsMatcher("SUCCESS"),
				"resultParams": mapFieldMatcher,
			}),
		},
		{
			name: "RequestParam params are additive",
			params: []audit2log.Param{
				audit2log.RequestParam("key1", "val1"),
				audit2log.RequestParams(map[string]interface{}{
					"key2": "val2",
				}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":          objmatcher.NewRegExpMatcher(".+"),
				"name":          objmatcher.NewEqualsMatcher("audited action name"),
				"type":          objmatcher.NewEqualsMatcher("audit.2"),
				"result":        objmatcher.NewEqualsMatcher("SUCCESS"),
				"requestParams": mapFieldMatcher,
			}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			logger.Audit("audited action name", audit2log.AuditResultSuccess, tc.params...)

			auditLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &auditLog)
			require.NoError(
				t,
				err,
				"Case %d: %s\nAudit log line is not a valid map: %v",
				i,
				tc.name,
				string(logEntry))
			assert.NoError(t, tc.want.Matches(auditLog), "Case %d: %s", i, tc.name)
		})
	}
}

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraRParamsDoNotAppear(t *testing.T, loggerProvider func(w io.Writer) audit2log.Logger) {
	const (
		resultParamsKey  = "resultParams"
		requestParamsKey = "requestParams"
	)

	for i, tc := range []struct {
		name       string
		paramKey   string
		paramFunc  func(key string, val interface{}) audit2log.Param
		paramsFunc func(map[string]interface{}) audit2log.Param
	}{
		{
			name:       "Params stay separate across calls for ResultParam",
			paramKey:   resultParamsKey,
			paramFunc:  audit2log.ResultParam,
			paramsFunc: audit2log.ResultParams,
		},
		{
			name:       "Params stay separate across calls for RequestParam",
			paramKey:   requestParamsKey,
			paramFunc:  audit2log.RequestParam,
			paramsFunc: audit2log.RequestParams,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			reusedParams := tc.paramsFunc(map[string]interface{}{"key1": "val1"})

			logger.Audit(
				"audited action name",
				audit2log.AuditResultSuccess,
				reusedParams,
				tc.paramFunc("key2", "val2"))
			want := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":   objmatcher.NewRegExpMatcher(".+"),
				"name":   objmatcher.NewEqualsMatcher("audited action name"),
				"type":   objmatcher.NewEqualsMatcher("audit.2"),
				"result": objmatcher.NewEqualsMatcher("SUCCESS"),
				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"key1": objmatcher.NewEqualsMatcher("val1"),
					"key2": objmatcher.NewEqualsMatcher("val2"),
				}),
			})
			auditLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := json.Unmarshal(logEntry, &auditLog)
			require.NoError(
				t,
				err,
				"Case %d: %s\nAudit log is not a valid map: %v",
				i,
				tc.name,
				string(logEntry))
			assert.NoError(t, want.Matches(auditLog), "Case %d: %s", i, tc.name)

			buf.Reset()
			logger.Audit("audited action name", audit2log.AuditResultSuccess, reusedParams)

			want = objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":   objmatcher.NewRegExpMatcher(".+"),
				"name":   objmatcher.NewEqualsMatcher("audited action name"),
				"type":   objmatcher.NewEqualsMatcher("audit.2"),
				"result": objmatcher.NewEqualsMatcher("SUCCESS"),
				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"key1": objmatcher.NewEqualsMatcher("val1"),
				}),
			})

			auditLog = map[string]interface{}{}
			logEntry = buf.Bytes()
			err = json.Unmarshal(logEntry, &auditLog)
			require.NoError(
				t,
				err,
				"Case %d: %s\nAudit log is not a valid map: %v",
				i,
				tc.name,
				string(logEntry))
			assert.NoError(t, want.Matches(auditLog), "Case %d: %s", i, tc.name)
		})
	}
}
