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
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/objecttest"
	"github.com/stretchr/testify/require"
)

// omitted lists generated fields that no parameter sets, so the log types leave them out.
var omitted = map[string][]string{
	"serviceLogV1": {"thread", "safe"},
	"metricLogV1":  {"samples"},
	"auditLogV3":   {"sequenceId", "requestParams", "resultParams"},
	"wrappedLogV1": {"service", "serviceId", "stack", "stackId"},
}

type pair struct {
	name          string
	conjure, wire any
}

// TestSchemaMatchesConjure requires that each fully populated log type encodes to the same JSON paths and value kinds
// as the generated type, apart from the omitted fields.
func TestSchemaMatchesConjure(t *testing.T) {
	pairs := logPairs(t)
	for _, p := range slices.Clone(pairs) {
		pairs = append(pairs, wrappedPair(t, p))
	}
	for _, p := range pairs {
		t.Run(p.name, func(t *testing.T) {
			require.Equal(t, shape(t, p.name, p.conjure, true), shape(t, p.name, p.wire, false))
		})
	}
}

func logPairs(t *testing.T) []pair {
	var generic logging.GenericDiagnostic
	objecttest.Fill(t, &generic)
	var threadDump logging.ThreadDumpV1
	objecttest.Fill(t, &threadDump)

	conjureEnums := []any{logging.New_AuditResult(logging.AuditResult_SUCCESS), logging.New_AuditProducer(logging.AuditProducer_SERVER)}
	pairs := []pair{
		fill(t, "serviceLogV1", &logging.ServiceLogV1{}, &logtypes.ServiceLogV1{}, nil, nil),
		fill(t, "eventLogV2", &logging.EventLogV2{}, &logtypes.EventLogV2{}, nil, nil),
		fill(t, "metricLogV1", &logging.MetricLogV1{}, &logtypes.MetricLogV1{}, nil, nil),
		fill(t, "traceLogV1", &logging.TraceLogV1{}, &logtypes.TraceLogV1{}, nil, nil),
		fill(t, "requestLogV2", &logging.RequestLogV2{}, &logtypes.RequestLogV2{}, nil, nil),
		fill(t, "auditLogV2", &logging.AuditLogV2{}, &logtypes.AuditLogV2{}, conjureEnums, nil),
		fill(t, "auditLogV3", &logging.AuditLogV3{}, &logtypes.AuditLogV3{}, conjureEnums, nil),
	}
	for variant, d := range map[string]struct {
		conjure logging.Diagnostic
		wire    logtypes.DiagnosticLogV1Diagnostic
	}{
		"generic":    {logging.NewDiagnosticFromGeneric(generic), logtypes.DiagnosticLogV1Diagnostic{Type: "generic", Generic: &generic}},
		"threadDump": {logging.NewDiagnosticFromThreadDump(threadDump), logtypes.DiagnosticLogV1Diagnostic{Type: "threadDump", ThreadDump: &threadDump}},
	} {
		p := fill(t, "diagnosticLogV1", &logging.DiagnosticLogV1{}, &logtypes.DiagnosticLogV1{}, []any{d.conjure}, []any{d.wire})
		p.name += "/" + variant
		pairs = append(pairs, p)
	}
	return pairs
}

func fill(t *testing.T, name string, conjure, wire any, conjureOverrides, wireOverrides []any) pair {
	objecttest.Fill(t, conjure, conjureOverrides...)
	objecttest.Fill(t, wire, wireOverrides...)
	return pair{name: name, conjure: conjure, wire: wire}
}

// wrappedPair wraps a log pair in each type's wrapped.1 payload variant.
func wrappedPair(t *testing.T, p pair) pair {
	var conjurePayload logging.WrappedLogV1Payload
	var wirePayload logtypes.WrappedLogV1Payload
	switch c := p.conjure.(type) {
	case *logging.ServiceLogV1:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromServiceLogV1(*c), logtypes.WrappedLogV1Payload{Type: "serviceLogV1", ServiceLogV1: p.wire.(*logtypes.ServiceLogV1)}
	case *logging.EventLogV2:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromEventLogV2(*c), logtypes.WrappedLogV1Payload{Type: "eventLogV2", EventLogV2: p.wire.(*logtypes.EventLogV2)}
	case *logging.MetricLogV1:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromMetricLogV1(*c), logtypes.WrappedLogV1Payload{Type: "metricLogV1", MetricLogV1: p.wire.(*logtypes.MetricLogV1)}
	case *logging.TraceLogV1:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromTraceLogV1(*c), logtypes.WrappedLogV1Payload{Type: "traceLogV1", TraceLogV1: p.wire.(*logtypes.TraceLogV1)}
	case *logging.RequestLogV2:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromRequestLogV2(*c), logtypes.WrappedLogV1Payload{Type: "requestLogV2", RequestLogV2: p.wire.(*logtypes.RequestLogV2)}
	case *logging.AuditLogV2:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromAuditLogV2(*c), logtypes.WrappedLogV1Payload{Type: "auditLogV2", AuditLogV2: p.wire.(*logtypes.AuditLogV2)}
	case *logging.AuditLogV3:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromAuditLogV3(*c), logtypes.WrappedLogV1Payload{Type: "auditLogV3", AuditLogV3: p.wire.(*logtypes.AuditLogV3)}
	case *logging.DiagnosticLogV1:
		conjurePayload, wirePayload = logging.NewWrappedLogV1PayloadFromDiagnosticLogV1(*c), logtypes.WrappedLogV1Payload{Type: "diagnosticLogV1", DiagnosticLogV1: p.wire.(*logtypes.DiagnosticLogV1)}
	default:
		t.Fatalf("no wrapped.1 variant for %T", p.conjure)
	}
	var conjure logging.WrappedLogV1
	objecttest.Fill(t, &conjure, conjurePayload)
	var wire logtypes.WrappedLogV1
	objecttest.Fill(t, &wire, wirePayload)
	return pair{name: "wrappedLogV1/" + p.name, conjure: &conjure, wire: &wire}
}

// shape maps each JSON path in the encoded value to the kind of value at that path. Generated values encode with JSON
// v1 and omit the fields in omitted.
func shape(t *testing.T, name string, value any, generated bool) map[string]string {
	var b []byte
	var err error
	if generated {
		b, err = jsonv1.Marshal(value)
	} else {
		b, err = json.Marshal(value)
	}
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(b, &decoded))
	if generated {
		strip(decoded, name)
	}
	out := map[string]string{}
	walk(decoded, "", out)
	return out
}

func strip(record map[string]any, name string) {
	root, inner, _ := strings.Cut(name, "/")
	for _, key := range omitted[root] {
		delete(record, key)
	}
	if root == "wrappedLogV1" {
		innerRoot, _, _ := strings.Cut(inner, "/")
		strip(record["payload"].(map[string]any)[innerRoot].(map[string]any), inner)
	}
}

func walk(value any, path string, out map[string]string) {
	switch v := value.(type) {
	case map[string]any:
		out[path] = "object"
		for _, k := range slices.Sorted(maps.Keys(v)) {
			walk(v[k], path+"."+k, out)
		}
	case []any:
		out[path] = "array"
		for i, e := range v {
			walk(e, fmt.Sprintf("%s[%d]", path, i), out)
		}
	case string:
		out[path] = "string"
	case float64:
		out[path] = "number"
	case bool:
		out[path] = "bool"
	case nil:
		out[path] = "null"
	default:
		out[path] = fmt.Sprintf("%T", v)
	}
}
