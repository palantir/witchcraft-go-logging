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

package extractor

import (
	"net/http"
	"strings"
)

// newIDsFromWSExtractor creates an extractor that sets the UIDKey, SIDKey, TokenIDKey and OrgIDKey keys to have the values
// parsed from the WebSocket request containing the bearer token in the "Sec-WebSocket-Protocol" header with the `Bearer-<token>` format.
// The JWT's "sub" field is used as the UID, the "sid" field is used as the SID, the "jti" field is used as the tokenID and the "org" field is used as the orgID.
func newIDsFromWSExtractor() IDsFromRequest {
	return &wsRequestIDsExtractor{}
}

type wsRequestIDsExtractor struct{}

func (e *wsRequestIDsExtractor) ExtractIDs(req *http.Request) map[string]string {
	const bearerTokenPrefix = "Bearer-"

	var uid, sid, tokenID, orgID string
	authContent := req.Header.Get("Sec-WebSocket-Protocol")
	if strings.HasPrefix(authContent, bearerTokenPrefix) {
		uid, sid, tokenID, orgID, _ = idsFromJWT(authContent[len(bearerTokenPrefix):])
	}
	return map[string]string{
		UIDKey:     uid,
		SIDKey:     sid,
		TokenIDKey: tokenID,
		OrgIDKey:   orgID,
	}
}
