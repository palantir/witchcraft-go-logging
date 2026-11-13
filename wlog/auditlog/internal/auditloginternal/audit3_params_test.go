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

package auditloginternal

import (
	"testing"

	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/stretchr/testify/require"
)

func TestAudit3IDs(t *testing.T) {
	const id = "c15487b9-ff6a-4bb1-8c25-2433a185c438"
	for _, tc := range []struct {
		name, eventID, logEntryID string
	}{
		{name: "canonical", eventID: id, logEntryID: id},
		{name: "missing"},
		{name: "invalid", eventID: "invalid", logEntryID: "invalid"},
		{name: "noncanonical", eventID: id, logEntryID: "C15487B9-FF6A-4BB1-8C25-2433A185C438"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := []AuditParam{Audit3EventID(tc.eventID)}
			if tc.logEntryID != "" {
				params = append(params, Audit3LogEntryID(tc.logEntryID))
			}
			log := new(logtypes.AuditLogV3)
			populateAuditLogV3(log, &atLogTimeValues{}, "test", AuditResultSuccess, params, func(p AuditParam) AuditParam { return p })
			require.Equal(t, tc.eventID, log.EventId)
			if tc.logEntryID != "" {
				require.Equal(t, tc.logEntryID, log.LogEntryId)
			} else {
				_, err := uuid.ParseUUID(log.LogEntryId)
				require.NoError(t, err)
			}
		})
	}
}

func TestAudit2RequestCategoryConversion(t *testing.T) {
	for _, value := range []any{"category-name", 42, nil} {
		var singular, plural logtypes.AuditLogV3
		Audit2RequestParam("category", value).Audit3ParamFn(&singular)
		Audit2RequestParams(map[string]any{"category": value}).Audit3ParamFn(&plural)
		require.Equal(t, plural, singular)
	}
}
