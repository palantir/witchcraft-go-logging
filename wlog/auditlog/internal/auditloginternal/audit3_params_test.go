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

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/stretchr/testify/require"
)

func TestAudit3DirectObjectFallback(t *testing.T) {
	const id = "c15487b9-ff6a-4bb1-8c25-2433a185c438"
	for _, tc := range []struct {
		name, eventID string
		logEntryID    *string
	}{
		{name: "direct", eventID: id},
		{name: "missing event ID"},
		{name: "invalid event ID", eventID: "invalid"},
		{name: "noncanonical log entry ID", eventID: id, logEntryID: new("C15487B9-FF6A-4BB1-8C25-2433A185C438")},
		{name: "invalid log entry ID", eventID: id, logEntryID: new("invalid")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			params := []Audit3Param{Audit3EventID(tc.eventID), {Param: AuditParam{Audit3ParamFn: func(log *logging.AuditLogV3) {
				calls++
				log.Name = "populated"
			}}}}
			if tc.logEntryID != nil {
				params = append(params, Audit3LogEntryID(*tc.logEntryID))
			}
			obj, fallback := Audit3ToDirectObject("test", AuditResultSuccess, params)
			if tc.name == "direct" {
				require.Nil(t, fallback)
			} else {
				require.NotNil(t, fallback)
				entry := wlog.NewMapLogEntry()
				wlog.ApplyParams(entry, []wlog.Param{fallback})
				record := entry.AllValues()
				require.Equal(t, obj.Name, record[Audit3NameKey])
				require.Equal(t, obj.Time.String(), record[wlog.TimeKey])
				if tc.logEntryID != nil {
					require.Equal(t, *tc.logEntryID, record[Audit3LogEntryIDKey])
				} else {
					require.NotNil(t, obj.LogEntryId)
					require.Equal(t, obj.LogEntryId.String(), record[Audit3LogEntryIDKey])
				}
				if tc.eventID == id {
					require.Equal(t, id, record[Audit3EventIDKey])
				} else {
					require.NotContains(t, record, Audit3EventIDKey)
				}
			}
			require.Equal(t, 1, calls)
		})
	}
}

func TestAudit2RequestCategoryConversion(t *testing.T) {
	for _, value := range []any{"category-name", 42, nil} {
		var singular, plural logging.AuditLogV3
		Audit2RequestParam("category", value).Param.Audit3ParamFn(&singular)
		Audit2RequestParams(map[string]any{"category": value}).Param.Audit3ParamFn(&plural)
		require.Equal(t, plural, singular)
	}
}
