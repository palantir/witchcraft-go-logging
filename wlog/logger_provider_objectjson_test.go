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

package wlog_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log/audit3logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log/wrapped1logtests"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

func TestObjectJSONTrc1LogPreservesSpanTags(t *testing.T) {
	logger := trc1log.NewFromCreator(io.Discard, wlog.NewJSONMarshalLoggerProvider().NewLogger)
	span := wtracing.SpanModel{Tags: map[string]string{"key": "value"}}
	logger.Log(span)
	if got := span.Tags["key"]; got != "value" {
		t.Fatalf("span tag after logging = %q, want value", got)
	}
}

func TestObjectJSONAudit3Log(t *testing.T) {
	audit3logtests.JSONConjureTestSuite(t, func(w io.Writer) audit3log.Logger {
		return audit3log.NewFromCreator(w, wlog.NewJSONMarshalLoggerProvider().NewLogger)
	})
}

func TestObjectJSONAudit3DualLog(t *testing.T) {
	audit3logtests.DualLoggingJSONTestSuite(t, func(audit3Writer, audit2Writer io.Writer) audit3log.Logger {
		return audit3log.NewDualLoggerFromCreator(audit3Writer, audit2Writer, wlog.NewJSONMarshalLoggerProvider().NewLogger)
	})
}

func TestObjectJSONWrapped1LogAudit3Log(t *testing.T) {
	entityName := "entity"
	entityVersion := "version"
	wrapped1logtests.Audit3TypedLogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) audit3log.Logger {
		return audit3log.WithParams(
			wrapped1log.NewFromProvider(w, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider(), entityName, entityVersion).AuditV3(),
			audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438"),
		)
	})
}

func TestObjectJSONWrapped1LogAudit3IDs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params []audit3log.Param
		want   string
		typed  bool
	}{
		{name: "generated", typed: true},
		{name: "supplied", params: []audit3log.Param{audit3log.LogEntryID("46ac025f-5d70-4a79-8e8e-270da9635a43")}, want: "46ac025f-5d70-4a79-8e8e-270da9635a43", typed: true},
		{name: "malformed after valid", params: []audit3log.Param{audit3log.LogEntryID("46ac025f-5d70-4a79-8e8e-270da9635a43"), audit3log.LogEntryID("invalid-uuid")}, want: "invalid-uuid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider(), "entity", "version").AuditV3()
			logger.Audit("test", audit3log.AuditResultSuccess, append([]audit3log.Param{audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438")}, tc.params...)...)
			var got struct {
				Payload struct {
					AuditLogV3 struct {
						LogEntryID string `json:"logEntryId"`
					} `json:"auditLogV3"`
				} `json:"payload"`
			}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if tc.want != "" && got.Payload.AuditLogV3.LogEntryID != tc.want {
				t.Fatalf("logEntryId = %q, want %q", got.Payload.AuditLogV3.LogEntryID, tc.want)
			}
			if tc.want == "" {
				if _, err := uuid.ParseUUID(got.Payload.AuditLogV3.LogEntryID); err != nil {
					t.Fatalf("generated logEntryId = %q: %v", got.Payload.AuditLogV3.LogEntryID, err)
				}
			}
			if typed := bytes.HasPrefix(buf.Bytes(), []byte(`{"type":"wrapped.1","payload":`)); typed != tc.typed {
				t.Fatalf("typed emission = %t, want %t: %s", typed, tc.typed, buf.Bytes())
			}
		})
	}
}

func TestObjectJSONWrapped1LogAudit3MissingEventIDFallback(t *testing.T) {
	wrapped1logtests.Audit3LogJSONTestSuite(t, "entity", "version", func(w io.Writer) audit3log.Logger {
		return wrapped1log.NewFromProvider(w, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider(), "entity", "version").AuditV3()
	})

	var buf bytes.Buffer
	logger := wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider(), "entity", "version").AuditV3()
	logger.Audit("test", audit3log.AuditResultSuccess)
	if bytes.HasPrefix(buf.Bytes(), []byte(`{"type":"wrapped.1","payload":`)) {
		t.Fatalf("missing event ID used typed emission: %s", buf.Bytes())
	}
	var got struct {
		Payload struct {
			AuditLogV3 map[string]any `json:"auditLogV3"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if _, hasEventID := got.Payload.AuditLogV3["eventId"]; hasEventID {
		t.Fatalf("fallback emitted invalid event ID: %s", buf.Bytes())
	}
	if _, hasDeployment := got.Payload.AuditLogV3["deployment"]; hasDeployment {
		t.Fatalf("fallback emitted absent deployment: %s", buf.Bytes())
	}
}
