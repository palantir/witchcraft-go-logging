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

package audit3logtests

import (
	"bytes"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name           string
	UID            string
	SID            string
	TokenID        string
	TraceID        string
	OtherUIDs      []string
	Origin         string
	AuditName      string
	AuditResult    audit3log.AuditResultType
	Deployment     string
	Host           string
	Product        string
	ProductVersion string
	RequestParams  map[string]audit3log.AuditSensitivityTaggedValueType
	ResultParams   map[string]audit3log.AuditSensitivityTaggedValueType
	JSONMatcher    objmatcher.MapMatcher
}

func (tc TestCase) Params() []audit3log.Param {
	return []audit3log.Param{
		audit3log.UID(tc.UID),
		audit3log.SID(tc.SID),
		audit3log.TokenID(tc.TokenID),
		audit3log.TraceID(tc.TraceID),
		audit3log.OtherUIDs(tc.OtherUIDs...),
		audit3log.Origin(tc.Origin),
		audit3log.RequestParams(tc.RequestParams),
		audit3log.ResultParams(tc.ResultParams),
	}
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name:           "basic audit log entry",
			UID:            "user-1",
			SID:            "session-1",
			TokenID:        "X-Y-Z",
			TraceID:        "trace-id-1",
			OtherUIDs:      []string{"user-2", "user-3"},
			Origin:         "0.0.0.0",
			AuditName:      "AUDITED_ACTION_NAME",
			AuditResult:    audit3log.AuditResultSuccess,
			Deployment:     "deployment-1",
			Host:           "host-1",
			Product:        "product-1",
			ProductVersion: "1.0.0",
			RequestParams: map[string]audit3log.AuditSensitivityTaggedValueType{
				"requestKey": {Level: []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput}, Payload: "requestValue"},
			},
			ResultParams: map[string]audit3log.AuditSensitivityTaggedValueType{
				"resultKey": {Level: []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData}, Payload: "resultValue"},
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"uid":            objmatcher.NewEqualsMatcher("user-1"),
				"sid":            objmatcher.NewEqualsMatcher("session-1"),
				"tokenId":        objmatcher.NewEqualsMatcher("X-Y-Z"),
				"traceId":        objmatcher.NewEqualsMatcher("trace-id-1"),
				"otherUids":      objmatcher.NewEqualsMatcher([]interface{}{"user-2", "user-3"}),
				"origin":         objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":           objmatcher.NewEqualsMatcher("AUDITED_ACTION_NAME"),
				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
				"host":           objmatcher.NewEqualsMatcher("host-1"),
				"product":        objmatcher.NewEqualsMatcher("product-1"),
				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
				"requestParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"requestKey": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher([]interface{}{"UserInput"}),
						"payload": objmatcher.NewEqualsMatcher("requestValue"),
					}),
				}),
				"resultParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"resultKey": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher([]interface{}{"Data"}),
						"payload": objmatcher.NewEqualsMatcher("resultValue"),
					}),
				}),
				"type": objmatcher.NewEqualsMatcher("audit.3"),
			}),
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
	jsonOutputTests(t, loggerProvider)
	rParamIsntOverwrittenByRParamsTest(t, loggerProvider)
	// extraRParamsDoNotAppear(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Audit(
				tc.AuditName,
				tc.AuditResult,
				tc.Deployment,
				tc.Host,
				tc.Product,
				tc.ProductVersion,
				audit3log.UID(tc.UID),
				audit3log.SID(tc.SID),
				audit3log.TokenID(tc.TokenID),
				audit3log.TraceID(tc.TraceID),
				audit3log.OtherUIDs(tc.OtherUIDs...),
				audit3log.Origin(tc.Origin),
				audit3log.RequestParams(tc.RequestParams),
				audit3log.ResultParams(tc.ResultParams))

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
func rParamIsntOverwrittenByRParamsTest(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
	mapFieldMatcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"key1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"level":   objmatcher.NewEqualsMatcher([]interface{}{"UserInput"}),
			"payload": objmatcher.NewEqualsMatcher("val1"),
		}),
		"key2": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"level":   objmatcher.NewEqualsMatcher([]interface{}{"Data"}),
			"payload": objmatcher.NewEqualsMatcher("val2"),
		}),
	})
	for i, tc := range []struct {
		name   string
		params []audit3log.Param
		want   objmatcher.MapMatcher
	}{
		{
			name: "ResultParam params are additive",
			params: []audit3log.Param{
				audit3log.ResultParam(
					"key1", audit3log.AuditSensitivityTaggedValueType{
						Level:   []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput},
						Payload: "val1",
					},
				),
				audit3log.ResultParams(map[string]audit3log.AuditSensitivityTaggedValueType{
					"key2": {Level: []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData}, Payload: "val2"},
				}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"name":           objmatcher.NewEqualsMatcher("audited action name"),
				"type":           objmatcher.NewEqualsMatcher("audit.3"),
				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
				"host":           objmatcher.NewEqualsMatcher("host-1"),
				"product":        objmatcher.NewEqualsMatcher("product-1"),
				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
				"resultParams":   mapFieldMatcher,
			}),
		},
		{
			name: "RequestParam params are additive",
			params: []audit3log.Param{
				audit3log.RequestParam(
					"key1", audit3log.AuditSensitivityTaggedValueType{
						Level:   []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput},
						Payload: "val1",
					},
				),
				audit3log.RequestParams(map[string]audit3log.AuditSensitivityTaggedValueType{
					"key2": {Level: []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData}, Payload: "val2"},
				}),
			},
			want: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"name":           objmatcher.NewEqualsMatcher("audited action name"),
				"type":           objmatcher.NewEqualsMatcher("audit.3"),
				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
				"host":           objmatcher.NewEqualsMatcher("host-1"),
				"product":        objmatcher.NewEqualsMatcher("product-1"),
				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
				"requestParams":  mapFieldMatcher,
			}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			logger.Audit(
				"audited action name",
				audit3log.AuditResultSuccess,
				"deployment-1",
				"host-1",
				"product-1",
				"1.0.0",
				tc.params...,
			)

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

// // Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// // are modified by making a logger call).
// func extraRParamsDoNotAppear(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
// 	const (
// 		resultParamsKey  = "resultParams"
// 		requestParamsKey = "requestParams"
// 	)

// 	for i, tc := range []struct {
// 		name       string
// 		paramKey   string
// 		paramFunc  func(key string, val interface{}) audit3log.Param
// 		paramsFunc func(map[string]interface{}) audit3log.Param
// 	}{
// 		{
// 			name:       "Params stay separate across calls for ResultParam",
// 			paramKey:   resultParamsKey,
// 			paramFunc:  audit3log.ResultParam,
// 			paramsFunc: audit3log.ResultParams,
// 		},
// 		{
// 			name:       "Params stay separate across calls for RequestParam",
// 			paramKey:   requestParamsKey,
// 			paramFunc:  audit3log.RequestParam,
// 			paramsFunc: audit3log.RequestParams,
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var buf bytes.Buffer
// 			logger := loggerProvider(&buf)

// 			reusedParams := tc.paramsFunc(map[string]interface{}{"key1": "val1"})

// 			logger.Audit(
// 				"audited action name",
// 				audit3log.AuditResultSuccess,
// 				"deployment-1",
// 				"host-1",
// 				"product-1",
// 				"1.0.0",
// 				reusedParams,
// 				tc.paramFunc("key2", "val2"))
// 			want := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
// 				"time":           objmatcher.NewRegExpMatcher(".+"),
// 				"name":           objmatcher.NewEqualsMatcher("audited action name"),
// 				"type":           objmatcher.NewEqualsMatcher("audit.3"),
// 				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
// 				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
// 				"host":           objmatcher.NewEqualsMatcher("host-1"),
// 				"product":        objmatcher.NewEqualsMatcher("product-1"),
// 				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
// 				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
// 					"key1": objmatcher.NewEqualsMatcher("val1"),
// 					"key2": objmatcher.NewEqualsMatcher("val2"),
// 				}),
// 			})
// 			auditLog := map[string]interface{}{}
// 			logEntry := buf.Bytes()
// 			err := json.Unmarshal(logEntry, &auditLog)
// 			require.NoError(
// 				t,
// 				err,
// 				"Case %d: %s\nAudit log is not a valid map: %v",
// 				i,
// 				tc.name,
// 				string(logEntry))
// 			assert.NoError(t, want.Matches(auditLog), "Case %d: %s", i, tc.name)

// 			buf.Reset()
// 			logger.Audit(
// 				"audited action name",
// 				audit3log.AuditResultSuccess,
// 				"deployment-1",
// 				"host-1",
// 				"product-1",
// 				"1.0.0",
// 				reusedParams,
// 			)

// 			want = objmatcher.MapMatcher(map[string]objmatcher.Matcher{
// 				"time":           objmatcher.NewRegExpMatcher(".+"),
// 				"name":           objmatcher.NewEqualsMatcher("audited action name"),
// 				"type":           objmatcher.NewEqualsMatcher("audit.3"),
// 				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
// 				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
// 				"host":           objmatcher.NewEqualsMatcher("host-1"),
// 				"product":        objmatcher.NewEqualsMatcher("product-1"),
// 				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
// 				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
// 					"key1": objmatcher.NewEqualsMatcher("val1"),
// 				}),
// 			})

// 			auditLog = map[string]interface{}{}
// 			logEntry = buf.Bytes()
// 			err = json.Unmarshal(logEntry, &auditLog)
// 			require.NoError(
// 				t,
// 				err,
// 				"Case %d: %s\nAudit log is not a valid map: %v",
// 				i,
// 				tc.name,
// 				string(logEntry))
// 			assert.NoError(t, want.Matches(auditLog), "Case %d: %s", i, tc.name)
// 		})
// 	}
// }
