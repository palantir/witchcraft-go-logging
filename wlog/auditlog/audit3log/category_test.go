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

package audit3log

import (
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/rid"
	v2 "github.com/palantir/witchcraft-go-logging/conjure/foundry/audit/api/category/v2"
	commonv2 "github.com/palantir/witchcraft-go-logging/conjure/foundry/audit/api/common/v2"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verifies that the "Category" function returns the appropriate audit logging parameters. This test is a baseline that
// verifies that the "Category" function properly populates fields for a category that includes request and response
// fields: it is not an exhaustive test of every declared category.
func TestCategory(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := FromContext(ctx)

	testRid, err := rid.ParseRID("ri.function-registry.main.function.19f5abc7-c299-4f4d-92e3-b109f1d8795b")
	require.NoError(t, err)

	query := "testQuery"
	categoryParams, err := Category(v2.NewAuditCategoryV2FromAppConfigSearch(v2.AppConfigSearch{
		AppConfigSearchQuery: &query,
		AppConfigSearchResults: []commonv2.ApplicationResource{
			{
				Id:      commonv2.NewIdentifierFromRid(testRid),
				Product: "test-product",
				Context: []commonv2.ResourceContext{
					{
						Value:       "test-value",
						Description: "test-description",
					},
				},
			},
		},
	}))
	require.NoError(t, err)
	logger.Audit("TEST_ENTRY", AuditResultSuccess, categoryParams...)

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))
	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":   objmatcher.NewRegExpMatcher(".+"),
		"type":   objmatcher.NewEqualsMatcher("audit.3"),
		"name":   objmatcher.NewEqualsMatcher("TEST_ENTRY"),
		"result": objmatcher.NewEqualsMatcher("SUCCESS"),
		"requestFields": objmatcher.NewEqualsMatcher(map[string]any{
			"appConfigSearchQuery": "testQuery",
		}),
		"resultFields": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"appConfigSearchResults": objmatcher.SliceMatcher{
				objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"context": objmatcher.SliceMatcher{
						objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"value":       objmatcher.NewEqualsMatcher("test-value"),
							"description": objmatcher.NewEqualsMatcher("test-description"),
						}),
					},
					"id": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"rid":  objmatcher.NewEqualsMatcher("ri.function-registry.main.function.19f5abc7-c299-4f4d-92e3-b109f1d8795b"),
						"type": objmatcher.NewEqualsMatcher("rid"),
					}),
					"product": objmatcher.NewEqualsMatcher("test-product"),
				}),
			},
		}),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}
