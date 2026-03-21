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

package audit3logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	v2 "github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/category/v2"
	commonv2 "github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common/v2"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	TestCaseName string

	Deployment     string
	Host           string
	Product        string
	ProductVersion string
	Stack          string
	Service        string
	Environment    string
	ProducerType   audit3log.AuditProducerType
	Organizations  []audit3log.Organization
	EventID        string
	LogEntryID     *string
	UserAgent      string
	Category       v2.AuditCategoryV2
	Entities       []any
	Users          []audit3log.ContextualizedUser
	Origins        []string
	SourceOrigin   string
	RequestFields  map[string]any
	ResultFields   map[string]any
	UID            string
	SID            string
	TokenID        string
	OrgID          string
	TraceID        string
	Origin         string
	Name           string
	Result         audit3log.AuditResultType
	JSONMatcher    objmatcher.MapMatcher

	AdditionalParams []audit3log.Param
}

func (tc TestCase) Params() []audit3log.Param {
	params := []audit3log.Param{
		audit3log.Deployment(tc.Deployment),
		audit3log.Host(tc.Host),
		audit3log.Product(tc.Product),
		audit3log.ProductVersion(tc.ProductVersion),
		audit3log.Stack(tc.Stack),
		audit3log.Service(tc.Service),
		audit3log.Environment(tc.Environment),
		audit3log.ProducerType(tc.ProducerType),
		audit3log.Organizations(tc.Organizations),
		audit3log.EventID(tc.EventID),
		audit3log.UserAgent(tc.UserAgent),
		audit3log.Entities(tc.Entities),
		audit3log.Users(tc.Users),
		audit3log.Origins(tc.Origins),
		audit3log.SourceOrigin(tc.SourceOrigin),
		audit3log.RequestFields(tc.RequestFields),
		audit3log.ResultFields(tc.ResultFields),
		audit3log.UID(tc.UID),
		audit3log.SID(tc.SID),
		audit3log.TokenID(tc.TokenID),
		audit3log.OrgID(tc.OrgID),
		audit3log.TraceID(tc.TraceID),
		audit3log.Origin(tc.Origin),
		audit3log.Category(tc.Category),
	}
	if tc.LogEntryID != nil {
		params = append(params, audit3log.LogEntryID(*tc.LogEntryID))
	}
	params = append(params, tc.AdditionalParams...)
	return params
}

func TestCases() []TestCase {
	return []TestCase{
		{
			TestCaseName: "basic audit log entry",

			Deployment:     "test-deployment",
			Host:           "test-host",
			Product:        "test-product",
			ProductVersion: "test-product-version",
			Stack:          "test-stack",
			Service:        "test-service",
			Environment:    "test-environment",
			ProducerType:   audit3log.AuditProducerServer,
			Organizations: []audit3log.Organization{
				{
					ID:     "d460d218-5768-43cb-888c-ebd27637e9a6",
					Reason: "test-reason-1",
				},
				{
					ID:     "851aa171-b783-4c11-8cbb-ed5ef31a9ac5",
					Reason: "test-reason-2",
				},
			},
			EventID:    "c15487b9-ff6a-4bb1-8c25-2433a185c438",
			LogEntryID: new("46ac025f-5d70-4a79-8e8e-270da9635a43"),
			UserAgent:  "test-user-agent",
			Category:   v2.NewAuditCategoryV2FromDataLoad(v2.DataLoad{}),
			Entities: []any{
				"test-entity-1",
				2,
				map[string]any{
					"key-1": "value-1",
					"key-2": 2,
				},
			},
			Users: []audit3log.ContextualizedUser{
				{
					UID:       "test-user-id",
					UserName:  new("test-username"),
					FirstName: new("test-firstname"),
					LastName:  new("test-lastname"),
					Groups: []string{
						"test-group-1",
						"test-group-2",
					},
					Realm: new("test-realm"),
				},
			},
			Origins: []string{
				"test-origin-1",
				"test-origin-2",
			},
			SourceOrigin: "test-source-origin",
			RequestFields: map[string]any{
				"key-1": "value-1",
				"key-2": 2,
			},
			ResultFields: map[string]any{
				"test-result-fields-key-1": "test-result-fields-value-1",
				"test-result-fields-key-2": 2,
				"test-result-fields-key-3": map[string]any{
					"key-1": "value-1",
					"key-2": 2,
				},
			},
			UID:     "test-uid",
			SID:     "test-sid",
			TokenID: "test-token-id",
			OrgID:   "91de1891-387a-405d-8a33-8543af87afbd",
			TraceID: "test-trace-id",
			Origin:  "0.0.0.0",
			Name:    "TEST_AUDITED_ACTION_NAME",
			Result:  audit3log.AuditResultSuccess,

			JSONMatcher: map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"deployment":     objmatcher.NewEqualsMatcher("test-deployment"),
				"host":           objmatcher.NewEqualsMatcher("test-host"),
				"product":        objmatcher.NewEqualsMatcher("test-product"),
				"productVersion": objmatcher.NewEqualsMatcher("test-product-version"),
				"stack":          objmatcher.NewEqualsMatcher("test-stack"),
				"service":        objmatcher.NewEqualsMatcher("test-service"),
				"environment":    objmatcher.NewEqualsMatcher("test-environment"),
				"producerType":   objmatcher.NewEqualsMatcher(string(audit3log.AuditProducerServer)),
				"organizations": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("d460d218-5768-43cb-888c-ebd27637e9a6"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-1"),
					}),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("851aa171-b783-4c11-8cbb-ed5ef31a9ac5"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-2"),
					}),
				}),
				"eventId":    objmatcher.NewEqualsMatcher("c15487b9-ff6a-4bb1-8c25-2433a185c438"),
				"logEntryId": objmatcher.NewEqualsMatcher("46ac025f-5d70-4a79-8e8e-270da9635a43"),
				"userAgent":  objmatcher.NewEqualsMatcher("test-user-agent"),
				"categories": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("dataLoad"),
				}),
				"entities": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-entity-1"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"key-1": objmatcher.NewEqualsMatcher("value-1"),
						"key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					}),
				}),
				"users": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid":       objmatcher.NewEqualsMatcher("test-user-id"),
						"userName":  objmatcher.NewEqualsMatcher("test-username"),
						"firstName": objmatcher.NewEqualsMatcher("test-firstname"),
						"lastName":  objmatcher.NewEqualsMatcher("test-lastname"),
						"groups": objmatcher.SliceMatcher([]objmatcher.Matcher{
							objmatcher.NewEqualsMatcher("test-group-1"),
							objmatcher.NewEqualsMatcher("test-group-2"),
						}),
						"realm": objmatcher.NewEqualsMatcher("test-realm"),
					}),
				}),
				"origins": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-origin-1"),
					objmatcher.NewEqualsMatcher("test-origin-2"),
				}),
				"sourceOrigin": objmatcher.NewEqualsMatcher("test-source-origin"),
				"requestFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"loadedResources": objmatcher.NewEqualsMatcher(nil),
					"key-1":           objmatcher.NewEqualsMatcher("value-1"),
					"key-2":           objmatcher.NewEqualsMatcher(json.Number("2")),
				}),
				"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"test-result-fields-key-1": objmatcher.NewEqualsMatcher("test-result-fields-value-1"),
					"test-result-fields-key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					"test-result-fields-key-3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"key-1": objmatcher.NewEqualsMatcher("value-1"),
						"key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					}),
				}),
				"uid":     objmatcher.NewEqualsMatcher("test-uid"),
				"sid":     objmatcher.NewEqualsMatcher("test-sid"),
				"tokenId": objmatcher.NewEqualsMatcher("test-token-id"),
				"orgId":   objmatcher.NewEqualsMatcher("91de1891-387a-405d-8a33-8543af87afbd"),
				"traceId": objmatcher.NewEqualsMatcher("test-trace-id"),
				"origin":  objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":    objmatcher.NewEqualsMatcher("TEST_AUDITED_ACTION_NAME"),
				"result":  objmatcher.NewEqualsMatcher(string(audit3log.AuditResultSuccess)),
				"type":    objmatcher.NewEqualsMatcher("audit.3"),
			},
		},
		{
			TestCaseName: "audit log entry sets LogEntryID to random UUID if not specified",

			Deployment:     "test-deployment",
			Host:           "test-host",
			Product:        "test-product",
			ProductVersion: "test-product-version",
			Stack:          "test-stack",
			Service:        "test-service",
			Environment:    "test-environment",
			ProducerType:   audit3log.AuditProducerServer,
			Organizations: []audit3log.Organization{
				{
					ID:     "d460d218-5768-43cb-888c-ebd27637e9a6",
					Reason: "test-reason-1",
				},
				{
					ID:     "851aa171-b783-4c11-8cbb-ed5ef31a9ac5",
					Reason: "test-reason-2",
				},
			},
			EventID:   "c15487b9-ff6a-4bb1-8c25-2433a185c438",
			UserAgent: "test-user-agent",
			Category:  v2.NewAuditCategoryV2FromDataLoad(v2.DataLoad{}),
			Entities: []any{
				"test-entity-1",
				2,
				map[string]any{
					"key-1": "value-1",
					"key-2": 2,
				},
			},
			Users: []audit3log.ContextualizedUser{
				{
					UID:       "test-user-id",
					UserName:  new("test-username"),
					FirstName: new("test-firstname"),
					LastName:  new("test-lastname"),
					Groups: []string{
						"test-group-1",
						"test-group-2",
					},
					Realm: new("test-realm"),
				},
			},
			Origins: []string{
				"test-origin-1",
				"test-origin-2",
			},
			SourceOrigin: "test-source-origin",
			RequestFields: map[string]any{
				"key-1": "value-1",
				"key-2": 2,
			},
			ResultFields: map[string]any{
				"test-result-fields-key-1": "test-result-fields-value-1",
				"test-result-fields-key-2": 2,
				"test-result-fields-key-3": map[string]any{
					"key-1": "value-1",
					"key-2": 2,
				},
			},
			UID:     "test-uid",
			SID:     "test-sid",
			TokenID: "test-token-id",
			OrgID:   "91de1891-387a-405d-8a33-8543af87afbd",
			TraceID: "test-trace-id",
			Origin:  "0.0.0.0",
			Name:    "TEST_AUDITED_ACTION_NAME",
			Result:  audit3log.AuditResultSuccess,

			JSONMatcher: map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"deployment":     objmatcher.NewEqualsMatcher("test-deployment"),
				"host":           objmatcher.NewEqualsMatcher("test-host"),
				"product":        objmatcher.NewEqualsMatcher("test-product"),
				"productVersion": objmatcher.NewEqualsMatcher("test-product-version"),
				"stack":          objmatcher.NewEqualsMatcher("test-stack"),
				"service":        objmatcher.NewEqualsMatcher("test-service"),
				"environment":    objmatcher.NewEqualsMatcher("test-environment"),
				"producerType":   objmatcher.NewEqualsMatcher(string(audit3log.AuditProducerServer)),
				"organizations": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("d460d218-5768-43cb-888c-ebd27637e9a6"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-1"),
					}),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("851aa171-b783-4c11-8cbb-ed5ef31a9ac5"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-2"),
					}),
				}),
				"eventId":    objmatcher.NewEqualsMatcher("c15487b9-ff6a-4bb1-8c25-2433a185c438"),
				"logEntryId": objmatcher.NewRegExpMatcher("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$"),
				"userAgent":  objmatcher.NewEqualsMatcher("test-user-agent"),
				"categories": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("dataLoad"),
				}),
				"entities": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-entity-1"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"key-1": objmatcher.NewEqualsMatcher("value-1"),
						"key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					}),
				}),
				"users": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid":       objmatcher.NewEqualsMatcher("test-user-id"),
						"userName":  objmatcher.NewEqualsMatcher("test-username"),
						"firstName": objmatcher.NewEqualsMatcher("test-firstname"),
						"lastName":  objmatcher.NewEqualsMatcher("test-lastname"),
						"groups": objmatcher.SliceMatcher([]objmatcher.Matcher{
							objmatcher.NewEqualsMatcher("test-group-1"),
							objmatcher.NewEqualsMatcher("test-group-2"),
						}),
						"realm": objmatcher.NewEqualsMatcher("test-realm"),
					}),
				}),
				"origins": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-origin-1"),
					objmatcher.NewEqualsMatcher("test-origin-2"),
				}),
				"sourceOrigin": objmatcher.NewEqualsMatcher("test-source-origin"),
				"requestFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"loadedResources": objmatcher.NewEqualsMatcher(nil),
					"key-1":           objmatcher.NewEqualsMatcher("value-1"),
					"key-2":           objmatcher.NewEqualsMatcher(json.Number("2")),
				}),
				"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"test-result-fields-key-1": objmatcher.NewEqualsMatcher("test-result-fields-value-1"),
					"test-result-fields-key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					"test-result-fields-key-3": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"key-1": objmatcher.NewEqualsMatcher("value-1"),
						"key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					}),
				}),
				"uid":     objmatcher.NewEqualsMatcher("test-uid"),
				"sid":     objmatcher.NewEqualsMatcher("test-sid"),
				"tokenId": objmatcher.NewEqualsMatcher("test-token-id"),
				"orgId":   objmatcher.NewEqualsMatcher("91de1891-387a-405d-8a33-8543af87afbd"),
				"traceId": objmatcher.NewEqualsMatcher("test-trace-id"),
				"origin":  objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":    objmatcher.NewEqualsMatcher("TEST_AUDITED_ACTION_NAME"),
				"result":  objmatcher.NewEqualsMatcher(string(audit3log.AuditResultSuccess)),
				"type":    objmatcher.NewEqualsMatcher("audit.3"),
			},
		},
		{
			TestCaseName: "audit log entry sets all appended values for multi-value fields",

			Deployment:     "test-deployment",
			Host:           "test-host",
			Product:        "test-product",
			ProductVersion: "test-product-version",
			Stack:          "test-stack",
			Service:        "test-service",
			Environment:    "test-environment",
			ProducerType:   audit3log.AuditProducerServer,
			Organizations: []audit3log.Organization{
				{
					ID:     "d460d218-5768-43cb-888c-ebd27637e9a6",
					Reason: "test-reason-1",
				},
			},
			EventID:   "c15487b9-ff6a-4bb1-8c25-2433a185c438",
			UserAgent: "test-user-agent",
			Category:  v2.NewAuditCategoryV2FromDataLoad(v2.DataLoad{}),
			Entities: []any{
				"test-entity-1",
			},
			Users: []audit3log.ContextualizedUser{
				{
					UID: "test-user-id-1",
				},
			},
			Origins: []string{
				"test-origin-1",
			},
			SourceOrigin: "test-source-origin",
			RequestFields: map[string]any{
				"key-1": "value-1",
			},
			ResultFields: map[string]any{
				"test-result-fields-key-1": "test-result-fields-value-1",
			},
			UID:     "test-uid",
			SID:     "test-sid",
			TokenID: "test-token-id",
			OrgID:   "91de1891-387a-405d-8a33-8543af87afbd",
			TraceID: "test-trace-id",
			Origin:  "0.0.0.0",
			Name:    "TEST_AUDITED_ACTION_NAME",
			Result:  audit3log.AuditResultSuccess,

			AdditionalParams: []audit3log.Param{
				audit3log.Organizations([]audit3log.Organization{
					{
						ID:     "851aa171-b783-4c11-8cbb-ed5ef31a9ac5",
						Reason: "test-reason-2",
					},
				}),
				audit3log.Category(v2.NewAuditCategoryV2FromContainerSearch(v2.ContainerSearch{})),
				audit3log.Entities([]any{
					2,
					map[string]any{
						"key-1": "value-1",
						"key-2": 2,
					},
				}),
				audit3log.Users([]audit3log.ContextualizedUser{
					{
						UID: "test-user-id-2",
					},
				}),
				audit3log.Origins([]string{
					"test-origin-2",
				}),
				audit3log.RequestFields(map[string]any{
					"key-2": "value-2",
				}),
				audit3log.ResultFields(map[string]any{
					"test-result-fields-key-2": "test-result-fields-value-2",
				}),
			},

			JSONMatcher: map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"deployment":     objmatcher.NewEqualsMatcher("test-deployment"),
				"host":           objmatcher.NewEqualsMatcher("test-host"),
				"product":        objmatcher.NewEqualsMatcher("test-product"),
				"productVersion": objmatcher.NewEqualsMatcher("test-product-version"),
				"stack":          objmatcher.NewEqualsMatcher("test-stack"),
				"service":        objmatcher.NewEqualsMatcher("test-service"),
				"environment":    objmatcher.NewEqualsMatcher("test-environment"),
				"producerType":   objmatcher.NewEqualsMatcher(string(audit3log.AuditProducerServer)),
				"organizations": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("d460d218-5768-43cb-888c-ebd27637e9a6"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-1"),
					}),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("851aa171-b783-4c11-8cbb-ed5ef31a9ac5"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-2"),
					}),
				}),
				"eventId":    objmatcher.NewEqualsMatcher("c15487b9-ff6a-4bb1-8c25-2433a185c438"),
				"logEntryId": objmatcher.NewRegExpMatcher("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$"),
				"userAgent":  objmatcher.NewEqualsMatcher("test-user-agent"),
				"categories": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("dataLoad"),
					objmatcher.NewEqualsMatcher("containerSearch"),
				}),
				"entities": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-entity-1"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"key-1": objmatcher.NewEqualsMatcher("value-1"),
						"key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
					}),
				}),
				"users": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid": objmatcher.NewEqualsMatcher("test-user-id-1"),
					}),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid": objmatcher.NewEqualsMatcher("test-user-id-2"),
					}),
				}),
				"origins": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-origin-1"),
					objmatcher.NewEqualsMatcher("test-origin-2"),
				}),
				"sourceOrigin": objmatcher.NewEqualsMatcher("test-source-origin"),
				"requestFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"containerSearchQuery": objmatcher.NewEqualsMatcher(nil),
					"loadedResources":      objmatcher.NewEqualsMatcher(nil),
					"key-1":                objmatcher.NewEqualsMatcher("value-1"),
					"key-2":                objmatcher.NewEqualsMatcher("value-2"),
				}),
				"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"containerSearchResults":   objmatcher.NewEqualsMatcher(nil),
					"test-result-fields-key-1": objmatcher.NewEqualsMatcher("test-result-fields-value-1"),
					"test-result-fields-key-2": objmatcher.NewEqualsMatcher("test-result-fields-value-2"),
				}),
				"uid":     objmatcher.NewEqualsMatcher("test-uid"),
				"sid":     objmatcher.NewEqualsMatcher("test-sid"),
				"tokenId": objmatcher.NewEqualsMatcher("test-token-id"),
				"orgId":   objmatcher.NewEqualsMatcher("91de1891-387a-405d-8a33-8543af87afbd"),
				"traceId": objmatcher.NewEqualsMatcher("test-trace-id"),
				"origin":  objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":    objmatcher.NewEqualsMatcher("TEST_AUDITED_ACTION_NAME"),
				"result":  objmatcher.NewEqualsMatcher(string(audit3log.AuditResultSuccess)),
				"type":    objmatcher.NewEqualsMatcher("audit.3"),
			},
		},
		{
			TestCaseName: "audit log entry sets values extracted from categories",

			Deployment:     "test-deployment",
			Host:           "test-host",
			Product:        "test-product",
			ProductVersion: "test-product-version",
			Stack:          "test-stack",
			Service:        "test-service",
			Environment:    "test-environment",
			ProducerType:   audit3log.AuditProducerServer,
			Organizations: []audit3log.Organization{
				{
					ID:     "d460d218-5768-43cb-888c-ebd27637e9a6",
					Reason: "test-reason-1",
				},
			},
			EventID:   "c15487b9-ff6a-4bb1-8c25-2433a185c438",
			UserAgent: "test-user-agent",
			Category: v2.NewAuditCategoryV2FromTokenRevoke(v2.TokenRevoke{
				// UID should be extracted as a user
				RevokedTokens: []commonv2.Token{
					{
						UserId: new("test-revoked-user"),
					},
				},
			}),
			Entities: []any{
				"test-entity-1",
			},
			Users: []audit3log.ContextualizedUser{
				{
					UID: "test-user-id-1",
				},
			},
			Origins: []string{
				"test-origin-1",
			},
			SourceOrigin: "test-source-origin",
			RequestFields: map[string]any{
				"key-1": "value-1",
			},
			ResultFields: map[string]any{
				"test-result-fields-key-1": "test-result-fields-value-1",
			},
			UID:     "test-uid",
			SID:     "test-sid",
			TokenID: "test-token-id",
			OrgID:   "91de1891-387a-405d-8a33-8543af87afbd",
			TraceID: "test-trace-id",
			Origin:  "0.0.0.0",
			Name:    "TEST_AUDITED_ACTION_NAME",
			Result:  audit3log.AuditResultSuccess,

			AdditionalParams: []audit3log.Param{
				// RID should be extracted as an entity
				audit3log.Category(v2.NewAuditCategoryV2FromDataTransform(v2.DataTransform{
					TransformTargets: []commonv2.DataResource{
						{
							Id: commonv2.NewIdentifierFromRid(rid.MustNew("service", "instance", "resource-type", "00d5fa8a-87d3-4416-ad86-b96cba130b19")),
						},
					},
				})),
			},

			JSONMatcher: map[string]objmatcher.Matcher{
				"time":           objmatcher.NewRegExpMatcher(".+"),
				"deployment":     objmatcher.NewEqualsMatcher("test-deployment"),
				"host":           objmatcher.NewEqualsMatcher("test-host"),
				"product":        objmatcher.NewEqualsMatcher("test-product"),
				"productVersion": objmatcher.NewEqualsMatcher("test-product-version"),
				"stack":          objmatcher.NewEqualsMatcher("test-stack"),
				"service":        objmatcher.NewEqualsMatcher("test-service"),
				"environment":    objmatcher.NewEqualsMatcher("test-environment"),
				"producerType":   objmatcher.NewEqualsMatcher(string(audit3log.AuditProducerServer)),
				"organizations": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"id":     objmatcher.NewEqualsMatcher("d460d218-5768-43cb-888c-ebd27637e9a6"),
						"reason": objmatcher.NewEqualsMatcher("test-reason-1"),
					}),
				}),
				"eventId":    objmatcher.NewEqualsMatcher("c15487b9-ff6a-4bb1-8c25-2433a185c438"),
				"logEntryId": objmatcher.NewRegExpMatcher("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$"),
				"userAgent":  objmatcher.NewEqualsMatcher("test-user-agent"),
				"categories": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("tokenRevoke"),
					objmatcher.NewEqualsMatcher("dataTransform"),
				}),
				"entities": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-entity-1"),
					// extracted from category
					objmatcher.NewEqualsMatcher("ri.service.instance.resource-type.00d5fa8a-87d3-4416-ad86-b96cba130b19"),
				}),
				"users": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid": objmatcher.NewEqualsMatcher("test-user-id-1"),
					}),
					objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"uid": objmatcher.NewEqualsMatcher("test-revoked-user"),
					}),
				}),
				"origins": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-origin-1"),
				}),
				"sourceOrigin": objmatcher.NewEqualsMatcher("test-source-origin"),
				"requestFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"revokeTokensDescription": objmatcher.NewEqualsMatcher(nil),
					"transformDescription":    objmatcher.NewEqualsMatcher(nil),
					"transformTargets": objmatcher.SliceMatcher([]objmatcher.Matcher{
						objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"context": objmatcher.NewEqualsMatcher([]any{}),
							"id": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
								"rid":  objmatcher.NewEqualsMatcher("ri.service.instance.resource-type.00d5fa8a-87d3-4416-ad86-b96cba130b19"),
								"type": objmatcher.NewEqualsMatcher("rid"),
							}),
						}),
					}),
					"key-1": objmatcher.NewEqualsMatcher("value-1"),
				}),
				"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"test-result-fields-key-1": objmatcher.NewEqualsMatcher("test-result-fields-value-1"),
					"revokedTokens": objmatcher.SliceMatcher([]objmatcher.Matcher{
						objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"userId": objmatcher.NewEqualsMatcher("test-revoked-user"),
							"type":   objmatcher.NewEqualsMatcher(""),
						}),
					}),
				}),
				"uid":     objmatcher.NewEqualsMatcher("test-uid"),
				"sid":     objmatcher.NewEqualsMatcher("test-sid"),
				"tokenId": objmatcher.NewEqualsMatcher("test-token-id"),
				"orgId":   objmatcher.NewEqualsMatcher("91de1891-387a-405d-8a33-8543af87afbd"),
				"traceId": objmatcher.NewEqualsMatcher("test-trace-id"),
				"origin":  objmatcher.NewEqualsMatcher("0.0.0.0"),
				"name":    objmatcher.NewEqualsMatcher("TEST_AUDITED_ACTION_NAME"),
				"result":  objmatcher.NewEqualsMatcher(string(audit3log.AuditResultSuccess)),
				"type":    objmatcher.NewEqualsMatcher("audit.3"),
			},
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) audit3log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.TestCaseName, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Audit(
				tc.Name,
				tc.Result,
				tc.Params()...,
			)

			gotAuditLog := map[string]any{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotAuditLog)
			require.NoError(t, err, "Case %d: %s\nAudit log line is not a valid map: %v", i, tc.TestCaseName, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotAuditLog), "Case %d: %s", i, tc.TestCaseName)
		})
	}
}

func DualLoggingJSONTestSuite(t *testing.T, dualLoggerProvider func(audit3Writer, audit2Writer io.Writer) audit3log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.TestCaseName, func(t *testing.T) {
			audit3Buf := &bytes.Buffer{}
			audit2Buf := &bytes.Buffer{}
			logger := dualLoggerProvider(audit3Buf, audit2Buf)

			logger.Audit(
				tc.Name,
				tc.Result,
				tc.Params()...,
			)

			var gotAuditV3Log logging.AuditLogV3
			err := safejson.Unmarshal(audit3Buf.Bytes(), &gotAuditV3Log)
			require.NoError(t, err, "Case %d: %s\nError unmarshalling output as v3 audit log", i, tc.TestCaseName)

			var gotAuditV2Log logging.AuditLogV2
			err = safejson.Unmarshal(audit2Buf.Bytes(), &gotAuditV2Log)
			require.NoError(t, err, "Case %d: %s\nError unmarshalling output as v2 audit log", i, tc.TestCaseName)

			assert.Equal(t, gotAuditV3Log.Time, gotAuditV2Log.Time)
			assertEqualNillableStringPointer(t, gotAuditV3Log.SourceOrigin, gotAuditV2Log.RequestParams, "_sourceOrigin")
			assertEqualNillableStringPointer(t, gotAuditV3Log.TokenId, gotAuditV2Log.RequestParams, "_tokenId")
			assertEqualNillableStringPointer(t, gotAuditV3Log.UserAgent, gotAuditV2Log.RequestParams, "_userAgent")
			assert.Equal(t, gotAuditV3Log.EventId.String(), gotAuditV2Log.RequestParams["_auditEventId"])
			assert.Equal(t, gotAuditV3Log.LogEntryId.String(), gotAuditV2Log.RequestParams["_auditLogEntryId"])
		})
	}
}

func assertEqualNillableStringPointer[T ~string](t *testing.T, want *T, gotMap map[string]any, mapKey string) {
	gotValue, exists := gotMap[mapKey]

	if want == nil {
		assert.Falsef(t, exists, "key %q exists in map when it should not", mapKey)
	} else {
		assert.Equal(t, string(*want), gotValue)
	}
}

//go:fix inline
func toPointer[T any](in T) *T {
	return new(in)
}
