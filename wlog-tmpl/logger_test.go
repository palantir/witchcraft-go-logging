// Copyright (c) 2020 Palantir Technologies. All rights reserved.
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

package wlogtmpl_test

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogtmpl "github.com/palantir/witchcraft-go-logging/wlog-tmpl"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logs"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	customFormatter := func(logType, tmpl string) map[logentryformatter.LogType]logentryformatter.Formatter {
		typ := logentryformatter.LogType(logType)
		f, err := logs.Formatter(typ, tmpl)
		require.NoError(t, err)
		return map[logentryformatter.LogType]logentryformatter.Formatter{typ: f}
	}

	for _, test := range []struct {
		Name     string
		Config   *wlogtmpl.Config
		LogFn    func(ctx context.Context)
		Expected *regexp.Regexp
	}{
		{
			Name: "svc1log info no substitution",
			LogFn: func(ctx context.Context) {
				svc1log.FromContext(ctx).Info("An info message about foo", svc1log.UnsafeParam("foo", "bar"))
			},
			Expected: regexp.MustCompile(`^INFO  \[[0-9TZ:.-]+] origin: An info message about foo \(foo: bar\)$`),
		},
		{
			Name: "svc1log error no substitution",
			LogFn: func(ctx context.Context) {
				svc1log.FromContext(ctx).Error("An error message about foo", svc1log.UnsafeParam("foo", "bar"))
			},
			Expected: regexp.MustCompile(`^ERROR \[[0-9TZ:.-]+] origin: An error message about foo \(foo: bar\)$`),
		},
		{
			Name: "svc1log info with substitution",
			LogFn: func(ctx context.Context) {
				svc1log.FromContext(ctx).Info("An info message about {}", svc1log.UnsafeParam("0", "bar"))
			},
			Expected: regexp.MustCompile(`^INFO  \[[0-9TZ:.-]+] origin: An info message about bar \(0: bar\)$`),
		},
		{
			Name: "evt2log",
			LogFn: func(ctx context.Context) {
				evt2log.FromContext(ctx).Event("MY_EVENT", evt2log.UnsafeParam("foo", "bar"))
			},
			Expected: regexp.MustCompile(`^\[[0-9TZ:.-]+] MY_EVENT \(foo: bar\)$`),
		},
		{
			Name: "custom config",
			Config: &wlogtmpl.Config{
				FormatterMap: customFormatter("service.1", `CUSTOM {{printf "%-5s" .Level}} {{printf "%-26s" (printf "[%s]" .Time)}} {{.Origin}}: {{.Message}}`),
			},
			LogFn: func(ctx context.Context) {
				svc1log.FromContext(ctx).Info("An info message about foo", svc1log.UnsafeParam("foo", "bar"))
			},
			Expected: regexp.MustCompile(`^CUSTOM INFO  \[[0-9TZ:.-]+] origin: An info message about foo$`),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			out := &bytes.Buffer{}
			ctx := context.Background()
			provider := wlogtmpl.LoggerProvider(test.Config)
			ctx = evt2log.WithLogger(ctx, evt2log.NewFromCreator(out, provider.NewLogger))
			ctx = svc1log.WithLogger(ctx, svc1log.NewFromCreator(out, wlog.InfoLevel, provider.NewLeveledLogger, svc1log.Origin("origin")))

			test.LogFn(ctx)

			outStr := strings.TrimSpace(out.String())
			assert.Regexp(t, test.Expected, outStr)
		})
	}
}
