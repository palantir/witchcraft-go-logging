// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package wrapped1logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	v2 "github.com/palantir/witchcraft-go-logging/conjure/foundry/audit/api/category/v2"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	typeValue       = "audit.3"
	resultValue     = "SUCCESS"
	uuidRegexpValue = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$"
)

type Audit3TestCase struct {
	Name          string
	UID           string
	SID           string
	TokenID       string
	OrgID         string
	TraceID       string
	Origin        string
	Category      v2.AuditCategoryV2
	AuditName     string
	AuditResult   audit3log.AuditResultType
	RequestFields map[string]any
	ResultFields  map[string]any
	JSONMatcher   objmatcher.MapMatcher
}

func (tc Audit3TestCase) Params() []audit3log.Param {
	return []audit3log.Param{
		audit3log.UID(tc.UID),
		audit3log.SID(tc.SID),
		audit3log.TokenID(tc.TokenID),
		audit3log.OrgID(tc.OrgID),
		audit3log.TraceID(tc.TraceID),
		audit3log.Origin(tc.Origin),
		audit3log.Category(tc.Category),
		audit3log.RequestFields(tc.RequestFields),
		audit3log.ResultFields(tc.ResultFields),
	}
}

func Audit3TestCases(entityName, entityVersion string) []Audit3TestCase {
	return []Audit3TestCase{
		{
			Name:          "basic audit log entry",
			UID:           "user-1",
			SID:           "session-1",
			TokenID:       "X-Y-Z",
			OrgID:         "org-1",
			TraceID:       "trace-id-1",
			Origin:        "0.0.0.0",
			Category:      v2.NewAuditCategoryV2FromDataLoad(v2.DataLoad{}),
			AuditName:     "AUDITED_ACTION_NAME",
			AuditResult:   audit3log.AuditResultSuccess,
			RequestFields: map[string]any{"requestKey": "requestValue"},
			ResultFields:  map[string]any{"resultKey": "resultValue"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("auditLogV3"),
					"auditLogV3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"uid":     objmatcher.NewEqualsMatcher("user-1"),
						"sid":     objmatcher.NewEqualsMatcher("session-1"),
						"tokenId": objmatcher.NewEqualsMatcher("X-Y-Z"),
						"orgId":   objmatcher.NewEqualsMatcher("org-1"),
						"traceId": objmatcher.NewEqualsMatcher("trace-id-1"),
						"origin":  objmatcher.NewEqualsMatcher("0.0.0.0"),
						"categories": objmatcher.SliceMatcher{
							objmatcher.NewEqualsMatcher("dataLoad"),
						},
						"name":       objmatcher.NewEqualsMatcher("AUDITED_ACTION_NAME"),
						"result":     objmatcher.NewEqualsMatcher(resultValue),
						"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
						"requestFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"loadedResources": objmatcher.NewEqualsMatcher(nil),
							"requestKey":      objmatcher.NewEqualsMatcher("requestValue"),
						}),
						"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"resultKey": objmatcher.NewEqualsMatcher("resultValue"),
						}),
						"type": objmatcher.NewEqualsMatcher(typeValue),
					}),
				}),
			},
		},
	}
}

func Audit3LogJSONTestSuite(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) audit3log.Logger) {
	audit3LogJSONOutputTests(t, entityName, entityVersion, loggerProvider)
	rFieldIsntOverwrittenByRFieldsTest(t, entityName, entityVersion, loggerProvider)
	extraRFieldsDoNotAppear(t, entityName, entityVersion, loggerProvider)
}

func audit3LogJSONOutputTests(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) audit3log.Logger) {
	for i, tc := range Audit3TestCases(entityName, entityVersion) {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Audit(
				tc.AuditName,
				tc.AuditResult,
				audit3log.UID(tc.UID),
				audit3log.SID(tc.SID),
				audit3log.TokenID(tc.TokenID),
				audit3log.OrgID(tc.OrgID),
				audit3log.TraceID(tc.TraceID),
				audit3log.Origin(tc.Origin),
				audit3log.Category(tc.Category),
				audit3log.RequestFields(tc.RequestFields),
				audit3log.ResultFields(tc.ResultFields))

			gotAuditLog := map[string]any{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotAuditLog)
			require.NoError(t, err, "Case %d: %s\nAudit log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotAuditLog), "Case %d: %s", i, tc.Name)
		})
	}
}

// Verifies that if different parameters are specified using ResultField/RequestField and ResultFields/RequestFields,
// all of the values are present in the final output (that is, these parameters should be additive).
func rFieldIsntOverwrittenByRFieldsTest(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) audit3log.Logger) {
	mapFieldMatcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"key1": objmatcher.NewEqualsMatcher("val1"),
		"key2": objmatcher.NewEqualsMatcher("val2"),
	})
	for i, tc := range []struct {
		name   string
		params []audit3log.Param
		want   objmatcher.MapMatcher
	}{
		{
			name: "ResultField params are additive",
			params: []audit3log.Param{
				audit3log.ResultField("key1", "val1"),
				audit3log.ResultFields(map[string]any{"key2": "val2"}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("auditLogV3"),
					"auditLogV3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"time":         objmatcher.NewRegExpMatcher(".+"),
						"name":         objmatcher.NewEqualsMatcher("audited action name"),
						"type":         objmatcher.NewEqualsMatcher(typeValue),
						"result":       objmatcher.NewEqualsMatcher(resultValue),
						"logEntryId":   objmatcher.NewRegExpMatcher(uuidRegexpValue),
						"resultFields": mapFieldMatcher,
					}),
				}),
			}),
		},
		{
			name: "RequestField params are additive",
			params: []audit3log.Param{
				audit3log.RequestField("key1", "val1"),
				audit3log.RequestFields(map[string]any{
					"key2": "val2",
				}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("auditLogV3"),
					"auditLogV3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"time":          objmatcher.NewRegExpMatcher(".+"),
						"name":          objmatcher.NewEqualsMatcher("audited action name"),
						"type":          objmatcher.NewEqualsMatcher(typeValue),
						"result":        objmatcher.NewEqualsMatcher(resultValue),
						"logEntryId":    objmatcher.NewRegExpMatcher(uuidRegexpValue),
						"requestFields": mapFieldMatcher,
					}),
				}),
			}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			logger.Audit("audited action name", audit3log.AuditResultSuccess, tc.params...)

			auditLog := map[string]any{}
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
func extraRFieldsDoNotAppear(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) audit3log.Logger) {
	const (
		resultFieldsKey  = "resultFields"
		requestFieldsKey = "requestFields"
	)

	for i, tc := range []struct {
		name       string
		paramKey   string
		paramFunc  func(key string, val any) audit3log.Param
		paramsFunc func(map[string]any) audit3log.Param
	}{
		{
			name:       "Params stay separate across calls for ResultParam",
			paramKey:   resultFieldsKey,
			paramFunc:  audit3log.ResultField,
			paramsFunc: audit3log.ResultFields,
		},
		{
			name:       "Params stay separate across calls for RequestParam",
			paramKey:   requestFieldsKey,
			paramFunc:  audit3log.RequestField,
			paramsFunc: audit3log.RequestFields,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			reusedParams := tc.paramsFunc(map[string]any{"key1": "val1"})

			logger.Audit(
				"audited action name",
				audit3log.AuditResultSuccess,
				reusedParams,
				tc.paramFunc("key2", "val2"))
			want := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("auditLogV3"),
					"auditLogV3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"time":       objmatcher.NewRegExpMatcher(".+"),
						"name":       objmatcher.NewEqualsMatcher("audited action name"),
						"type":       objmatcher.NewEqualsMatcher(typeValue),
						"result":     objmatcher.NewEqualsMatcher(resultValue),
						"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
						tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key1": objmatcher.NewEqualsMatcher("val1"),
							"key2": objmatcher.NewEqualsMatcher("val2"),
						}),
					}),
				}),
			})
			auditLog := map[string]any{}
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
			logger.Audit("audited action name", audit3log.AuditResultSuccess, reusedParams)

			want = map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("auditLogV3"),
					"auditLogV3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"time":       objmatcher.NewRegExpMatcher(".+"),
						"name":       objmatcher.NewEqualsMatcher("audited action name"),
						"type":       objmatcher.NewEqualsMatcher(typeValue),
						"result":     objmatcher.NewEqualsMatcher(resultValue),
						"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
						tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key1": objmatcher.NewEqualsMatcher("val1"),
						}),
					}),
				}),
			}

			auditLog = map[string]any{}
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
