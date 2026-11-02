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

package wlog

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
)

// Generated union marshalers encode through JSON v1. Use their visitors for
// variant validation, then encode the selected value with JSON v2.
var jsonMarshalers = json.JoinMarshalers(
	errorMarshalers,
	paramMapMarshalers,
	json.MarshalToFunc(func(enc *jsontext.Encoder, diagnostic logging.Diagnostic) error {
		return diagnostic.Accept(&payloadEncoder{enc: enc})
	}),
	json.MarshalToFunc(func(enc *jsontext.Encoder, payload logging.WrappedLogV1Payload) error {
		return payload.Accept(&payloadEncoder{enc: enc})
	}),
)

type payloadEncoder struct {
	enc *jsontext.Encoder
}

func marshalPayload[T any](enc *jsontext.Encoder, name string, value T) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.String("type")); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.String(name)); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.String(name)); err != nil {
		return err
	}
	if err := json.MarshalEncode(enc, value); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndObject)
}

func (p *payloadEncoder) VisitServiceLogV1(value logging.ServiceLogV1) error {
	return marshalPayload(p.enc, "serviceLogV1", value)
}

func (p *payloadEncoder) VisitRequestLogV2(value logging.RequestLogV2) error {
	return marshalPayload(p.enc, "requestLogV2", value)
}

func (p *payloadEncoder) VisitTraceLogV1(value logging.TraceLogV1) error {
	return marshalPayload(p.enc, "traceLogV1", value)
}

func (p *payloadEncoder) VisitEventLogV2(value logging.EventLogV2) error {
	return marshalPayload(p.enc, "eventLogV2", value)
}

func (p *payloadEncoder) VisitMetricLogV1(value logging.MetricLogV1) error {
	return marshalPayload(p.enc, "metricLogV1", value)
}

func (p *payloadEncoder) VisitAuditLogV2(value logging.AuditLogV2) error {
	return marshalPayload(p.enc, "auditLogV2", value)
}

func (p *payloadEncoder) VisitAuditLogV3(value logging.AuditLogV3) error {
	return marshalPayload(p.enc, "auditLogV3", value)
}

func (p *payloadEncoder) VisitDiagnosticLogV1(value logging.DiagnosticLogV1) error {
	return marshalPayload(p.enc, "diagnosticLogV1", value)
}

func (p *payloadEncoder) VisitGeneric(value logging.GenericDiagnostic) error {
	return marshalPayload(p.enc, "generic", value)
}

func (p *payloadEncoder) VisitThreadDump(value logging.ThreadDumpV1) error {
	return marshalPayload(p.enc, "threadDump", value)
}

func (p *payloadEncoder) VisitUnknown(name string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", name)
}
