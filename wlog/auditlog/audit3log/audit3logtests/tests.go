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
	"encoding/json"
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
	Stack          string
	Service        string
	Environment    string
	ProducerType   audit3log.AuditProducerType
	Organizations  []audit3log.AuditOrganizationType
	EventId        string
	UserAgent      string
	Categories     []string
	Entities       []interface{}
	Users          []audit3log.AuditContextualizedUserType
	Origins        []string
	SourceOrigin   string
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
		audit3log.Stack(tc.Stack),
		audit3log.Service(tc.Service),
		audit3log.Environment(tc.Environment),
		audit3log.ProducerType(tc.ProducerType),
		audit3log.Organizations(tc.Organizations...),
		audit3log.EventID(tc.EventId),
		audit3log.UserAgent(tc.UserAgent),
		audit3log.Categories(tc.Categories...),
		audit3log.Entities(tc.Entities...),
		audit3log.Users(tc.Users...),
		audit3log.Origins(tc.Origins...),
		audit3log.SourceOrigin(tc.SourceOrigin),
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
			Stack:          "stack-1",
			Service:        "service-1",
			Environment:    "environment-1",
			ProducerType:   audit3log.AuditProducerServer,
			Organizations:  []audit3log.AuditOrganizationType{{ID: "organization-1", Reason: "reason"}},
			EventId:        "event-id-1",
			UserAgent:      "user-agent-1",
			Categories:     []string{"DATA_LOAD", "USER_LOGIN"},
			Entities:       []interface{}{"entity-1", "entity-2"},
			Users: []audit3log.AuditContextualizedUserType{{
				UID:       "user-1",
				UserName:  "username",
				FirstName: "User",
				LastName:  "Name",
				Groups:    []string{"group-1", "group-2"},
				Realm:     "realm-1",
			}},
			Origins:      []string{"0.0.0.0", "1.2.3.4"},
			SourceOrigin: "0.0.0.0",
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
				"stack":          objmatcher.NewEqualsMatcher("stack-1"),
				"service":        objmatcher.NewEqualsMatcher("service-1"),
				"environment":    objmatcher.NewEqualsMatcher("environment-1"),
				"producerType":   objmatcher.NewEqualsMatcher("SERVER"),
				"organizations":  objmatcher.NewEqualsMatcher([]interface{}{map[string]interface{}{"id": "organization-1", "reason": "reason"}}),
				"eventId":        objmatcher.NewEqualsMatcher("event-id-1"),
				"userAgent":      objmatcher.NewEqualsMatcher("user-agent-1"),
				"categories":     objmatcher.NewEqualsMatcher([]interface{}{"DATA_LOAD", "USER_LOGIN"}),
				"entities":       objmatcher.NewEqualsMatcher([]interface{}{"entity-1", "entity-2"}),
				"users": objmatcher.NewEqualsMatcher([]interface{}{map[string]interface{}{
					"uid":       "user-1",
					"userName":  "username",
					"firstName": "User",
					"lastName":  "Name",
					"groups":    []interface{}{"group-1", "group-2"},
					"realm":     "realm-1",
				}}),
				"origins":      objmatcher.NewEqualsMatcher([]interface{}{"0.0.0.0", "1.2.3.4"}),
				"sourceOrigin": objmatcher.NewEqualsMatcher("0.0.0.0"),
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
	extraRParamsDoNotAppear(t, loggerProvider)
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
				audit3log.ResultParams(tc.ResultParams),
				audit3log.Stack(tc.Stack),
				audit3log.Service(tc.Service),
				audit3log.Environment(tc.Environment),
				audit3log.ProducerType(tc.ProducerType),
				audit3log.Organizations(tc.Organizations...),
				audit3log.EventID(tc.EventId),
				audit3log.UserAgent(tc.UserAgent),
				audit3log.Categories(tc.Categories...),
				audit3log.Entities(tc.Entities...),
				audit3log.Users(tc.Users...),
				audit3log.Origins(tc.Origins...),
				audit3log.SourceOrigin(tc.SourceOrigin),
			)

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
				audit3log.ResultParam("key1", []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput}, "val1"),
				audit3log.ResultParams(map[string]audit3log.AuditSensitivityTaggedValueType{
					"key2": {
						Level:   []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData},
						Payload: "val2",
					},
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
				audit3log.RequestParam("key1", []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput}, "val1"),
				audit3log.RequestParams(map[string]audit3log.AuditSensitivityTaggedValueType{
					"key2": {
						Level:   []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData},
						Payload: "val2",
					},
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

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraRParamsDoNotAppear(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
	const (
		resultParamsKey  = "resultParams"
		requestParamsKey = "requestParams"
	)

	for i, tc := range []struct {
		name       string
		paramKey   string
		paramFunc  func(key string, levels []audit3log.AuditSensitivityType, payload interface{}) audit3log.Param
		paramsFunc func(map[string]audit3log.AuditSensitivityTaggedValueType) audit3log.Param
	}{
		{
			name:       "Params stay separate across calls for ResultParam",
			paramKey:   resultParamsKey,
			paramFunc:  audit3log.ResultParam,
			paramsFunc: audit3log.ResultParams,
		},
		{
			name:       "Params stay separate across calls for RequestParam",
			paramKey:   requestParamsKey,
			paramFunc:  audit3log.RequestParam,
			paramsFunc: audit3log.RequestParams,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := loggerProvider(&buf)

			reusedParams := tc.paramsFunc(map[string]audit3log.AuditSensitivityTaggedValueType{"key1": {
				Level:   []audit3log.AuditSensitivityType{audit3log.AuditSensitivityUserInput},
				Payload: "val1",
			}})

			logger.Audit(
				"audited action name",
				audit3log.AuditResultSuccess,
				"deployment-1",
				"host-1",
				"product-1",
				"1.0.0",
				reusedParams,
				tc.paramFunc("key2", []audit3log.AuditSensitivityType{audit3log.AuditSensitivityData}, "val2"))
			want := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"name":           objmatcher.NewEqualsMatcher("audited action name"),
				"type":           objmatcher.NewEqualsMatcher("audit.3"),
				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
				"host":           objmatcher.NewEqualsMatcher("host-1"),
				"product":        objmatcher.NewEqualsMatcher("product-1"),
				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"key1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher([]interface{}{"UserInput"}),
						"payload": objmatcher.NewEqualsMatcher("val1"),
					}),
					"key2": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher([]interface{}{"Data"}),
						"payload": objmatcher.NewEqualsMatcher("val2"),
					}),
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
			logger.Audit(
				"audited action name",
				audit3log.AuditResultSuccess,
				"deployment-1",
				"host-1",
				"product-1",
				"1.0.0",
				reusedParams,
			)

			want = objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"name":           objmatcher.NewEqualsMatcher("audited action name"),
				"type":           objmatcher.NewEqualsMatcher("audit.3"),
				"result":         objmatcher.NewEqualsMatcher("SUCCESS"),
				"deployment":     objmatcher.NewEqualsMatcher("deployment-1"),
				"host":           objmatcher.NewEqualsMatcher("host-1"),
				"product":        objmatcher.NewEqualsMatcher("product-1"),
				"productVersion": objmatcher.NewEqualsMatcher("1.0.0"),
				tc.paramKey: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"key1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher([]interface{}{"UserInput"}),
						"payload": objmatcher.NewEqualsMatcher("val1"),
					}),
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
