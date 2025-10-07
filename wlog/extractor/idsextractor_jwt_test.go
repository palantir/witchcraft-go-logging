// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testJWTJohn = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwic2lkIjoiSm9obiBEb2UiLCJqdGkiOiIzNTBhY2MwYS0zNmU2LTQ3MWItOTUxOC1lODBiY2VlYjFmYzQiLCJvcmciOiJhNTY0OWJhNS04MDAyLTRhYzAtYTZkZi1hNTE3ZThlMDFjMTEiLCJpYXQiOjE3NTk4NDgwMzQsImV4cCI6MTc1OTg1MTYzNH0.hL32VOaNMzKxggOtQ_j4Yp4WxfInub3Sqt2VypWeHdY"
	testJWTJane = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwic2lkIjoiSmFuZSBEb2UiLCJqdGkiOiIzNTBhY2MwYS0zNmU2LTQ3MWItOTUxOC1lODBiY2VlYjFmYzQiLCJvcmciOiJhNTY0OWJhNS04MDAyLTRhYzAtYTZkZi1hNTE3ZThlMDFjMTEiLCJpYXQiOjE3NTk4NDgwMzQsImV4cCI6MTc1OTg1MTYzNH0.vkI_jtLu53gak3doWhMNjsFmHL5gymvO10pfZN301YQ"
)

func Test_uuidFromBase64StdEncodedString(t *testing.T) {
	in := "vp9kXVLgSem6MdsyknYV2w=="
	expected := "be9f645d-52e0-49e9-ba31-db32927615db"
	assert.Equal(t, expected, uuidFromBase64StdEncodedString(in))
}

func Test_jwtRequestIDsExtractor_ExtractIDs(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want map[string]string
	}{
		{
			name: "empty",
			req:  &http.Request{},
			want: map[string]string{
				UIDKey:     "",
				SIDKey:     "",
				TokenIDKey: "",
				OrgIDKey:   "",
			},
		},
		{
			name: "authorization header",
			req: &http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {"Bearer " + testJWTJohn},
				},
			},
			want: map[string]string{
				UIDKey:     "1234567890",
				SIDKey:     "John Doe",
				TokenIDKey: "350acc0a-36e6-471b-9518-e80bceeb1fc4",
				OrgIDKey:   "a5649ba5-8002-4ac0-a6df-a517e8e01c11",
			},
		},
		{
			name: "websocket protocols header",
			req: &http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey("Sec-WebSocket-Protocol"): {"Bearer-" + testJWTJohn},
				},
			},
			want: map[string]string{
				UIDKey:     "1234567890",
				SIDKey:     "John Doe",
				TokenIDKey: "350acc0a-36e6-471b-9518-e80bceeb1fc4",
				OrgIDKey:   "a5649ba5-8002-4ac0-a6df-a517e8e01c11",
			},
		},
		{
			name: "websocket protocols and authorization header",
			req: &http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey("Sec-WebSocket-Protocol"): {"Bearer-" + testJWTJohn},
					textproto.CanonicalMIMEHeaderKey("Authorization"):          {"Bearer " + testJWTJane},
				},
			},
			want: map[string]string{
				UIDKey:     "1234567890",
				SIDKey:     "Jane Doe",
				TokenIDKey: "350acc0a-36e6-471b-9518-e80bceeb1fc4",
				OrgIDKey:   "a5649ba5-8002-4ac0-a6df-a517e8e01c11",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &jwtRequestIDsExtractor{}
			assert.Equal(t, tt.want, e.ExtractIDs(tt.req))
		})
	}
}
