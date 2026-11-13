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

package trc1log_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wlogtests"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/stretchr/testify/require"
)

func TestEntryEndpointPresence(t *testing.T) {
	provider := wlogtests.MapJSONProvider()
	for _, tc := range []struct {
		name     string
		endpoint *wtracing.Endpoint
		want     map[string]any
	}{
		{name: "missing"},
		{name: "empty", endpoint: &wtracing.Endpoint{}, want: map[string]any{"serviceName": ""}},
		{name: "empty IP", endpoint: &wtracing.Endpoint{IPv4: net.IP{}}, want: map[string]any{"serviceName": ""}},
		{name: "service", endpoint: &wtracing.Endpoint{ServiceName: "service"}, want: map[string]any{"serviceName": "service"}},
		{name: "address", endpoint: &wtracing.Endpoint{IPv4: net.IPv4(127, 0, 0, 1)}, want: map[string]any{"serviceName": "", "ipv4": "127.0.0.1"}},
	} {
		for name, emit := range map[string]func(io.Writer, wtracing.SpanModel){
			"ToParams": func(w io.Writer, span wtracing.SpanModel) {
				provider.NewLogger(w).Log(trc1log.ToParams(span, nil)...)
			},
			"logger": func(w io.Writer, span wtracing.SpanModel) {
				trc1log.NewFromCreator(w, provider.NewLogger).Log(span)
			},
			"wrapped": func(w io.Writer, span wtracing.SpanModel) {
				wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, "entity", "version").Trace().Log(span)
			},
		} {
			t.Run(tc.name+"/"+name, func(t *testing.T) {
				var buf bytes.Buffer
				emit(&buf, wtracing.SpanModel{Kind: wtracing.Server, Timestamp: time.Unix(1, 0), LocalEndpoint: tc.endpoint})
				var record map[string]any
				require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
				if name == "wrapped" {
					record = record["payload"].(map[string]any)["traceLogV1"].(map[string]any)
				}
				annotations := record["span"].(map[string]any)["annotations"].([]any)
				require.Len(t, annotations, 2)
				for _, annotation := range annotations {
					fields := annotation.(map[string]any)
					if tc.want == nil {
						require.NotContains(t, fields, "endpoint")
					} else {
						require.Equal(t, tc.want, fields["endpoint"])
					}
				}
			})
		}
	}
}
