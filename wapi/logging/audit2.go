// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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

package logging

import (
	"github.com/palantir/pkg/datetime"
)

type AuditResult string

const (
	AuditResultSUCCESS      AuditResult = "SUCCESS"
	AuditResultUNAUTHORIZED AuditResult = "UNAUTHORIZED"
	AuditResultERROR        AuditResult = "ERROR"
)

// audit.2

type AuditLogV2 struct {
	// "audit.2"
	Type string `json:"type"`
	// RFC3339Nano timestamp when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// Name of the audit event, e.g. PUT_FILE
	Name string `json:"name"`
	// Indicates whether the request was successful or the type of failure, e.g. ERROR or UNAUTHORIZED
	Result AuditResult `json:"result"`
	// User id (if available). This is the most downstream caller.
	Uid *UserId `json:"uid,omitempty"`
	// Session id (if available)
	Sid *SessionId `json:"sid,omitempty"`
	// Token id (if available)
	TokenId *TokenId `json:"tokenId,omitempty"`
	// Organization ID (if available)
	OrgId *OrgId `json:"orgId,omitempty"`
	// Zipkin trace id (if available)
	TraceId *TraceId `json:"traceId,omitempty"`
	// All users upstream of the user currently taking an action. The first element in this list is the uid of the most upstream caller. This list does not include the `uid`.
	OtherUids []UserId `json:"otherUids,omitempty"`
	// Best-effort identifier of the originating machine, e.g. an IP address, a Kubernetes node identifier, or similar
	Origin *string `json:"origin,omitempty"`
	// The parameters known at method invocation time.
	RequestParams map[string]any `json:"requestParams,omitempty"`
	// Information derived within a method, commonly parts of the return value.
	ResultParams map[string]any `json:"resultParams,omitempty"`
}

func (log *AuditLogV2) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Uid = nil
	log.Sid = nil
	log.TokenId = nil
	log.TraceId = nil
	log.OtherUids = nil
	log.Origin = nil
	log.Name = ""
	log.Result = ""
	clear(log.ResultParams)
	clear(log.RequestParams)
}
